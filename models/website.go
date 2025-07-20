package models

import (
	"context"
	"time"

	"github.com/ankur12345678/uptime-monitor/pkg/constants"
	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"github.com/ankur12345678/uptime-monitor/utils"
	"gorm.io/gorm"
)

type Website struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UUID          string    `gorm:"unique;not null;" json:"uuid"`
	WebsiteURL    string    `gorm:"not null" json:"website_url"`
	UserId        uint      `gorm:"not null" json:"user_id"`
	LastCheckedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_last_checked_at" json:"last_checked_at"`

	User User `gorm:"foreignKey:UserId;References:ID"`
}

type websiteRepo struct {
	db *gorm.DB
}

func (w *Website) BeforeCreate(tx *gorm.DB) error {
	w.UUID = utils.UUIDGen(constants.WEBISTE_TYPE)
	return nil
}

// Create implements IWebsite.
func (wr *websiteRepo) Create(w *Website) error {
	return wr.CreateWithTx(wr.db, w)
}

// CreateWithTx implements IWebsite.
func (wr *websiteRepo) CreateWithTx(tx *gorm.DB, w *Website) error {
	err := tx.Model(&Website{}).Create(&w).Error
	return err
}

// GetWithTx implements IWebsite.
func (wr *websiteRepo) GetWithTx(where *Website, tx *gorm.DB) (*Website, error) {
	var website Website
	err := tx.Model(&Website{}).Where(where).First(&website).Error
	return &website, err
}

// Update implements IWebsite.
func (wr *websiteRepo) Update(where *Website, w *Website) error {
	return wr.UpdateWithTx(wr.db, where, w)
}

// UpdateWithTx implements IWebsite.
func (wr *websiteRepo) UpdateWithTx(tx *gorm.DB, where *Website, w *Website) error {
	err := tx.
		Model(&Website{}).
		Where(where).Updates(&w).Error
	if err != nil {
		logger.Error("unable to update website | err: ", err)
		return err
	}
	return nil
}

// Delete implements IWebsite.
func (wr *websiteRepo) Delete(where *Website) error {
	return wr.DeleteWithTx(wr.db, where)
}

// DeleteWithTx implements IWebsite.
func (wr *websiteRepo) DeleteWithTx(tx *gorm.DB, where *Website) error {
	err := wr.db.Model(&Website{}).
		Where(where).
		Delete(&Website{}).Error
	if err != nil {
		logger.Error("error in deleting user | err: ", err)
		return err
	}
	return nil
}

func (wr *websiteRepo) FetchWebsitesInBulk(ctx context.Context, limit int) ([]Website, *gorm.DB, error) {
	var websites []Website

	tx := wr.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, nil, tx.Error
	}

	err := tx.WithContext(ctx).Raw(`
		SELECT * FROM websites
		WHERE last_checked_at <= $1 AND deleted_at is NULL
		ORDER BY last_checked_at
		LIMIT $2
		FOR UPDATE SKIP LOCKED
	`, time.Now().Add(-3*time.Minute), limit).Scan(&websites).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	return websites, tx, nil
}
