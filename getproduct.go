package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
)

type product struct {
	Brands      string `json:"brands" gorm:"column:brands"`
	ProductName string `json:"product_name" gorm:"column:product_name"`
	ProductLink string `json:"product_link_details" 	gorm:"column:product_link"`
}

type categoriesID struct {
	Id int `json:"id" gorm:"column:product_id"`
}

func main() {
	baseURL := "https://www.maisononline.vn/search?q=filter=((collectionid:product=%s))&view=filter&page="

	//connect to database
	dsn := "host=localhost user=postgres password=secret dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	//range to get element data
	c := colly.NewCollector()
	c.OnHTML("h2.product-loop-name", func(c *colly.HTMLElement) {
		product := product{
			Brands:      c.ChildText("a.pro-vendor"),
			ProductName: c.ChildText("a.pro-title"),
			ProductLink: c.ChildAttr("a.pro-title", "href"),
		}
		db.Table("product_link").CreateInBatches(product, 3)
	})

	// get product id from categories table in database
	var item []categoriesID
	if err := db.Table("categories").Select("product_id").Find(&item).Error; err != nil {
		return
	}

	//parse data in slice categoriesID into data to use
	data := make([]int, len(item))
	for i, itm := range item {
		data[i] = itm.Id
	}

	//add collectionid and page tp visit in go colly1
	for i, _ := range data {
		newData := strconv.Itoa(data[i])
		url := fmt.Sprintf(baseURL, newData)
		for i := 1; i <= 10; i++ {
			newUrlID := strconv.Itoa(i)
			c.Visit(url + newUrlID)
		}
	}
	fmt.Println("Get product successful")
}
