package controllers

import (
	"net/http"

	models "github.com/ankur12345678/uptime-monitor/Models"
	"github.com/ankur12345678/uptime-monitor/pkg/constants"
	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (b *BaseController) UpdateAlertConfig(c *gin.Context) {
	var (
		request         = UpdateAlertConfigRequest{}
		alertConfigRepo = models.InitAlertConfigRepo(b.DB)
	)

	err := c.ShouldBindJSON(&request)
	if err != nil {
		logger.Error("error in binding request | err: ", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, UpdateAlertConfigResponse{
			Status:  constants.GENERIC_FAILURE_RESPONSE,
			Message: "Something went wrong. Please try again",
		})
		return
	}

	//TODO: validate request

	config, err := alertConfigRepo.GetWithTx(b.DB, &models.AlertConfig{WebsiteID: request.WebsiteId})
	if err != nil {
		logger.Error("error in fetching config from DB | err: ", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, UpdateAlertConfigResponse{
			Status:  constants.GENERIC_FAILURE_RESPONSE,
			Message: "Something went wrong. Please try again",
		})
		return
	}

	//TODO: validate the updation of falsy values like FALSE
	err = alertConfigRepo.Update(&models.AlertConfig{ID: config.ID}, &models.AlertConfig{IsEnabled: request.IsEnabled, LatencyThreshold: int(request.LatencyThreshold), FailureThreshold: int(request.FailureThreshold)})
	if err != nil {
		logger.Error("error in updating config in DB | err: ", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, UpdateAlertConfigResponse{
			Status:  constants.GENERIC_FAILURE_RESPONSE,
			Message: "Something went wrong. Please try again",
		})
		return
	}

	c.JSON(http.StatusOK, UpdateAlertConfigResponse{
		Status:  constants.GENERIC_SUCCESS_RESPONSE,
		Message: "Updated successfully.",
	})
}
