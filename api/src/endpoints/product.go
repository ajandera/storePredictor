package endpoints

import (
	"encoding/json"
	"log"
	pr "main/structs/product"
	"main/utils"
	"net/http"
	"strconv"

	model "github.com/ajandera/sp_model"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

// GetProducts function to handle API requests to sell products

// GET
// @tags Products
// @Summary Endpoint to return sell products
// @Description Endpoint to return sell products
// @Accept  json
// @Produce  json
// @Success 200 {object} []rdbsClientData.TopSellProduct
// @Router /products/{storeId}/{limit}/{offset}  [get]
func GetProducts(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	limit := vars["limit"]
	offset := vars["offset"]
	if isAuthorized(w, r, m, storeId) == true {
		response := simplejson.New()
		response.Set("success", true)
		response.Set("products", m.GetProducts(map[string]interface{}{"store_id": storeId, "limit": limit, "offset": offset}))

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Printf(err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}

// GetProductsWarehouse function to handle API requests to return products stock

// GET
// @tags Products
// @Summary Endpoint to return products stock
// @Description Endpoint to return products stock
// @Accept  json
// @Produce  json
// @Success 200 {object} []pr.Product
// @Router /warehouse/products/{storeId}/{limit}/{offset}  [get]
func GetProductsWarehouse(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	limit := vars["limit"]
	offset := vars["offset"]
	if isAuthorized(w, r, m, storeId) == true {
		response := simplejson.New()
		response.Set("success", true)
		l, _ := strconv.ParseInt(limit, 10, 64)
		o, _ := strconv.ParseInt(offset, 10, 64)
		response.Set("products", m.GetProductsWarehouse(storeId, int(l), int(o)))

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Printf(err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}

// GetProductsWarehouse function to handle API requests to return products stock by product code

// GET
// @tags Products
// @Summary Endpoint to return products stock by product code
// @Description Endpoint to return products stock by product code
// @Accept  json
// @Produce  json
// @Success 200 {object} []pr.Product
// @Router /warehouse/product/{code}/{storeId}  [get]
func GetProductWarehouse(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	productCode := vars["code"]
	if isAuthorized(w, r, m, storeId) == true {
		response := simplejson.New()
		response.Set("success", true)
		response.Set("product", m.GetProduct(productCode, storeId))

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Printf(err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}

// GetProductsToStore function to handle API requests to return products to be ordered

// GET
// @tags Products
// @Summary Endpoint to return products to be ordered
// @Description Endpoint to return products to be ordered
// @Accept  json
// @Produce  json
// @Success 200 {object} []rdbsClientData.ProductToStore
// @Router /products/needs/{storeId}/{limit}/{offset}  [get]
func GetProductsToStore(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	limit := vars["limit"]
	offset := vars["offset"]
	if isAuthorized(w, r, m, storeId) == true {
		response := simplejson.New()
		response.Set("success", true)
		l, _ := strconv.ParseInt(limit, 10, 64)
		o, _ := strconv.ParseInt(offset, 10, 64)
		response.Set("products", m.GetProductsToStore(storeId, int(l), int(o)))

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Printf(err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}

// GetProductsToStore function to handle API requests to return products to be ordered by product code

// GET
// @tags Products
// @Summary Endpoint to return products to be orderred by product code
// @Description Endpoint to return products to be ordered by product code
// @Accept  json
// @Produce  json
// @Success 200 {object} []rdbsClientData.ProductToStore
// @Router /products/needs/{storeId}/{productCode}  [get]
func GetProductToStore(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	code := vars["productCode"]
	if isAuthorized(w, r, m, storeId) == true {
		response := simplejson.New()
		response.Set("success", true)
		response.Set("product", m.GetProductToStore(code, storeId))

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Printf(err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}

// CreateProduct function to handle API requests to store product in database

// POST
// @tags Products
// @Summary Endpoint to store product data
// @Description Endpoint to store product data
// @Param  pr.Product body  pr.Product true  "Product object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /warehouse/product  [post]
func CreateProduct(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if isAuthorized(w, r, m, "") == true {
		// Declare a new Visitor struct.
		var product pr.Product

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf(err.Error())
			return
		}

		response := simplejson.New()
		m.CreateProduct(product.ProductCode, product.Name, product.Quantity, product.StoreId)
		response.Set("success", true)

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Printf(err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	}
}

// UpdateProduct function to handle API requests to update product in database

// PUT
// @tags Products
// @Summary Endpoint to update product data
// @Description Endpoint to update product data
// @Param  pr.Product body  pr.Product true  "Peoduct object to update in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /warehouse/product  [put]
func UpdateProduct(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	// Declare a new Stores struct.
	var p pr.Product

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf(err.Error())
		return
	}

	if isAuthorized(w, r, m, p.StoreId) == true {
		response := simplejson.New()
		m.UpdateProduct(p.ProductCode, p.Name, p.StoreId, p.Quantity)
		response.Set("success", true)

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Printf(err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	}
}
