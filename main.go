package main

import (
	"fmt"
	"net/http"
	"time"

	config "github.com/ankur12345678/uptime-monitor/Config"
	controllers "github.com/ankur12345678/uptime-monitor/Controllers"
	migration "github.com/ankur12345678/uptime-monitor/Migration"
	Router "github.com/ankur12345678/uptime-monitor/Router"
	"github.com/ankur12345678/uptime-monitor/jobs"
	notification "github.com/ankur12345678/uptime-monitor/jobs/Notification"
	websitepicker "github.com/ankur12345678/uptime-monitor/jobs/WebsitePicker"
	"github.com/ankur12345678/uptime-monitor/pkg/graceful"
	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"github.com/ankur12345678/uptime-monitor/pkg/validator"
	"github.com/gin-gonic/gin"
)

func main() {
	//defining routes
	logger.Info("Starting Server....")

	//loading cfg
	cfg := config.LoadConfig()
	db := migration.InitDB(cfg)
	ctrl := controllers.BaseController{
		DB:     db,
		Config: cfg,
	}
	controllers.Ctrl = ctrl

	//init redis client
	redisClient := migration.InitRedisClient(ctrl.Config)
	controllers.Ctrl.RedisClient = redisClient

	//seeding data for test
	migration.SeedDB(db)

	job := config.GetJob()

	logger.Info("job: ", job)

	switch jobs.JobName(job) {
	case jobs.MonitorWesbitesJob:
		logger.Infof("****** Starting Job: %s ******", job)
		websitepicker.ProcessWebsitesJob(ctrl)
		logger.Infof("****** Completed Job: %s ******", job)
	case jobs.NotificationJob:
		logger.Infof("****** Starting Job: %s ******", job)
		notification.Start(&ctrl)
		logger.Infof("****** Completed Job: %s ******", job)
	default:
		router := gin.New()
		ctrl.Router = router

		validate, trans, err := validator.InitValidator()
		if err != nil {
			logger.Fatal("Unable to init validator ", err)
		}

		//adding remaining values to the controller
		ctrl.Translator = &trans
		ctrl.Validator = validate

		router.Use(gin.Logger())
		router.Use(gin.Recovery())

		//initializing routes
		Router.RegisterRoutes(ctrl)

		httpServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", ctrl.Config.ServerPort),
			Handler: ctrl.Router,
		}

		graceful := graceful.Graceful{
			HTTPServer:      httpServer,
			ShutdownTimeout: time.Duration(3 * time.Second),
			State:           &graceful.ServerState{},
		}

		banner := `
	
░███     ░███                        ░██████                      
░████   ░████                       ░██   ░██                     
░██░██ ░██░██  ░███████  ░████████  ░██   ░██  ░███████  ░██░████ 
░██ ░████ ░██ ░██    ░██ ░██    ░██  ░██████  ░██    ░██ ░███     
░██  ░██  ░██ ░██    ░██ ░██    ░██ ░██   ░██ ░██    ░██ ░██      
░██       ░██ ░██    ░██ ░██    ░██ ░██   ░██ ░██    ░██ ░██      
░██       ░██  ░███████  ░██    ░██  ░██████   ░███████  ░██      
                                                                  
                                                                  
	`

		graceful.ListenAndServe(banner)
	}

}
