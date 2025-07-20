package migration

import (
	"gorm.io/gorm"
)

func SeedDB(DB *gorm.DB) {
	// var (
	// 	userRepo        = models.InitUserRepo(DB)
	// 	websiteRepo     = models.InitWebsiteRepo(DB)
	// 	alertConfigRepo = models.InitAlertConfigRepo(DB)
	// 	alertTargetRepo = models.InitAlertTargetRepo(DB)
	// )

	// user := &models.User{UserUUID: "user_1", FirstName: "Ankur", LastName: "Tiwari", Email: "a@gmail.com", Password: "hehe", ProfilePicture: "yo.svg"}

	// err := userRepo.Create(user)
	// if err != nil {
	// 	logger.Error("err in creating user entry ", err)
	// 	// //return
	// }

	// //make entries for the webiste
	// web1 := &models.Website{UUID: "w_1", WebsiteURL: "google.com", UserId: user.ID}
	// err = websiteRepo.Create(web1)
	// if err != nil {
	// 	logger.Error("err in creating website entry ", err)
	// 	// //return
	// }
	// config1 := &models.AlertConfig{WebsiteID: web1.ID, IsEnabled: true}
	// err = alertConfigRepo.Create(config1)
	// if err != nil {
	// 	logger.Error("err in creating alert config entry ", err)
	// 	////return
	// }
	// err = alertTargetRepo.Create(&models.AlertTarget{TargetType: "email", TargetValue: "ankurtiwari613@gmail.com", IsActive: true, AlertConfigID: config1.ID})
	// if err != nil {
	// 	logger.Error("err in creating alert target entry ", err)
	// 	////return
	// }

	// logger.Info("-----SEEDING SUCCESS-----")
}
