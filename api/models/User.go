package models

import(
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"size:255;not null;unique" json:"name"`
	Email string `gorm:"size:255;not null;unique" json:"email"`
	BirthDate time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"birth_date"`
	Password string `gorm:"size:100;not null;" json:"password`
	BCBalance uint64 `gorm:"default:0;not null" json:"bc_balance"`
	USDBalance uint64 `gorm:"default:0;not null" json:"usd_balance"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (user *User) BeforeSave() error {
	hashedPassword, err := Hash(user.Password)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return nil
}

func (user *User) Prepare() {
	user.ID = 0
	user.Name = strings.TrimSpace(user.Name)
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.BCBalance = 0
	user.USDBalance = 0
}

func (user *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if user.Password == "" {
			return errors.New("Password required")
		}
		if user.Email == "" {
			return errors.New("Email required")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid email")
		}
		if user.BirthDate.After(time.Now()) {
			return errors.New("Invalid birthdate")
		}
		return nil

	case "login":
		if user.Password == "" {
			return errors.New("Password required")
		}
		if user.Email == "" {
			return errors.New("Email required")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid email")
		}
		return nil
	
	default:
		if user.Password == "" {
			return errors.New("Password required")
		}
		if user.Email == "" {
			return errors.New("Email required")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid email")
		}
		return nil
	}
}

func (user *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Find(&users).Error
	if err != nil{
		return &[]User{}, err
	}
	return &users, err
}

func (user *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User not found")
	}
	return user, err
}

func (user *User) UpdateAUser (db *gorm.DB, uid uint32) (*User, error) {
	err := user.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(user).UpdateColumns(
		map[string]interface{}{
			"name": user.Name,
			"password": user.Password,
			"email": user.Email,
			"birth_date": user.BirthDate,
			"bc_balance": user.BCBalance,
			"usd_balance": user.USDBalance,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}

	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&user).Error
	
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}