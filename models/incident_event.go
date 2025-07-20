package models

import (
	"time"

	"github.com/ankur12345678/uptime-monitor/pkg/constants"
	"github.com/ankur12345678/uptime-monitor/utils"
	"gorm.io/gorm"
)

type EventStatus string

const (
	EventStatusPending   EventStatus = "PENDING"
	EventStatusFailed    EventStatus = "FAILED"
	EventStatusDelivered EventStatus = "DELIVERED"
)

type IncidentEvent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UUID         string      `gorm:"not null;index:idx_incident_event_uuid" json:"uuid"`
	HealthStatus string      `gorm:"not null" json:"health_status"`
	WebsiteURL   string      `gorm:"not null" json:"website_url"`
	EventStatus  EventStatus `gorm:"not null" json:"event_status"`

	AlertTargetId uint `gorm:"not null" json:"alert_target_id"`

	AlertTarget AlertTarget `gorm:"foreignKey:AlertTargetId;references:ID" json:"alert_target"`
}

type incidentEventsRepo struct {
	db *gorm.DB
}

func (w *IncidentEvent) BeforeCreate(tx *gorm.DB) error {
	w.UUID = utils.UUIDGen(constants.INCIDENT_EVENT_TYPE)
	return nil
}

func (r *incidentEventsRepo) CreateWithTx(tx *gorm.DB, i *IncidentEvent) error {
	err := tx.Model(&IncidentEvent{}).Create(&i).Error
	return err
}

func (r *incidentEventsRepo) GetWithTx(tx *gorm.DB, where *IncidentEvent) (*IncidentEvent, error) {
	var incidentEvent IncidentEvent
	err := tx.Model(&IncidentEvent{}).Where(where).First(&incidentEvent).Error
	return &incidentEvent, err
}

func (ir *incidentEventsRepo) UpdateWithTx(tx *gorm.DB, where *IncidentEvent, updates *IncidentEvent) error {
	updateMap := map[string]interface{}{
		"event_status":  updates.EventStatus,
		"health_status": updates.HealthStatus,
	}

	return tx.Model(&IncidentEvent{}).
		Where("uuid = ?", where.UUID). // Make sure you're targeting the right row
		Updates(updateMap).Error
}
