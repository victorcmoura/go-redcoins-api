package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Transaction struct {
	ID uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Type string `gorm:"size:255;not null" json:"type"`
	BCValue uint64 `gorm:"not null" json:"bc_value"`
	USDValue uint64 `gorm:"not null" json:"usd_value"`
	OwnerID uint32 `gorm:"nor null" json:"owner_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (transaction *Transaction) Prepare() {
	transaction.ID = 0
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()
}

func (transaction *Transaction) Validate() error {
	if transaction.Type != "buy" && transaction.Type != "sell" {
		return errors.New("Invalid transaction type (only sell and buy available)")
	}
	return nil
}

func (transaction *Transaction) SaveTransaction(db *gorm.DB) (*Transaction, error) {
	var err error
	err = db.Debug().Model(&Transaction{}).Create(&transaction).Error
	if err != nil {
		return &Transaction{}, err
	}

	var owner User

	if transaction.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", transaction.OwnerID).Take(&owner).Error
		if err != nil {
			return &Transaction{}, err
		}
	}
	return transaction, nil
}

func (transaction *Transaction) FindAllTransactions(db *gorm.DB) (*[]Transaction, error) {
	var err error
	transactions := []Transaction{}
	err = db.Debug().Model(&Transaction{}).Find(&transactions).Error
	if err != nil {
		return &[]Transaction{}, err
	}
	return &transactions, nil
}