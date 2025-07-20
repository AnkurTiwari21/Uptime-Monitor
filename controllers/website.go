package controllers

import (
	"net/http"
	"time"

	models "github.com/ankur12345678/uptime-monitor/Models"
	"github.com/ankur12345678/uptime-monitor/pkg/constants"
	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (b *BaseController) RegisterWebsite(ctx *gin.Context) {
	var (
		request         RegisterWebsiteRequest
		websiteRepo     = models.InitWebsiteRepo(b.DB)
		userRepo        = models.InitUserRepo(b.DB)
		alertConfigRepo = models.InitAlertConfigRepo(b.DB)
	)

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		logger.Error("error in binding request | err: ", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, RegisterWebsiteResponse{
			Status:  constants.GENERIC_FAILURE_RESPONSE,
			Message: "Something went wrong. Please try again",
		})
		return
	}

	if request.WebsiteURL == "" {
		logger.Error("invalid request")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, RegisterWebsiteResponse{
			Status:  constants.GENERIC_FAILURE_RESPONSE,
			Message: "Please enter valid details",
		})
		return
	}

	email, err := GetEmailFromContext(ctx)
	if err != nil {
		logger.Error("error in getting email from context | err: ", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, RegisterWebsiteResponse{
			Status:  constants.GENERIC_FAILURE_RESPONSE,
			Message: "Something went wrong. Please try again",
		})
		return
	}

	user, err := userRepo.GetByEmail(email)
	if err != nil {
		logger.Error("error in getting user from DB | err: ", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, RegisterWebsiteResponse{
			Status:  constants.GENERIC_FAILURE_RESPONSE,
			Message: "Something went wrong. Please try again",
		})
		return
	}

	tx := b.DB.Begin()

	website := &models.Website{WebsiteURL: request.WebsiteURL, UserId: user.ID}
	err = websiteRepo.Create(website)
	if err != nil {
		logger.Error("error in registering website | err: ", err)
		tx.Rollback()
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, RegisterWebsiteResponse{
			Status:  constants.GENERIC_FAILURE_RESPONSE,
			Message: "Something went wrong. Please try again",
		})
		return
	}

	err = alertConfigRepo.Create(&models.AlertConfig{WebsiteID: website.ID})
	if err != nil {
		logger.Error("error in creating alert config for this website | err: ", err)
		tx.Rollback()
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, RegisterWebsiteResponse{
			Status:  constants.GENERIC_FAILURE_RESPONSE,
			Message: "Something went wrong. Please try again",
		})
		return
	}

	err = tx.Commit().Error
	if err != nil {
		logger.Error("error while commiting transaction | err: ", err)
		tx.Rollback()
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, RegisterWebsiteResponse{
			Status:  constants.GENERIC_FAILURE_RESPONSE,
			Message: "Something went wrong. Please try again",
		})
		return
	}

	ctx.JSON(http.StatusOK, RegisterWebsiteResponse{
		Status:  constants.GENERIC_SUCCESS_RESPONSE,
		Message: "Website registered successfully.",
	})
}

func (b *BaseController) TestWebsiteLiveliness(ctx *gin.Context) {
	var (
		request RegisterWebsiteRequest
	)

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		logger.Error("error in binding request | err: ", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, WebsiteLivelinessResponse{
			Status:  constants.GENERIC_FAILURE_RESPONSE,
			Message: "Something went wrong. Please try again",
		})
		return
	}

	if request.WebsiteURL == "" {
		logger.Error("invalid request")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, WebsiteLivelinessResponse{
			Status:  constants.GENERIC_FAILURE_RESPONSE,
			Message: "Please enter valid details",
		})
		return
	}

	//check if webiste is live or not before registering it
	isLive, _, err := b.IsWebsiteLive(request.WebsiteURL)
	if err != nil {
		logger.Error("error in testing website liveliness | err: ", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, WebsiteLivelinessResponse{
			Status:  constants.GENERIC_FAILURE_RESPONSE,
			Message: "Something went wrong. Please try again",
		})
		return
	}
	if !isLive {
		logger.Errorf("website:%s is not live!", request.WebsiteURL)
		ctx.AbortWithStatusJSON(http.StatusConflict, WebsiteLivelinessResponse{
			Status:  constants.GENERIC_FAILURE_RESPONSE,
			Message: "Entered webiste is not live",
		})
		return
	}

	ctx.JSON(http.StatusOK, WebsiteLivelinessResponse{
		Status:  constants.GENERIC_SUCCESS_RESPONSE,
		Message: "Website is live!",
	})
}

func (b *BaseController) IsWebsiteLive(url string) (bool, int, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false, 0, err
	}
	defer resp.Body.Close()

	// Treat 2xx and 3xx as live
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return true, resp.StatusCode, nil
	}

	return false, resp.StatusCode, nil
}
