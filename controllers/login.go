package controllers

import (
	"net/http"

	models "github.com/ankur12345678/uptime-monitor/Models"
	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"github.com/ankur12345678/uptime-monitor/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (base *BaseController) LoginHandler(c *gin.Context) {
	var (
		request  = SignInRequest{}
		userRepo = models.InitUserRepo(Ctrl.DB)
	)

	err := c.ShouldBindJSON(&request)
	if err != nil {
		logger.Error("error in binding request | err: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Please check details and try again",
		})
		return
	}

	//check if user present in the db
	user, err := userRepo.GetByEmail(request.Email)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{
			"error_message": "No records found, please checkout signup!",
		})
		return
	}
	if err != nil {
		logger.Error("error in getting user from DB | err: ", err)
		c.JSON(http.StatusOK, gin.H{
			"error_message": "Something went wrong. Please try again",
		})
		return
	}

	//verify password
	match := utils.VerifyPassword(request.Password, user.Password)
	if !match {
		c.JSON(http.StatusOK, gin.H{
			"error_message": "Invalid password!",
		})
		return
	}

	//generate a jwt with time 10 min
	accessToken, err := utils.GenerateJWT(Ctrl.Config.JwtSecret, request.Email, Ctrl.Config.JwtExpiryTime)
	if err != nil {
		logger.Error("error generating access token | err: ", err)
		c.JSON(http.StatusOK, gin.H{
			"error_message": "Something went wrong. Please try again",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"expires_in":   Ctrl.Config.JwtExpiryTime,
	})
}
