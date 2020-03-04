package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Transaction struct {
	ID uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Type string `gorm:"size:255;not null" json:"type"`
	BCValue float64 `gorm:"not null" json:"bc_value"`
	USDValue float64 `gorm:"not null" json:"usd_value"`
	OwnerID uint32 `gorm:"nor null" json:"owner_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (transaction *Transaction) Prepare(db *gorm.DB) error {
	transaction.ID = 0
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	if transaction.BCValue < 0 || transaction.USDValue < 0 {
		var latest_price Price
		var err error

		err = db.Last(&Price{}).Take(&latest_price).Error

		if err == nil {
			if transaction.USDValue < 0 {
				transaction.USDValue = latest_price.BCToUSD * transaction.BCValue
			}else{
				transaction.BCValue = latest_price.USDToBC * transaction.USDValue
			}	
		}else{
			return err
		}
	}else{
		return errors.New("Invalid currency values. Set one to -1 and the other to the exchanged amount.")
	}

	return nil
}

func (transaction *Transaction) Validate() error {
	if transaction.Type != "buy" && transaction.Type != "sell" {
		return errors.New("Invalid transaction type (only sell and buy available)")
	}

	if transaction.BCValue < 0 && transaction.USDValue < 0 {
		return errors.New("Invalid currency values. Set one to -1 and the other to the exchanged amount.")
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

		if transaction.Type == "buy" {
			owner.BCBalance += transaction.BCValue
			owner.USDBalance += transaction.USDValue
		}else {
			owner.BCBalance -= transaction.BCValue
			owner.USDBalance -= transaction.USDValue
		}
	
		db.Debug().Save(&owner)
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