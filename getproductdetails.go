package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"time"
)

// product link struct
type categoriesLink struct {
	ProductLink string `json:"ProductLink" gorm:"column:product_link"`
}

// design database product details
type ProductDetails struct {
	Available            bool        `json:"available"`
	CompareAtPriceMax    int         `json:"compare_at_price_max"`
	CompareAtPriceMin    int         `json:"compare_at_price_min"`
	CompareAtPriceVaries bool        `json:"compare_at_price_varies"`
	CompareAtPrice       int         `json:"compare_at_price"`
	Content              interface{} `json:"content"`
	Description          string      `json:"description"`
	Handle               string      `json:"handle"`
	Id                   int         `json:"id"`
	Media                []struct {
		Alt      interface{} `json:"alt"`
		Id       int         `json:"id"`
		Position int         `json:"position"`
	} `json:"media"`
	Options []struct {
		Name      string   `json:"name"`
		Position  int      `json:"position"`
		ProductId int      `json:"product_id"`
		Values    []string `json:"values"`
	} `json:"options"`
	Price           int      `json:"price"`
	PriceMax        int      `json:"price_max"`
	PriceMin        int      `json:"price_min"`
	PriceVaries     bool     `json:"price_varies"`
	Tags            []string `json:"tags"`
	Title           string   `json:"title"`
	Type            string   `json:"type"`
	Url             string   `json:"url"`
	Pagetitle       string   `json:"pagetitle"`
	Metadescription string   `json:"metadescription"`
	Variants        []struct {
		Id                   int         `json:"id"`
		Barcode              string      `json:"barcode"`
		Available            bool        `json:"available"`
		Price                int         `json:"price"`
		Sku                  string      `json:"sku"`
		Option1              string      `json:"option1"`
		Option2              string      `json:"option2"`
		Option3              string      `json:"option3"`
		Options              []string    `json:"options"`
		InventoryQuantity    int         `json:"inventory_quantity"`
		OldInventoryQuantity int         `json:"old_inventory_quantity"`
		Title                string      `json:"title"`
		Weight               int         `json:"weight"`
		CompareAtPrice       int         `json:"compare_at_price"`
		InventoryManagement  string      `json:"inventory_management"`
		InventoryPolicy      string      `json:"inventory_policy"`
		Selected             bool        `json:"selected"`
		Url                  interface{} `json:"url"`
		FeaturedImage        interface{} `json:"featured_image"`
	} `json:"variants"`
	Vendor            string    `json:"vendor"`
	PublishedAt       time.Time `json:"published_at"`
	CreatedAt         time.Time `json:"created_at"`
	NotAllowPromotion bool      `json:"not_allow_promotion"`
}

func main() {
	// Tạo một client HTTP
	client := &http.Client{}
	// user tje connected client to perform database operations
	baseurl := "https://www.maisononline.vn"
	//Connect to MongoDB database
	clients, err := connectToMongoDB()
	database := clients.Database("Product-Details")

	if err != nil {
		fmt.Println("Failed to connect to mongodb:", err)
		return
	}
	defer clients.Disconnect(context.Background())

	//Connect to Gorm database
	dsn := "host=localhost user=postgres password=secret dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// get all product link and store into struct
	var link []categoriesLink
	if err := db.Table("product_link").Select("product_link").Find(&link).Error; err != nil {
		return
	}
	//parse data inside database and store into data variables
	data := make([]string, len(link))
	for i, itm := range link {
		data[i] = itm.ProductLink
	}

	for _, url := range data {
		jsString := ".js"
		newURL := baseurl + url + jsString
		// Tạo một request GET mới
		req, err := http.NewRequest("GET", newURL, nil)
		if err != nil {
			fmt.Println("Lỗi khi tạo request:", err)
			return
		}

		// Thực hiện request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Lỗi khi thực hiện request:", err)
			return
		}
		defer resp.Body.Close()

		// Đọc dữ liệu từ response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Lỗi khi đọc dữ liệu từ response:", err)
			return
		}
		//store data inside database
		var metaData ProductDetails
		err = json.Unmarshal(body, &metaData)
		if err != nil {
			return
		}

		//insert into database inside mongodb
		err = insertDataProductDetails(database, metaData)
		if err != nil {
			fmt.Println("Failed to insert product details:", err)
			return
		}

	}
}

// function get connect URL
func getConnectionOptions() *options.ClientOptions {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27023")
	return clientOptions
}

// function connect into mongodb
func connectToMongoDB() (*mongo.Client, error) {
	clientOptions := getConnectionOptions()

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// function insert data into database in mongodb
func insertDataProductDetails(database *mongo.Database, product ProductDetails) error {
	collection := database.Collection("Product-Information")
	_, err := collection.InsertOne(context.Background(), product)
	if err != nil {
		return err
	}

	return nil
}
