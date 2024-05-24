package endpoints

import (
	"encoding/json"
	"log"
	o "main/structs/order"
	"main/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	model "github.com/ajandera/sp_model"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	_ "github.com/rabbitmq/amqp091-go"
	_ "github.com/streadway/amqp"
	"github.com/wagslane/go-rabbitmq"
)

// POST
// @tags Order
// @Summary Endpoint to store Order data
// @Description Endpoint to store Order data
// @Param  o.Order body  o.Order true  "Order object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /order  [post]
func AddOrder(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	// Declare a new Visitor struct.
	var order o.Order

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf(err.Error())
		return
	}

	response := simplejson.New()
	account := m.EditAccount("", order.Name, order.Email, order.Street, order.City, order.Zip, order.CountryCode,
		order.CompanyNumber, order.VatNumber, order.PaidTo, order.PlanRefer, "", "", order.Password, true)
	response.Set("success", true)
	response.Set("account", account)

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Printf(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(payload)
}

// POST
// @tags Order
// @Summary Endpoint to store Order data from account
// @Description Endpoint to store Order data from account
// @Param  o.Order body  o.Order true  "Order object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /account/order/{accountId}/{storeId}  [post]
func AccountOrder(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	accountId := vars["accountId"]
	storeId := vars["storeId"]
	if isAuthorized(w, r, m, storeId) == true {
		var order o.Prediction
		err := json.NewDecoder(r.Body).Decode(&order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf(err.Error())
			return
		}

		// prepare data to response
		response := simplejson.New()
		response.Set("success", true)
		response.Set("order", m.CreateOrder(accountId, storeId, order.PlanRefer, order.Amount, order.Paid))

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

// POST
// @tags Order
// @Summary Endpoint to store Order for one time prediction data
// @Description Endpoint to store Order for one time prediction data
// @Param  o.Order body  o.Order true  "Order object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /prediction/onetime/{orderId}  [post]
func OnetimePrediction(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	orderId := vars["orderId"]
	order := m.GetOrderById(orderId)
	log.Println("Start publishing")
	conn, err := rabbitmq.NewConn(
		"amqp://"+os.Getenv("RABBIT_USER")+":"+os.Getenv("RABBIT_PASS")+"@"+os.Getenv("RABBIT_HOST"),
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
	)

	if err != nil {
		log.Println(err)
	}

	log.Println("Predict store id: " + order.StoreRefer)
	message := simplejson.New()
	message.Set("storeId", order.StoreRefer)
	message.Set("d0", time.Now().AddDate(0, 0, -1).Format(time.RFC3339))
	message.Set("type", 0)
	message.Set("product", "")
	message.Set("plan", false)
	res, _ := message.MarshalJSON()
	errPublisher := publisher.Publish(
		res,
		[]string{"prediction"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
	)
	if errPublisher != nil {
		log.Printf(errPublisher.Error())
		return
	}

	// prepare data to response
	response := simplejson.New()
	response.Set("success", true)

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Printf(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// GET
// @tags Order
// @Summary Endpoint to get orders for each Client
// @Description Endpoint to get orders for each Client
// @Param  o.Order body  o.Order true  "Order object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /orders/{storeId}/{limit}/{offset}  [get]
func GetClientOrders(w http.ResponseWriter, r *http.Request, m model.Repository) {
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
		lim, _ := strconv.Atoi(limit)
		off, _ := strconv.Atoi(offset)
		response.Set("orders", m.GetOrders(map[string]interface{}{"store_id": storeId}, lim, off))

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

// GET
// @tags Order
// @Summary Endpoint to get Orders for each product by code
// @Description Endpoint to get Orders for each product by code
// @Param  o.Order body  o.Order true  "Order object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /orders/product/{storeId}/{productCode}  [get]
func GetOrdersForProduct(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	productCode := vars["productCode"]
	if isAuthorized(w, r, m, storeId) == true {
		response := simplejson.New()
		response.Set("success", true)
		response.Set("orders", m.GetOrdersWithProduct(productCode, storeId))

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

// POST
// @tags Order
// @Summary Endpoint to start manualy prediction for set day
// @Description Endpoint to start manualy prediction for set day
// @Param  o.Order body  o.Order true  "Order object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /prediction/manual/{storeId}/{date}  [post]
func ManualPrediction(w http.ResponseWriter, r *http.Request, publisher *rabbitmq.Publisher) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	date := vars["date"]

	log.Println("Start manual publishing")

	// prepare data to response
	response := simplejson.New()

	predictDate, parseErr := time.Parse("2006-01-02", date)
	if parseErr != nil {
		log.Println(parseErr.Error())
		response.Set("parseErr", parseErr.Error())
	} else {
		message := simplejson.New()
		message.Set("storeId", storeId)
		message.Set("d0", predictDate.Format(time.RFC3339))
		message.Set("type", 0)
		message.Set("product", "")
		message.Set("plan", false)
		res, _ := message.MarshalJSON()
		errPublisher := publisher.Publish(
			res,
			[]string{"prediction"},
			rabbitmq.WithPublishOptionsContentType("application/json"),
		)
		if errPublisher != nil {
			log.Println(errPublisher.Error())
			response.Set("errPublisher", errPublisher.Error())
		}
	}

	response.Set("success", true)

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// POST
// @tags Order
// @Summary Endpoint to store Order data
// @Description Endpoint to store Order data
// @Param  o.Order body  o.Order true  "Order object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /upgrade/{accountId}  [post]
func UpgradeOrder(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	// Declare a new Visitor struct.
	var upgrade o.Upgrade
	vars := mux.Vars(r)
	accountId := vars["accountId"]

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&upgrade)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf(err.Error())
		return
	}

	response := simplejson.New()
	account := m.GetAccountById(accountId)
	paidTo := time.Now()
	if upgrade.Year == true {
		paidTo.AddDate(1, 0, 0)
	} else {
		paidTo.AddDate(0, 1, 0)
	}
	response.Set("upgrade", upgrade)

	updatedAccount := m.EditAccount(account.Id.String(), account.Name, account.Email, account.Street, account.City, account.Zip, account.CountryCode,
		account.CompanyNumber, account.VatNumber, paidTo.Format(time.RFC3339), upgrade.PlanRefer, "", "", "", true)

	response.Set("success", true)
	response.Set("account", updatedAccount)

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Printf(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(payload)
}
