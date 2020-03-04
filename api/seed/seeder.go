package seed

import (
	"github.com/jinzhu/gorm"
	"github.com/victorcmoura/go-redcoins-api/api/models"
)

func Load(db *gorm.DB) {
	models.UpdatePrice(db)
}