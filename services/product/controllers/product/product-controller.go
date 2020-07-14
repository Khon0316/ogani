package productController

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"ogani.com/services/product/models"
	"strconv"
)

// Items godoc
// @Summary Get all catalogs
// @Description Get all catalogs
// @Accept json
// @Produce json
// @Param  size query int true "it's page size"
// @Param  pageIndex query int true "it's page index"
// @Success 200
// @BadRequest 400
// @Router /product/items [get]
func Items(w http.ResponseWriter, r *http.Request) {
	pageSize,err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err !=nil{
		log.Fatalln(err)
		badRequest(w)
		return
	}

	pageIndex, err := strconv.Atoi(r.URL.Query().Get("pageIndex"))
	if err !=nil{
		log.Fatalln(err)
		badRequest(w)
		return
	}

	db := connectDbHandler()
	defer db.Close()

	totalItemsCh := make(chan int)
	go func(){
		var _totalItems = 0
		db.Model(&models.ProductItem{}).Count(_totalItems)
		totalItemsCh <- _totalItems
	}()
	itemsOnPageCh := make(chan []models.ProductItem)
	go func() {
		var _itemsOnPage []models.ProductItem
		db.Model(&models.ProductItem{}).Order("name asc").Offset(pageIndex * pageSize).Limit(pageSize).Find(_itemsOnPage)
		itemsOnPageCh <- _itemsOnPage
	}()

	totalItems := <-totalItemsCh
	itemsOnPage := <- itemsOnPageCh

	fmt.Println(totalItems)
	fmt.Println(itemsOnPage)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		pageIndex int
		pageSize int
		totalItems int
		itemsOnPage []models.ProductItem
	}{
		pageIndex: pageIndex,
		pageSize: pageSize,
		totalItems: totalItems,
		itemsOnPage: itemsOnPage,
	})
}

func ItemById(w http.ResponseWriter, r *http.Request) {
	id, err :=  strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		badRequest(w)
		return
	}

	db := connectDbHandler()
	defer db.Close()

	var product models.ProductItem
	db.Model(&models.ProductItem{}).Where("id = ?", id).Find(&product)

	product.FillProductUrl("base url")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func ItemsWithName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageSize,err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err !=nil{
		log.Fatalln(err)
		badRequest(w)
		return
	}

	pageIndex, err := strconv.Atoi(r.URL.Query().Get("pageIndex"))
	if err !=nil{
		log.Fatalln(err)
		badRequest(w)
		return
	}

	db:= connectDbHandler()

	totalItems := 0
	db.Model(&models.ProductItem{}).Where("name LIKE %?",vars["name"]).Count(totalItems)

	var itemsOnPage []models.ProductItem
	db.Model(&models.ProductItem{}).Order("name asc").Offset(pageIndex * pageSize).Limit(pageSize).Find(itemsOnPage)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		pageIndex int
		pageSize int
		totalItems int
		itemsOnPage []models.ProductItem
	}{
		pageIndex: pageIndex,
		pageSize: pageSize,
		totalItems: totalItems,
		itemsOnPage: itemsOnPage,
	})
}

func ItemsByTypeIdAndBrandId(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func ItemsByBrandId(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func ProductTypes(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func ProductBrands(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func badRequest(w http.ResponseWriter){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
}

func connectDbHandler() *gorm.DB{
	db, err := gorm.Open("postgres","host=0.0.0.0 port=5432 user=postgres dbname=ogani password=postgres sslmode=disable")
	if err != nil{
		panic(err)
	}
	return db
}