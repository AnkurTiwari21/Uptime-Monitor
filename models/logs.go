package models

import (
	"context"
	"time"

	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"gorm.io/gorm"
)

type HealthStatus string

const (
	Healthy   HealthStatus = "HEALTHY"
	Unhealthy HealthStatus = "UNHEALTHY"
)

type Log struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"index:idx_website_created_at,sort:desc" json:"created_at"`

	WebsiteId    uint   `gorm:"not null;index:idx_website_created_at" json:"website_id"`
	StatusCode   uint   `gorm:"not null" json:"status_code"`
	LatencyInMS  uint   `gorm:"not null" json:"latency_in_ms"`
	HealthStatus string `gorm:"not null" json:"health_status"`
}

type logsRepo struct {
	db *gorm.DB
}

func (lr *logsRepo) Create(ctx context.Context, log Log) error {
	err := lr.db.WithContext(ctx).Create(&log).Error
	if err != nil {
		logger.Error("error in creating log entry | err: ", err)
		return err
	}
	return nil
}

func (lr *logsRepo) FetchPastRecordStatusByWebsiteID(ctx context.Context, limit uint, webisteID uint) ([]string, error) {
	var (
		statusLogs []string
	)
	err := lr.db.WithContext(ctx).Raw(`
	SELECT health_status FROM logs
	WHERE website_id = $1 
	ORDER BY created_at DESC
	LIMIT $2
	`, webisteID, limit).Scan(&statusLogs).Error
	if err != nil {
		logger.Error("error in fetching past logs | err", err)
		return nil, err
	}
	return statusLogs, err
}
