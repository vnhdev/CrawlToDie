package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
)

// create struct to store data
type ProductPrice struct {
	Compareatpricemax int    `json:"compareatpricemax" gorm:"column:compare_price_max"`
	Compareatpricemin int    `json:"compareatpricemin" gorm:"column:compare_price_min"`
	Compareatprice    int    `json:"compareatprice" gorm:"column:compare_price"`
	ProductID         int    `json:"id" gorm:"column:product_id"`
	Price             int    `json:"price" gorm:"column:price"`
	Pricemax          int    `json:"pricemax" gorm:"column:price_max"`
	Pricemin          int    `json:"pricemin" gorm:"column:price_min"`
	Title             string `json:"title" gorm:"column:title"`
	Type              string `json:"type" gorm:"column:type"`
	Vendor            string `json:"vendor" gorm:"column:vendor"`
}

func main() {
	//connected postgres
	dsn := "host=localhost user=postgres password=secret dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// add filepath to crawl
	filepath := "/home/thinkpad/Desktop/Code/testCrawl/Product-Details.Product-Information.json"
	//start to crawl
	err := crawlJSONAndStoreDB(filepath, db)
	if err != nil {
		fmt.Println("Error while crawling JSON and storing in DB:", err)
		return
	}
	fmt.Println("Data inserted successfully!")

}

// crawl file json to store data
func crawlJSONAndStoreDB(filepath string, db *gorm.DB) error {
	//read the json file
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	//parse json file
	var data []ProductPrice
	err = json.Unmarshal([]byte(byteValue), &data)
	if err != nil {
		return err
	}

	for i := 1; i < len(data); i++ {
		data[i].Price /= 100
		data[i].Pricemin /= 100
		data[i].Pricemax /= 100
		data[i].Compareatpricemax /= 100
		data[i].Compareatpricemin /= 100
		data[i].Compareatprice /= 100

	}

	for _, item := range data {
		err := db.Table("product_details_price").Create(&item).Error
		if err != nil {
			return err
		}
	}
	return nil
}
