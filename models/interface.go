package models

import (
	"context"

	"gorm.io/gorm"
)

type IUser interface {
	GetById(UserUUID string) (*User, error)
	GetByEmail(Email string) (*User, error)
	GetWithTx(where *User, tx *gorm.DB) (*User, error)
	Create(u *User) error
	CreateWithTx(tx *gorm.DB, u *User) error
	Update(where *User, u *User) error
	UpdateWithTx(tx *gorm.DB, where *User, u *User) error
	Delete(where *User) error
	DeleteWithTx(tx *gorm.DB, where *User) error
}

type IWebiste interface {
	Create(w *Website) error
	CreateWithTx(tx *gorm.DB, w *Website) error
	GetWithTx(where *Website, tx *gorm.DB) (*Website, error)
	Update(where *Website, w *Website) error
	UpdateWithTx(tx *gorm.DB, where *Website, w *Website) error
	Delete(where *Website) error
	DeleteWithTx(tx *gorm.DB, where *Website) error
	FetchWebsitesInBulk(ctx context.Context, limit int) ([]Website, *gorm.DB, error)
}

type IAlertConfig interface {
	Create(ac *AlertConfig) error
	CreateWithTx(tx *gorm.DB, ac *AlertConfig) error
	GetWithTx(tx *gorm.DB, where *AlertConfig) (*AlertConfig, error)
	Update(where *AlertConfig, a *AlertConfig) error
	UpdateWithTx(tx *gorm.DB, where *AlertConfig, a *AlertConfig) error
	Delete(where *AlertConfig) error
	DeleteWithTx(tx *gorm.DB, where *AlertConfig) error
}

type IAlertTarget interface {
	Create(at *AlertTarget) error
	CreateWithTx(tx *gorm.DB, at *AlertTarget) error
	GetWithTx(tx *gorm.DB, where *AlertTarget) (*AlertTarget, error)
	Update(where *AlertTarget, a *AlertTarget) error
	UpdateWithTx(tx *gorm.DB, where *AlertTarget, a *AlertTarget) error
	Delete(where *AlertTarget) error
	DeleteWithTx(tx *gorm.DB, where *AlertTarget) error
	GetAllByAlertConfigID(alertConfigID uint) ([]AlertTarget, error)
}

type ILog interface {
	Create(ctx context.Context, log Log) error
	FetchPastRecordStatusByWebsiteID(ctx context.Context, limit uint, webisteID uint) ([]string, error)
}

type IIncident interface {
	Create(tx *gorm.DB, log Incident) error
	GetWithTx(tx *gorm.DB, where *Incident) (*Incident, error)
	DeleteWithTx(tx *gorm.DB, where *Incident) error
}

type IIncidentEvent interface {
	CreateWithTx(tx *gorm.DB, i *IncidentEvent) error
	GetWithTx(tx *gorm.DB, where *IncidentEvent) (*IncidentEvent, error)
	UpdateWithTx(tx *gorm.DB, where *IncidentEvent, i *IncidentEvent) error
}
