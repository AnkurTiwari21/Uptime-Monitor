package websitepicker

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	controllers "github.com/ankur12345678/uptime-monitor/Controllers"
	models "github.com/ankur12345678/uptime-monitor/Models"
	"github.com/ankur12345678/uptime-monitor/jobs"
	"github.com/ankur12345678/uptime-monitor/pkg/aws"
	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"gorm.io/gorm"
)

// Config holds job configuration parameters
type Config struct {
	BatchSize          int
	WorkerCount        int
	ChannelBuffer      int
	JobTimeout         time.Duration
	HealthCheckTimeout time.Duration
}

// DefaultConfig returns default configuration values
func DefaultConfig() Config {
	return Config{
		BatchSize:          10,
		WorkerCount:        128,
		ChannelBuffer:      1000,
		JobTimeout:         2 * time.Minute,
		HealthCheckTimeout: 50 * time.Second,
	}
}

type websitePickerJob struct {
	jobs.JobInput
	config       Config
	websitesChan chan models.Website
	wg           sync.WaitGroup
	httpClient   http.Client
	sqsClient    *sqs.Client
}

func New(input jobs.JobInput, config Config) *websitePickerJob {
	//TODO: figure out a good timeout period
	// transport := &http.Transport{
	// 	// Connection pool sizes (tweak to taste)
	// 	MaxIdleConns:        100,
	// 	MaxIdleConnsPerHost: 10,
	// 	IdleConnTimeout:     30 * time.Second,

	// 	// Hard limits for the slow parts
	// 	DialContext: (&net.Dialer{
	// 		Timeout:   3 * time.Second, // TCP handshake (and DNS) timeout
	// 		KeepAlive: 30 * time.Second,
	// 	}).DialContext,
	// 	TLSHandshakeTimeout:   3 * time.Second, // separate cap for TLS
	// 	ResponseHeaderTimeout: 5 * time.Second, // wait this long for first byte
	// 	// (body read still governed by ctx or client.Timeout)
	// }
	awsCfg := aws.LoadAWSConfig(input.BaseController.Config.AwsProfile, input.BaseController.Config.AwsRegion)

	sqsClient := aws.NewClient(awsCfg)
	return &websitePickerJob{
		JobInput:     input,
		config:       config,
		websitesChan: make(chan models.Website, config.ChannelBuffer),
		wg:           sync.WaitGroup{},
		httpClient: http.Client{
			Timeout: config.HealthCheckTimeout,
			// Transport: transport,
		},
		sqsClient: sqsClient,
	}
}

func (w *websitePickerJob) UpdateAllWebsiteLastCheckedTime(ctx context.Context, tx *gorm.DB, websites []models.Website) error {
	ids := make([]uint, 0, len(websites))
	for _, w := range websites {
		ids = append(ids, w.ID)
	}

	err := tx.WithContext(ctx).Model(&models.Website{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{"last_checked_at": time.Now()}).
		Error

	if err != nil {
		logger.Error("error in updating last_checked_at | err: ", err)
		return err
	}

	return nil
}

func normalizeURL(url string) string {
	if !strings.Contains(url, "://") {
		return "https://" + url
	}
	return url
}

func (w *websitePickerJob) FetchWebsitesForJob(ctx context.Context) error {
	var (
		websiteRepo = models.InitWebsiteRepo(w.DB)
	)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		websites, tx, err := websiteRepo.FetchWebsitesInBulk(ctx, 10)

		if err != nil {
			logger.Error("error in fetching webistes | err: ", err)
			tx.Rollback()
			return err
		}

		if len(websites) == 0 {
			break
		}

		err = w.UpdateAllWebsiteLastCheckedTime(ctx, tx, websites)
		if err != nil {
			logger.Error("error in committing transactions | err: ", err)
			tx.Rollback()
			return err
		}

		err = tx.Commit().Error
		if err != nil {
			logger.Error("error in committing transactions | err: ", err)
			tx.Rollback()
			return err
		}

		for _, site := range websites {
			select {
			case w.websitesChan <- site:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

	}
	return nil
}

func (w *websitePickerJob) CloseChannel() {
	close(w.websitesChan)
}

func identifyCummulativeStatusBasedOnPastRecords(statusRecords []string) models.HealthStatus {
	var (
		flag = models.Unhealthy
	)

	for _, val := range statusRecords {
		if val == string(models.Healthy) {
			flag = models.Healthy
			break
		}
	}

	return flag
}

func isUnhealthyStatus(code int) bool {
	return code >= 400 || code == 0
}

func (w *websitePickerJob) CreateOrResolveIncident(ctx context.Context, webisteID uint, statusCode int, latency time.Duration) {
	var (
		alertConfigRepo = models.InitAlertConfigRepo(w.DB)
		logsRepo        = models.InitLogsRepo(w.DB)
		incidentsRepo   = models.InitIncidentsRepo(w.DB)
		status          models.HealthStatus
	)

	alertConfig, err := alertConfigRepo.GetWithTx(w.DB.WithContext(ctx), &models.AlertConfig{WebsiteID: webisteID})
	if err != nil {
		logger.Error("error in fetching alert config for this webiste | err: ", err)
		return
	}

	if latency.Milliseconds() >= int64(alertConfig.LatencyThreshold) || isUnhealthyStatus(statusCode) {
		status = models.Unhealthy
	} else {
		status = models.Healthy
	}

	err = logsRepo.Create(ctx, models.Log{
		WebsiteId:    webisteID,
		StatusCode:   uint(statusCode),
		LatencyInMS:  uint(latency.Milliseconds()),
		HealthStatus: string(status),
	})
	if err != nil {
		logger.Error("error in creating log | err: ", err)
		return
	}

	statusRecords, err := logsRepo.FetchPastRecordStatusByWebsiteID(ctx, uint(alertConfig.FailureThreshold), webisteID)
	if err != nil {
		logger.Error("error in fetching log records for incidents | err: ", err)
		return
	}

	if len(statusRecords) != alertConfig.FailureThreshold {
		//since we dont have enough data to create incidents/notify therefore get back from here
		return
	}

	currentCummulativeStatus := identifyCummulativeStatusBasedOnPastRecords(statusRecords)

	//fetch past incident
	pastStatus, err := incidentsRepo.GetWithTx(w.DB.WithContext(ctx), &models.Incident{WebsiteId: webisteID})
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error("error in fetching previous incidents | err: ", err)
		return
	}

	if err == gorm.ErrRecordNotFound {
		if currentCummulativeStatus == models.Unhealthy {
			//enter record in incident table and notify to user
			err := incidentsRepo.Create(w.DB.WithContext(ctx), models.Incident{WebsiteId: webisteID, HealthStatus: string(status)})
			if err != nil {
				logger.Error("error in creating incident record | err: ", err)
				return
			}
			logger.Info("notifying user that website is down!")
			w.notifyUser(ctx, alertConfig.ID, status, webisteID)
		}
	} else {
		if pastStatus.HealthStatus == string(models.Unhealthy) && status == (models.Unhealthy) {
			logger.Info("notifying user that website is down!")
			w.notifyUser(ctx, alertConfig.ID, status, webisteID)
		} else if pastStatus.HealthStatus == string(models.Unhealthy) && status == (models.Healthy) {
			//notufy user that webiste is up and delete the incident
			err := incidentsRepo.DeleteWithTx(w.DB.WithContext(ctx), &models.Incident{ID: pastStatus.ID})
			if err != nil {
				logger.Error("error in deleting incident record | err: ", err)
				return
			}

			//push to SQS for notification
			logger.Info("notifying user that website is up!")
			w.notifyUser(ctx, alertConfig.ID, status, webisteID)
		}
	}

}

func (w *websitePickerJob) DoHealthCheck(parentCtx context.Context, website models.Website) {
	//defining a child context
	childCtx, cancel := context.WithTimeout(parentCtx, w.config.HealthCheckTimeout)
	defer cancel()

	req, _ := http.NewRequestWithContext(childCtx, http.MethodGet, normalizeURL(website.WebsiteURL), nil)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	start := time.Now()
	resp, err := w.httpClient.Do(req)
	latency := time.Since(start)

	if resp != nil {
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)

		//check if incident should be created/already present and notify them
		w.CreateOrResolveIncident(childCtx, website.ID, resp.StatusCode, latency)
	}
	if err != nil {
		logger.Error("error while checking website's health | err: ", err)
		return
	}

}

func ProcessWebsitesJob(ctrl controllers.BaseController) {
	config := DefaultConfig()

	job := New(jobs.JobInput{
		BaseController: ctrl,
	}, config)

	ctx, cancel := context.WithTimeout(context.Background(), config.JobTimeout)
	defer cancel()

	go func() {
		defer job.CloseChannel()
		defer func() {
			if r := recover(); r != nil {
				// r is the panic payload (error, string, etc.)
				logger.Error("goroutine panicked | err: ", r)
			}
		}()
		err := job.FetchWebsitesForJob(ctx)
		if err != nil {
			logger.Error("website fetching error: ", err)
		}
	}()

	for i := 0; i < config.WorkerCount; i++ {
		job.wg.Add(1)
		go func(workerId int) {
			defer job.wg.Done()
			defer func() {
				if r := recover(); r != nil {
					// r is the panic payload (error, string, etc.)
					logger.Error("goroutine panicked | err: ", r)
				}
			}()
			for {
				select {
				case site, ok := <-job.websitesChan:
					if !ok {
						// channel closed
						return
					}
					job.DoHealthCheck(ctx, site)

				case <-ctx.Done():
					// job timeout or cancellation
					return
				}
			}
		}(i + 1)
	}

	job.wg.Wait()
}

func (w *websitePickerJob) notifyUser(ctx context.Context, alertConfigID uint, healthStatus models.HealthStatus, websiteID uint) {
	var (
		alertTargetRepo    = models.InitAlertTargetRepo(w.DB.WithContext(ctx))
		incidentEventsRepo = models.InitIncidentEventsRepo(w.DB)
		websiteRepo        = models.InitWebsiteRepo(w.DB)
	)

	website, err := websiteRepo.GetWithTx(&models.Website{ID: websiteID}, w.DB.WithContext(ctx))
	if err != nil {
		logger.Error("error in getting the webiste with given websiteId | err: ", err)
		return
	}

	alertTargets, err := alertTargetRepo.GetAllByAlertConfigID(alertConfigID)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error("error in getting the alert targets for this config | err: ", err)
		return
	}

	if err == gorm.ErrRecordNotFound {
		logger.Error("no alert targets are present for this webiste")
		return
	}

	for _, target := range alertTargets {
		if target.TargetType == models.TargetTypeSMS {
			logger.Error("SMS notifications is not supported currently!")
		} else {
			incidentEventMsgForQueue := jobs.SQSIncidentEventType{
				WebsiteURL: website.WebsiteURL,
				Phone:      "",
				Email:      target.TargetValue,
				Status:     string(healthStatus),
			}
			incidentEvent := models.IncidentEvent{
				HealthStatus:  string(healthStatus),
				WebsiteURL:    website.WebsiteURL,
				EventStatus:   models.EventStatusPending,
				AlertTargetId: target.ID,
			}

			err := incidentEventsRepo.CreateWithTx(w.DB.WithContext(ctx), &incidentEvent)
			if err != nil {
				logger.Error("error in creating incident event | err: ", err)
				return
			}

			incidentEventMsgForQueue.IncidentEventID = incidentEvent.UUID

			err = aws.SendMessage(w.sqsClient, w.BaseController.Config.AwsQueueUrl, &incidentEventMsgForQueue)
			if err != nil {
				logger.Error("error in sending incident event to SQS for notification | err: ", err)
				return
			}
		}
	}

}
