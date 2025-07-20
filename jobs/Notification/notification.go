package notification

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	controllers "github.com/ankur12345678/uptime-monitor/Controllers"
	models "github.com/ankur12345678/uptime-monitor/Models"
	"github.com/ankur12345678/uptime-monitor/jobs"
	"github.com/ankur12345678/uptime-monitor/pkg/aws"
	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"github.com/ankur12345678/uptime-monitor/pkg/sendgrid"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type notificationJob struct {
	jobs.JobInput
	config    Config
	channel   chan *types.Message
	wg        sync.WaitGroup
	sqsClient *sqs.Client
}

type Config struct {
	WorkerCount   int
	JobTimeout    time.Duration
	ChannelBuffer int
}

func DefaultConfig() Config {
	return Config{
		WorkerCount:   128,
		JobTimeout:    2 * time.Minute,
		ChannelBuffer: 1000,
	}
}

func New(input jobs.JobInput, config Config) *notificationJob {
	return &notificationJob{
		JobInput: input,
		config:   config,
		channel:  make(chan *types.Message, 1000),
		wg:       sync.WaitGroup{},
	}
}

func (nj *notificationJob) PrepareSQSClient(ctx context.Context) *sqs.Client {
	awsConfig := aws.LoadAWSConfig(nj.BaseController.Config.AwsProfile, nj.BaseController.Config.AwsRegion)
	sqsClient := aws.NewClient(awsConfig)

	return sqsClient
}

func (nj *notificationJob) StartPullingNotificationsFromQueue(ctx context.Context) {
	sqsClient := nj.PrepareSQSClient(ctx)

	nj.sqsClient = sqsClient

	for {
		//this message wont be visbile to other consumers for visibilty period (30sec)
		msg, err := aws.ReceiveMessage(sqsClient, nj.BaseController.Config.AwsQueueUrl)
		if err != nil {
			logger.Error("error in fetching message from queue | err: ", err)
		}
		if msg == nil {
			continue
		}

		select {
		case nj.channel <- msg:
		case <-ctx.Done():
			logger.Error("context error | err: ", ctx.Err())
			return
		}
	}
}

func HandleMessage(msg *types.Message) (*jobs.SQSIncidentEventType, *string) {
	var formattedMsg jobs.SQSIncidentEventType
	err := json.Unmarshal([]byte(*msg.Body), &formattedMsg)
	if err != nil {
		log.Println("Failed to parse message body:", err)
		return nil, nil
	}
	logger.Infof("Message received : %+v", formattedMsg)
	return &formattedMsg, msg.ReceiptHandle
}

func Start(ctrl *controllers.BaseController) {
	cfg := DefaultConfig()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.JobTimeout)
	defer cancel()

	nj := New(jobs.JobInput{BaseController: *ctrl}, cfg)

	go nj.StartPullingNotificationsFromQueue(ctx)

	for worker := 0; worker < nj.config.WorkerCount; worker++ {
		nj.wg.Add(1)
		go func(ctx context.Context, workerID int) {
			defer nj.wg.Done()
			//recovery code
			defer func() {
				if r := recover(); r != nil {
					// r is the panic payload (error, string, etc.)
					logger.Error("goroutine panicked | err: ", r)
				}
			}()

			for {
				select {
				case msg := <-nj.channel:
					var (
						err error
					)
					formattedMsg, receiptHandle := HandleMessage(msg)
					if formattedMsg.Phone != "" {
						logger.Error("not supporting notification on phone number!")
					}
					if formattedMsg.Email != "" {
						err = nj.handleEmail(ctx, formattedMsg)
					}

					if err == nil {
						//since processing of email is done therfore delete this msg from queue
						if receiptHandle != nil {
							err := aws.DeleteMessage(nj.sqsClient, nj.BaseController.Config.AwsQueueUrl, receiptHandle)
							if err != nil {
								logger.Error("error in deleting msg from sqs | err: ", err)
								return
							}
						}
					}
				case <-ctx.Done():
					logger.Error("context error | err: ", ctx.Err())
					return
				}
			}

		}(ctx, worker)
	}
	nj.wg.Wait()

}

func (nj *notificationJob) handleEmail(ctx context.Context, formattedMsg *jobs.SQSIncidentEventType) error {
	var (
		incidentEventsRepo = models.InitIncidentEventsRepo(nj.DB)
	)
	err := sendgrid.SendEmail(formattedMsg.Email, "", nj.BaseController.Config, sendgrid.EmailData{WebsiteURL: formattedMsg.WebsiteURL, Status: formattedMsg.Status, Year: time.Now().Year()})
	if err != nil {
		logger.Error("error in sending email notification | err: ", err)
		eventUpdateErr := incidentEventsRepo.UpdateWithTx(nj.DB.WithContext(ctx), &models.IncidentEvent{UUID: formattedMsg.IncidentEventID}, &models.IncidentEvent{EventStatus: models.EventStatusFailed})
		if eventUpdateErr != nil {
			logger.Error("error updating incident status to success | err: ", err)
			return eventUpdateErr
		}
		return err
	} else {
		err := incidentEventsRepo.UpdateWithTx(nj.DB.WithContext(ctx), &models.IncidentEvent{UUID: formattedMsg.IncidentEventID}, &models.IncidentEvent{EventStatus: models.EventStatusDelivered})
		if err != nil {
			logger.Error("error updating incident status to success | err: ", err)
			return err
		}

	}

	return nil
}
