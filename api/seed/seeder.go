package seed

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/victorcmoura/go-redcoins-api/api/models"
)

var users = []models.User{
	models.User{
		Name: "Victor Moura",
		Email: "victor@email.com",
		Password: "password",
		BirthDate: time.Now().AddDate(-30, 0, 0),
	},
	models.User{
		Name: "Moura Victor",
		Email: "moura@email.com",
		Password: "password",
		BirthDate: time.Now().AddDate(-31, 0, 0),
	},
}

var transactions = []models.Transaction{
	models.Transaction{
		Type: "buy",
		BCValue: 1000000,
		USDValue: 932932,
	},
	models.Transaction{
		Type: "buy",
		BCValue: 1000000,
		USDValue: 932932,
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Transaction{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Transaction{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Transaction{}).AddForeignKey("owner_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		transactions[i].OwnerID = users[i].ID

		err = db.Debug().Model(&models.Transaction{}).Create(&transactions[i]).Error
		if err != nil {
			log.Fatalf("cannot seed transactions table: %v", err)
		}
	}
}