package models

import (
	"time"

	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"gorm.io/gorm"
)

type Incident struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	WebsiteId    uint   `gorm:"not null;index:idx_website_id" json:"website_id"`
	HealthStatus string `gorm:"not null" json:"health_status"`
}

type incidentsRepo struct {
	db *gorm.DB
}

func (ir *incidentsRepo) Create(tx *gorm.DB, log Incident) error {
	err := tx.Create(&log).Error
	if err != nil {
		logger.Error("error in creating log entry | err: ", err)
		return err
	}
	return nil
}

func (ir *incidentsRepo) GetWithTx(tx *gorm.DB, where *Incident) (*Incident, error) {
	var incident Incident
	err := tx.Model(&Incident{}).Where(where).First(&incident).Error
	return &incident, err
}

func (ir *incidentsRepo) DeleteWithTx(tx *gorm.DB, where *Incident) error {
	err := ir.db.Model(&Incident{}).
		Where(where).
		Delete(&Incident{}).Error
	if err != nil {
		logger.Error("error in deleting Incident | err: ", err)
		return err
	}
	return nil
}
