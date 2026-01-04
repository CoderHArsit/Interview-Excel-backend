package models

import (
	"time"

	"gorm.io/gorm"
)

type WalletTransaction struct {
	ID       uint `gorm:"primaryKey"`
	WalletID uint `gorm:"index;not null"`

	Wallet        *Wallet `gorm:"foreignKey:WalletID;references:ID" json:"wallet,omitempty"`
	AmountInPaise int64   `gorm:"not null"` // +credit / -debit
	Type          string  `gorm:"index"`    // credit | debit
	Source        string  `gorm:"index"`    // session | refund | payout
	ReferenceID   string  `gorm:"index"`    // session_uuid / payout_id

	Description string

	CreatedAt time.Time
}

type walletTransactionRepo struct {
	DB *gorm.DB
}

func (r *walletTransactionRepo) Create(tx *gorm.DB, wt *WalletTransaction) error {
	return tx.Create(wt).Error
}

func (r *walletTransactionRepo) GetByWalletID(walletID uint) ([]WalletTransaction, error) {
	var transactions []WalletTransaction
	err := r.DB.Where("wallet_id = ?", walletID).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *walletTransactionRepo) GetByReferenceID(referenceID string) ([]WalletTransaction, error) {
	var transactions []WalletTransaction
	err := r.DB.Where("reference_id = ?", referenceID).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
