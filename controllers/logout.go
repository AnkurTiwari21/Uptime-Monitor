package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
)

func (base *BaseController) HandleLogOut(c *gin.Context) {
	email, _ := c.Get("email")

	prevAuthToken, _ := c.Get("accessToken")
	//blacklisting prev token with key as email of the user
	Ctrl.RedisClient.Set(Ctrl.RedisClient.Context(), email.(string), prevAuthToken.(string), time.Second*600)

	c.JSON(200, gin.H{
		"message": "Thanks for using our website, see you again!",
	})
}
