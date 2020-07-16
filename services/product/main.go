package main

import (
	"encoding/csv"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"io"
	"log"
	productController "ogani.com/services/product/controllers/product"
	"os"
	"strconv"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "ogani.com/services/product/docs" // docs is generated by Swag CLI, you have to import it.

	models "ogani.com/services/product/models"
)

func init() {
	db, err := gorm.Open("postgres", "host=product-pg port=5444 user=postgres dbname=ogani_product password=postgres sslmode=disable")
	defer db.Close()
	db.AutoMigrate(&models.ProductType{}, &models.ProductItem{}, &models.ProductBrand{})
	seedData(db)
	if err != nil {
		panic(err)
	}
}

// @title Catalog Swagger API
// @version 1.0
// @description  This is Catalog service
// @termsOfService http://swagger.io/terms/
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
func main() {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx){
		c.Set("Content-Type","application/json")

		c.Next()
	})
	api := app.Group("/api", cors.New())

	api.Get("/v1/product/items", productController.Items)
	api.Get("/v1/product/items:id", productController.ItemById)
	api.Get("/v1/product/items/withname/:name", productController.ItemsWithName)
	api.Get("/v1/product/items/type/:catalogTypeId/brand/:catalogBrandId", productController.ItemsByTypeIdAndBrandId)
	api.Get("/v1/product/items/type/all/brand/:catalogBrandId", productController.ItemsByBrandId)
	api.Get("/v1/product/catalogtypes", productController.ProductTypes)
	api.Get("/v1/product/catalogbrands", productController.ProductBrands)
	api.Put("/v1/product/items", productController.UpdateProduct)
	api.Post("/v1/product/items", productController.CreateProduct)
	api.Delete("/v1/product/:id", productController.DeleteProduct)


	app.Use("/swagger", swagger.Handler)
	log.Fatal(app.Listen(":5000"))
}

func seedData(db *gorm.DB){
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		var count int
		db.Model(&models.ProductType{}).Count(&count)
		if count == 0 {
			fmt.Println("import sample ProductType data")
			seedProductTypes(db)
		}
		wg.Done()
	}()
	
	go func() {
		var count int
		db.Model(&models.ProductBrand{}).Count(&count)
		if count == 0 {
			fmt.Println("import sample ProductBrand data")
			seedProductBrands(db)
		}
		wg.Done()
	}()
	wg.Wait()

	var wgpd sync.WaitGroup
	wgpd.Add(1)
	go func() {
		var count int
		db.Model(&models.ProductItem{}).Count(&count)
		if count == 0 {
			fmt.Println("import sample ProductItem data")
			seedProductItems(db)
		}
		wgpd.Done()
	}()
	wgpd.Wait()
}

func seedProductTypes(db *gorm.DB){
	readProductTypesCSV(db)
}

func seedProductBrands(db *gorm.DB){
	readProductBrandsCSV(db)
}

func seedProductItems(db *gorm.DB){
	readProductItemsCSV(db)
}

func readProductTypesCSV(db *gorm.DB){
	f, err := os.Open("./data-sample/ProductTypes.csv")
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		db.Model(&models.ProductType{}).Create(&models.ProductType{
			Type: record[0],
		})
	}
}

func readProductBrandsCSV(db *gorm.DB){
	f, err := os.Open("./data-sample/ProductBrands.csv")
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		db.Model(&models.ProductBrand{}).Create(&models.ProductBrand{
			Brand: record[0],
		})
	}
}

func readProductItemsCSV(db *gorm.DB){
	f, err := os.Open("./data-sample/ProductItems.csv")
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		var tp models.ProductType
		db.Model(&models.ProductType{}).Where("type = ?",record[0]).First(&tp)

		var brand models.ProductBrand
		db.Model(&models.ProductBrand{}).Where("brand = ?",record[1]).First(&brand)

		price,_ := strconv.ParseFloat(record[4],32)
		as , _ := strconv.Atoi(record[6])
		or,_ := strconv.ParseBool(record[7])
		db.Model(&models.ProductItem{}).Create(&models.ProductItem{
			ProductBrandId: int(brand.ID),
			ProductBrand: brand,
			ProductTypeId: int(tp.ID),
			ProductType: tp,
			Description: record[2],
			Name: record[3],
			Price: float32(price),
			PictureFileName: record[5],
			AvailableStock: as,
			OnReorder: or,
		})
	}
}