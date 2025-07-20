package models

import (
	"time"

	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"gorm.io/gorm"
)


type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserUUID       string `gorm:"unique;not null;" json:"user_uuid"`
	FirstName      string `gorm:"not null" json:"first_name"`
	LastName       string `gorm:"not null" json:"last_name"`
	Email          string `gorm:"unique" json:"email"`
	ProfilePicture string `json:"profile_picture"`
	Password       string `gorm:"not null" json:"omit"`
}

type userRepo struct {
	db *gorm.DB
}

// Create implements IUser.
func (r *userRepo) Create(u *User) error {
	return r.CreateWithTx(r.db, u)
}

// CreateWithTx implements IUser.
func (r *userRepo) CreateWithTx(tx *gorm.DB, u *User) error {
	err := tx.Model(&User{}).Create(&u).Error
	return err
}

// GetWithTx implements IUser.
func (r *userRepo) GetWithTx(where *User, tx *gorm.DB) (*User, error) {
	var user User
	err := tx.Model(&User{}).Where(where).First(&user).Error
	return &user, err
}

// Update implements IUser.
func(r *userRepo) Update(where *User, u *User) error {
	return r.UpdateWithTx(r.db, where, u)
}

// UpdateWithTx implements IUser.
func(r *userRepo) UpdateWithTx(tx *gorm.DB, where *User, u *User) error {
	err := tx.
		Model(&User{}).
		Where(where).Updates(&u).Error
	if err != nil {
		logger.Error("unable to update user | err: ", err)
		return err
	}
	return nil
}

// Delete implements IUser.
func (r *userRepo) Delete(where *User) error {
	return r.DeleteWithTx(r.db, where)
}

// DeleteWithTx implements IUser.
func(r *userRepo) DeleteWithTx(tx *gorm.DB, where *User) error {
	err := r.db.Model(&User{}).
		Where(where).
		Delete(&User{}).Error
	if err != nil {
		logger.Error("error in deleting user | err: ", err)
		return err
	}
	return nil
}


func (r *userRepo) GetById(UserUUID string) (*User, error) {
	return r.GetWithTx(&User{UserUUID: UserUUID}, r.db)
}

func (r *userRepo) GetByEmail(Email string) (*User, error) {
	return r.GetWithTx(&User{Email: Email}, r.db)
}
