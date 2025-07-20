package models

import (
	"time"

	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"gorm.io/gorm"
)

type AlertConfig struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	WebsiteID        uint `gorm:"not null;index" json:"website_id"`
	FailureThreshold int  `gorm:"not null;default:3" json:"failure_threshold"`
	LatencyThreshold int  `gorm:"not null;default:5000" json:"latency_threshold"`

	IsEnabled bool `gorm:"default:false" json:"is_enabled"`

	Website Website `gorm:"foreignKey:WebsiteID;References:ID"`
}

type alertConfigRepo struct {
	db *gorm.DB
}

// Create implements IAlertConfig.
func (acr *alertConfigRepo) Create(ac *AlertConfig) error {
	return acr.CreateWithTx(acr.db, ac)
}

// CreateWithTx implements IAlertConfig.
func (acr *alertConfigRepo) CreateWithTx(tx *gorm.DB, ac *AlertConfig) error {
	return tx.Model(&AlertConfig{}).Create(&ac).Error
}

// GetWithTx implements IAlertConfig.
func (acr *alertConfigRepo) GetWithTx(tx *gorm.DB, where *AlertConfig) (*AlertConfig, error) {
	var ac AlertConfig
	err := tx.Model(&AlertConfig{}).Where(where).First(&ac).Error
	return &ac, err
}

// Update implements IAlertConfig.
func (acr *alertConfigRepo) Update(where *AlertConfig, a *AlertConfig) error {
	return acr.UpdateWithTx(acr.db, where, a)
}

// UpdateWithTx implements IAlertConfig.
func (acr *alertConfigRepo) UpdateWithTx(tx *gorm.DB, where *AlertConfig, a *AlertConfig) error {
	err := tx.
		Model(&AlertConfig{}).
		Where(where).Updates(&a).Error
	if err != nil {
		logger.Error("unable to update alert config | err: ", err)
		return err
	}
	return nil
}

// Delete implements IAlertConfig.
func (acr *alertConfigRepo) Delete(where *AlertConfig) error {
	return acr.DeleteWithTx(acr.db, where)
}

// DeleteWithTx implements IAlertConfig.
func (acr *alertConfigRepo) DeleteWithTx(tx *gorm.DB, where *AlertConfig) error {
	err := acr.db.Model(&AlertConfig{}).
		Where(where).
		Delete(&AlertConfig{}).Error
	if err != nil {
		logger.Error("error in deleting alert config | err: ", err)
		return err
	}
	return nil
}
