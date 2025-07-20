package models

import (
	"time"

	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"gorm.io/gorm"
)

type TargetType string

const (
	TargetTypeSMS   TargetType = "sms"
	TargetTypeEmail TargetType = "email"
)

type AlertTarget struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	TargetType  TargetType `gorm:"not null" json:"target_type"`
	TargetValue string     `gorm:"not null" json:"target_value"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`

	AlertConfigID uint `gorm:"not null;index" json:"alert_config_id"`
}

type alertTargetRepo struct {
	db *gorm.DB
}

// Create implements IAlertTarget.
func (atr *alertTargetRepo) Create(at *AlertTarget) error {
	return atr.CreateWithTx(atr.db, at)
}

// CreateWithTx implements IAlertTarget.
func (atr *alertTargetRepo) CreateWithTx(tx *gorm.DB, at *AlertTarget) error {
	return tx.Model(&AlertTarget{}).Create(&atr).Error
}

// GetWithTx implements IAlertTarget.
func (atr *alertTargetRepo) GetWithTx(tx *gorm.DB, where *AlertTarget) (*AlertTarget, error) {
	var at AlertTarget
	err := tx.Model(&AlertTarget{}).Where(where).First(&at).Error
	return &at, err
}

// Update implements IAlertTarget.
func (atr *alertTargetRepo) Update(where *AlertTarget, a *AlertTarget) error {
	return atr.UpdateWithTx(atr.db, where, a)
}

// UpdateWithTx implements IAlertTarget.
func (atr *alertTargetRepo) UpdateWithTx(tx *gorm.DB, where *AlertTarget, a *AlertTarget) error {
	err := tx.
		Model(&AlertTarget{}).
		Where(where).Updates(&a).Error
	if err != nil {
		logger.Error("unable to update alert config | err: ", err)
		return err
	}
	return nil
}

// Delete implements IAlertTarget.
func (atr *alertTargetRepo) Delete(where *AlertTarget) error {
	return atr.DeleteWithTx(atr.db, where)
}

// DeleteWithTx implements IAlertTarget.
func (atr *alertTargetRepo) DeleteWithTx(tx *gorm.DB, where *AlertTarget) error {
	err := atr.db.Model(&AlertTarget{}).
		Where(where).
		Delete(&AlertTarget{}).Error
	if err != nil {
		logger.Error("error in deleting alert config | err: ", err)
		return err
	}
	return nil
}

func (atr *alertTargetRepo) GetAllByAlertConfigID(alertConfigID uint) ([]AlertTarget, error) {
	var targets []AlertTarget
	err := atr.db.
		Model(&AlertTarget{}).
		Where("alert_config_id = ? AND is_active = true", alertConfigID).
		Find(&targets).Error

	if err != nil {
		logger.Error("error in fetching alert targets by alert config id | err: ", err)
		return nil, err
	}
	return targets, nil
}
