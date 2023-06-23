package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Product struct {
	ID              int    `bson:"id" json:"product_id" gorm:"column:product_id"`
	Vendor          string `bson:"vendor" json:"vendor" gorm:"column:vendor"`
	Type            string `bson:"type" json:"type" gorm:"column:type"`
	Title           string `bson:"title" json:"title" gorm:"column:title"`
	Price           int    `bson:"price" json:"price" gorm:"column:price"`
	PriceMax        int    `bson:"pricemax" json:"pricemax" gorm:"column:price_max"`
	PriceMin        int    `bson:"pricemin" json:"pricemin" gorm:"column:price_min"`
	ComparePrice    int    `bson:"compareatprice" json:"compareatprice" gorm:"compare_price"`
	ComparePriceMax int    `bson:"compareatpricemax" json:"compareatpricemax" gorm:"column:compare_price_max"`
	ComparePriceMin int    `bson:"compareatpricemin" json:"compareatpricemin" gorm:"column:compare_price_min"`
}

func main() {
	//connecting to gorm
	dsn := "host=localhost user=postgres password=secret dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//connecting to mongodb
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27023"))
	if err != nil {
		panic(err)
	}
	// connect to mongodb database and collection name
	userCollection := client.Database("Product-Details").Collection("Product-Information")

	//create a filter to search
	filter := bson.M{}
	//create a cursor
	cur, err := userCollection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())

	var products []Product
	for cur.Next(context.Background()) {
		var product Product
		err := cur.Decode(&product)
		if err != nil {
			log.Fatal(err)
		}
		products = append(products, product)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(products); i++ {
		products[i].Price /= 100
		products[i].PriceMin /= 100
		products[i].PriceMax /= 100
		products[i].ComparePriceMax /= 100
		products[i].ComparePriceMin /= 100
		products[i].ComparePrice /= 100
	}

	for _, item := range products {
		err := db.Table("product_details_price").Create(&item).Error
		if err != nil {
			return
		}
	}
}
