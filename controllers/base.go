package controllers

import (
	"errors"

	config "github.com/ankur12345678/uptime-monitor/Config"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type BaseController struct {
	DB          *gorm.DB
	Config      *config.Creds
	RedisClient *redis.Client
	Router      *gin.Engine
	Validator   *validator.Validate
	Translator  *ut.Translator
}

var Ctrl BaseController //GLOBAL instance for controllers only

func GetEmailFromContext(ctx *gin.Context) (string, error) {
	email, exists := ctx.Get("email")
	if !exists {
		return "", errors.New("email does not exists")
	}

	return email.(string), nil
}
