package controllers

import (
	"net/http"

	models "github.com/ankur12345678/uptime-monitor/Models"
	"github.com/ankur12345678/uptime-monitor/pkg/constants"
	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"github.com/ankur12345678/uptime-monitor/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func (base *BaseController) SignUpHandler(c *gin.Context) {
	var (
		signUpRequest = SignUpRequest{}
		userRepo      = models.InitUserRepo(Ctrl.DB)
	)

	err := c.ShouldBindJSON(&signUpRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Please check the details and try again",
		})
	}

	//TODO: validate request
	err = base.Validator.Struct(&signUpRequest)
	if err != nil {
		validationErrors := []constants.Error{}
		for _, err := range err.(validator.ValidationErrors) {
			logger.Error("error in validating the request | err", err)
			e := constants.Error{}
			e.Description = err.Translate(*base.Translator)
			e.Field = err.Field()
			validationErrors = append(validationErrors, e)
			c.JSON(http.StatusBadRequest,
				map[string]interface{}{
					"errors": validationErrors,
				},
			)
			return
		}
	}

	//check if username exists previously
	_, err = userRepo.GetByEmail(signUpRequest.Email)
	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "user already exists with this email, please try to login",
		})
		return
	}

	user := models.User{
		FirstName: signUpRequest.FirstName,
		LastName:  signUpRequest.LastName,
		Email:     signUpRequest.Email,
	}

	user.UserUUID = utils.UUIDGen(constants.USER_TYPE)
	hashedPassword, err := utils.HashPassword(signUpRequest.Password)
	if err != nil {
		logger.Error("error while hashing password | err: ", err)
		c.JSON(http.StatusOK, gin.H{
			"error_message": "Something went wrong. Please try again",
		})
		return
	}

	user.Password = hashedPassword

	err = userRepo.Create(&user)
	if err != nil {
		logger.Error("error while creating user | err: ", err)
		c.JSON(http.StatusOK, gin.H{
			"error_message": "Something went wrong. Please try again",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User Created! Please visit login endpoint.",
	})
}
