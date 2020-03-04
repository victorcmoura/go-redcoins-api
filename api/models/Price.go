package models

import (
	"errors"
	"time"
	"math"
	"net/http"
	"log"
	"fmt"
	"os"
	"net/url"
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type Price struct {
	ID uint64 `gorm:"primary_key;auto_increment" json:"id"`
	BCToUSD float64 `gorm:"not null" json:"bc_usd"`
	USDToBC float64 `gorm:"not null" json:"usd_bc"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (price *Price) Prepare() {
	price.ID = 0
	price.CreatedAt = time.Now()
	price.UpdatedAt = time.Now()
}

func (price *Price) Validate() error {
	diff := math.Abs((price.BCToUSD * price.USDToBC) - 1)

	fmt.Println(diff)

	if math.Abs((price.BCToUSD * price.USDToBC) - 1) > 0.001 {
		return errors.New("Incompatible conversions. Values should compensate each other")
	}
	return nil
}

func (price *Price) SavePrice(db *gorm.DB) (*Price, error) {
	var err error
	err = db.Debug().Create(&price).Error
	if err != nil {
		return &Price{}, err
	}
	return price, nil
}

func (price *Price) GetPrice(db *gorm.DB) (*Price, error) {
	var err error

	err = db.Last(&Price{}).Take(price).Error

	last_update := price.CreatedAt

	if gorm.IsRecordNotFoundError(err) || (err == nil && last_update.Add(time.Minute * 30).After(time.Now())){
		return &Price{}, errors.New("No valid price was found")
	}

	return price, err
}

func UpdatePrice(db *gorm.DB) {
	client := &http.Client{}
	req, err := http.NewRequest("GET","https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := url.Values{}
	q.Add("convert", "USD")
	q.Add("id", "1")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", os.Getenv("EXTERNAL_API_TOKEN"))
	req.URL.RawQuery = q.Encode()


	resp, err := client.Do(req);
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	fmt.Println(resp.Status);

	var respMap map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&respMap)

	if respMap["data"] == nil {
		fmt.Println("Coin Marked API response came empty, check access token");
	}

	data := respMap["data"].(map[string]interface{})
	bitcoin := data["1"].(map[string]interface{})
	quote := bitcoin["quote"].(map[string]interface{})
	usd := quote["USD"].(map[string]interface{})
	usd_price := usd["price"].(float64)
	
	var price Price
	price.BCToUSD = usd_price
	price.USDToBC = 1.0/usd_price

	err_db := db.Debug().Model(&Price{}).Create(&price).Error

	if err_db != nil {
		fmt.Println("Error while updating price: %v", err_db);
	}else{
		fmt.Println("Price updated");
	}
}