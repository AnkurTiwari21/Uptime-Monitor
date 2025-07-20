package migration

import (
	"context"
	"fmt"

	config "github.com/ankur12345678/uptime-monitor/Config"
	models "github.com/ankur12345678/uptime-monitor/Models"
	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Creds) *gorm.DB {
	config := cfg
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=%s TimeZone=%s", config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DbPort, config.DbSslMode, config.DbTimezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error connecting DB...Exiting!")
	}
	db.AutoMigrate(&models.User{}, &models.Website{}, &models.AlertConfig{}, &models.Log{}, &models.Incident{}, &models.AlertTarget{}, &models.IncidentEvent{})
	logger.Info("Connected to DB!")
	return db
}

func InitRedisClient(c *config.Creds) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     c.RedisConnectionAddress,
		Password: c.RedisConnectionPassword, // no password set
		DB:       0,                         // use default DB
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic("REDIS CONNECTION: FAILED...")
	}
	logger.Info("REDIS CONNECTION: SUCCESS...")
	return client
}
