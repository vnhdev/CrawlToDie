package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
)

type JsonCrawlCategories struct {
	Page        int `json:"page"`
	CurrentPage int `json:"current_page"`
	ListCol     []struct {
		Title       string `json:"title" gorm:"column:title"`
		HandleTitle string `json:"handle_title" gorm:"column:handle_title"`
		ProductID   int    `json:"id" gorm:"column:product_id"`
		Gender      string `json:"gender" gorm:"column:gender"`
	} `json:"list_col"`
}

func main() {
	//Connect to client
	client := &http.Client{}
	//Connecting to database
	dsn := "host=localhost user=postgres password=secret dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("can't connect to db")
	}
	//Get all data from json
	for page := 1; page <= 30; page++ {
		url := fmt.Sprintf("https://www.maisononline.vn/collections?page=%d", page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Lỗi khi tạo request:", err)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Lỗi khi thực hiện request:", err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Lỗi khi đọc dữ liệu từ response:", err)
			return
		}

		//parse data json unmarshal
		var data JsonCrawlCategories
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("Lỗi khi phân tích dữ liệu:", err)
			return
		}

		//store data in database
		var categories JsonCrawlCategories
		err = json.Unmarshal([]byte(body), &categories)
		if err != nil {
			fmt.Println("error parsing json", err)
			return
		}

		for _, item := range categories.ListCol {
			db.Table("categories").CreateInBatches(item, 6)
		}
	}
	fmt.Println("Get collection successful")
}
