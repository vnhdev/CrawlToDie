package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type ProductItem struct {
	Vendor string `json:"vendor" gorm:"column:vendor"`
	Type   string `json:"type" gorm:"column:type"`
	Title  string `json:"title" gorm:"column:title"`
	Price  int    `json:"price" gorm:"column:price"`
}

type ProductPriceHistory struct {
	Vendor       string    `json:"vendor" gorm:"column:vendor"`
	Type         string    `json:"type" gorm:"column:type"`
	Title        string    `json:"title" gorm:"column:title"`
	Price        int       `json:"price" gorm:"column:price"`
	ComparePrice int       `json:"compare_price" gorm:"column:compare_price"`
	CurrentDay   time.Time `json:"current_day" gorm:"column:current_day"`
}

func main() {
	dsn := "host=localhost user=postgres password=secret dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("can't connect to database")
	}
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.GET("get-product-info", GetProductInfo(db))
		v1.GET("get-product-change-price", GetHistoryChange(db))
	}
	router.Run()
}

func GetProductInfo(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data []ProductItem
		if err := db.Table("product_details_price").Select("vendor,type,title,price").Find(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": data})
	}
}

func GetHistoryChange(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data []ProductPriceHistory
		if err := db.Table("product_details_price").Select("vendor,type,title,price,compare_price,current_day").Where("compare_price != ?", 0).Find(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": data})
	}
}
