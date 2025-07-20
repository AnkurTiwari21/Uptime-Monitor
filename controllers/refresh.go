package controllers

import (
	"net/http"
	"time"

	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"github.com/ankur12345678/uptime-monitor/utils"
	"github.com/gin-gonic/gin"
)

func (base *BaseController) HandleRefresh(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		logger.Error("error setting user in context")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error_message": "Something went wrong. Please try again",
		})
		return
	}
	jti, exists := c.Get("jti")
	if !exists {
		logger.Error("error setting user in context")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error_message": "Something went wrong. Please try again",
		})
		return
	}
	emailStr, ok := email.(string)
	if !ok {
		logger.Error("error in type converison")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error_message": "Something went wrong. Please try again",
		})
		return
	}
	jtiStr, ok := jti.(string)
	if !ok {
		logger.Error("error in type converison")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error_message": "Something went wrong. Please try again",
		})
		return
	}
	accessToken, err := utils.GenerateJWT(Ctrl.Config.JwtSecret, emailStr, Ctrl.Config.JwtExpiryTime)
	if err != nil {
		logger.Error("error in generating access token. Try again")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error_message": "Something went wrong. Please try again",
		})
		return
	}
	prevAuthToken, _ := c.Get("accessToken")
	//blacklisting prev token with key as jti of prev token
	Ctrl.RedisClient.Set(Ctrl.RedisClient.Context(), jtiStr, prevAuthToken.(string), time.Second*600)
	c.JSON(200, gin.H{
		"access_token": accessToken,
		"expiresIn":    600,
	})

}
