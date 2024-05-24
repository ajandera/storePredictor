package endpoints

import (
	"encoding/json"
	"log"
	"main/utils"
	"net"
	"net/http"
	"net/url"
	"strings"

	model "github.com/ajandera/sp_model"
	"github.com/ajandera/sp_model/rdbsClientData"

	"github.com/bitly/go-simplejson"
)

// Order struct to store order data
type Order struct {
	Code       string
	OrderId    string
	Items      []rdbsClientData.Item
	TotalPrice float64
	Currency   string
	Url        string
	Tag        string
}

// Order struct to store order data
type ShopifyOrder struct {
	Currency          string
	CurrentTotalPrice float64
	LineItems         []ShopifyOrderItem
	OrderNumber       string
	OrderStatusUrl    string
}

type ShopifyOrderItem struct {
	Price       float64
	Sku         string
	Quantity    int8
	ProductCode string
	Name        string
	Tag         string
}

// Orders function to handle API requests to store orders in database

// POST
// @tags Orders
// @Summary Endpoint to store orders data for tracking
// @Description Endpoint to store orders data for tracking
// @Param  Order body Order true  "Order object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /orders  [post]
func Orders(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	// Declare a new Order struct.
	var o Order

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	// check code
	u, err := url.Parse(o.Url)
	if err != nil {
		log.Println(err)
	}

	response := simplejson.New()
	// for subdomain
	storeIdSubdomain := m.CheckStoreCode(o.Code, u.Hostname())
	if utils.IsValidUUID(storeIdSubdomain) == true {
		m.SaveOrder(o.TotalPrice, o.Currency, storeIdSubdomain, o.Items, o.OrderId, o.Tag)
		response.Set("success", true)
		payload, err := response.MarshalJSON()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	} else {
		parts := strings.Split(u.Hostname(), ".")
		domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
		storeId := m.CheckStoreCode(o.Code, domain)

		if utils.IsValidUUID(storeId) == true {
			m.SaveOrder(o.TotalPrice, o.Currency, storeId, o.Items, o.OrderId, o.Tag)
			response.Set("success", true)
		} else {
			response.Set("success", false)
			response.Set("error", "Not valid store code")
		}

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Println(err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	}
}

// Shopify Orders function to handle API requests to store orders from shopify in database

// POST
// @tags Orders
// @Summary Endpoint to store shopify orders data for tracking
// @Description Endpoint to store shopify orders data for tracking
// @Param  Order body Order true  "Order object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /shopify-orders  [post]
func ShopifyOrders(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	// Declare a new ShopifyStore struct.
	var o ShopifyOrder

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	storeUrl := o.OrderStatusUrl
	u, err := url.Parse(storeUrl)
	if err != nil {
		log.Println(err.Error())
		return
	}
	host, _, _ := net.SplitHostPort(u.Host)

	response := simplejson.New()
	// for subdomain
	storeId, err := m.GetStoreByUrl(host)
	var items []rdbsClientData.Item

	for _, element := range o.LineItems {
		items = append(items, rdbsClientData.Item{
			UnitPrice:   element.Price,
			Quantity:    element.Quantity,
			ProductCode: element.Sku,
			ProductName: element.Name,
			Tag:         "",
		})
	}

	m.SaveOrder(o.CurrentTotalPrice, o.Currency, storeId, items, o.OrderNumber, "")
	response.Set("success", true)
	payload, err := response.MarshalJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(payload)
}
