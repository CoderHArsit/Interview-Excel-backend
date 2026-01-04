package models

import (
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	ID       uint   `gorm:"primaryKey"`
	UserUUID string `gorm:"uniqueIndex;not null"` // expert user

	User *User `gorm:"foreignKey:UserUUID;references:UserUUID" json:"user,omitempty"`

	BalanceInPaise int64 `gorm:"not null;default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type walletRepo struct {
	DB *gorm.DB
}


func (r *walletRepo) GetByUserUUID(userUUID string) (*Wallet, error) {
	var wallet Wallet
	err := r.DB.Where("user_uuid = ?", userUUID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepo) Create(wallet *Wallet) error {
	return r.DB.Create(wallet).Error
}

func (r *walletRepo) UpdateBalance(userUUID string, newBalanceInPaise int64) error {
	return r.DB.Model(&Wallet{}).Where("user_uuid = ?", userUUID).
		Update("balance_in_paise", newBalanceInPaise).Error
}
