// Main package to handle store REST API
package main

import (
	"io/ioutil"
	"log"
	"main/bot/adapters/storage"
	_ "main/docs"
	endpoints "main/endpoints"
	structs "main/structs/bot"
	"main/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	model "github.com/ajandera/sp_model"
	"github.com/wagslane/go-rabbitmq"

	"github.com/gorilla/mux"
	_ "github.com/streadway/amqp"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title storePredictor shop REST API
// @version 1.0
// @description Store REST API to track visitorrs and orders

// @contact.name Ales Jandera
// @contact.url http://storepredictor.com
// @contact.email ales@storepredictor.com

// @license.name Commercial

// @host localhost:9999
// @BasePath /v1/
// @query.collection.format multi
func main() {
	clientsDataDsn := "host=" + os.Getenv("CLIENTS_DATA_HOST") + " user=" + os.Getenv("CLIENTS_DATA_USER") + " password=" + os.Getenv("CLIENTS_DATA_PASSWORD") + " dbname=" + os.Getenv("CLIENTS_DATA_DATABASE") + " port=" + os.Getenv("CLIENTS_DATA_PORT") + " sslmode=disable"
	clientsInformationDataDsn := "host=" + os.Getenv("CLIENTS_INFORMATION_HOST") + " user=" + os.Getenv("CLIENTS_INFORMATION_USER") + " password=" + os.Getenv("CLIENTS_INFORMATION_PASSWORD") + " dbname=" + os.Getenv("CLIENTS_INFORMATION_DATABASE") + " port=" + os.Getenv("CLIENTS_INFORMATION_PORT") + " sslmode=disable"
	repository := model.ClientsInit(
		clientsDataDsn,
		clientsInformationDataDsn)
	influxURL := os.Getenv("INFLUX_HOST")
	influx := model.ClientPredictedDataInit(
		influxURL,
		os.Getenv("INFLUX_TOKEN"))

	// connect to rabbitmq
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
		log.Println(err.Error())
	}

	var storageStruct []structs.StorageForAI

	// define storage for generative AI chatbot
	chatBotData, _ := filepath.Abs("./")
	files, err := ioutil.ReadDir(chatBotData + "/data/")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".gob") {
			store, err := storage.NewSeparatedMemoryStorage(chatBotData + "/data/" + file.Name())
			if err != nil {
				log.Println(err)
			}
			name := strings.Replace(file.Name(), ".gob", "", 1)
			storageStruct = append(storageStruct, structs.StorageForAI{Name: name, Storage: store})
		}
	}

	r := mux.NewRouter()
	// router refix
	api := r.PathPrefix("/v1").Subrouter()
	// accounts
	api.HandleFunc("/accounts/{lang}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.Accounts(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetAccounts(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		endpoints.UpdateAccounts(w, r, repository)
	}).Methods(http.MethodPut, http.MethodOptions)
	api.HandleFunc("/accounts/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.DeleteAccounts(w, r, repository)
	}).Methods(http.MethodDelete, http.MethodOptions)
	api.HandleFunc("/account/email/{email}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetAccountByEmail(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)
	// accounts
	api.HandleFunc("/accounts/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetAccount(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/account/child/{id}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetChildAccount(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)

	// prediction
	api.HandleFunc("/prediction/visitors/{storeId}/{from}/{to}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetPredictedVisitors(w, r, repository, influx)
	}).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/prediction/orders/{storeId}/{from}/{to}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetPredictedOrders(w, r, repository, influx)
	}).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/prediction/products/{productCode}/{storeId}/{from}/{to}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetPredictedProducts(w, r, repository, influx)
	}).Methods(http.MethodGet, http.MethodOptions)

	api.HandleFunc("/prediction/query/visitors/{storeId}/{from}/{to}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetPredictedVisitorsQuery(w, r, repository, influx)
	}).Methods(http.MethodGet, http.MethodOptions)

	api.HandleFunc("/prediction/query/orders/{storeId}/{from}/{to}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetPredictedOrdersQuery(w, r, repository, influx)
	}).Methods(http.MethodGet, http.MethodOptions)

	api.HandleFunc("/prediction/query/products/{productCode}/{storeId}/{from}/{to}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetPredictedProductsQuery(w, r, repository, influx)
	}).Methods(http.MethodGet, http.MethodOptions)

	// openData
	api.HandleFunc("/open-data", func(w http.ResponseWriter, r *http.Request) {
		endpoints.AddOpenData(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/open-data/{storeId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetOpenData(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)

	// stores
	api.HandleFunc("/stores", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetStores(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/stores", func(w http.ResponseWriter, r *http.Request) {
		endpoints.Stores(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/stores/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetStoresByAccount(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/stores", func(w http.ResponseWriter, r *http.Request) {
		endpoints.UpdateStores(w, r, repository)
	}).Methods(http.MethodPut, http.MethodOptions)
	api.HandleFunc("/stores/{storeId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.DeleteStores(w, r, repository)
	}).Methods(http.MethodDelete, http.MethodOptions)

	// storeWeights
	api.HandleFunc("/store-weights", func(w http.ResponseWriter, r *http.Request) {
		endpoints.CreateStoreWeights(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/store-weights", func(w http.ResponseWriter, r *http.Request) {
		endpoints.UpdateStoreWeights(w, r, repository)
	}).Methods(http.MethodPut, http.MethodOptions)
	api.HandleFunc("/store-weights/{storeId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetStoreWeights(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)

	// authorization
	api.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		utils.Auth(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		utils.Logout(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		utils.Refresh(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/account/me", func(w http.ResponseWriter, r *http.Request) {
		endpoints.AccountMe(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)

	// health check
	api.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		endpoints.HealthCheck(w, r)
	}).Methods(http.MethodGet, http.MethodOptions)

	// forgot
	api.HandleFunc("/forgot/{accountId}/{lang}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.SendNewPw(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/restore/{token}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.UpdatePassword(w, r, repository)
	}).Methods(http.MethodPut, http.MethodOptions)

	// order
	api.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		endpoints.AddOrder(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)

	// upgrade
	api.HandleFunc("/upgrade/{accountId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.UpgradeOrder(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)

	// account orders
	api.HandleFunc("/account/order/{accountId}/{storeId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetAccountOrders(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)

	api.HandleFunc("/account/order/{accountId}/{storeId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.AccountOrder(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)

	api.HandleFunc("/prediction/onetime/{orderId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.OnetimePrediction(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)

	api.HandleFunc("/prediction/manual/{storeId}/{date}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.ManualPrediction(w, r, publisher)
	}).Methods(http.MethodPost, http.MethodOptions)

	// orders
	api.HandleFunc("/orders/{storeId}/{limit}/{offset}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetClientOrders(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/orders/product/{storeId}/{productCode}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetOrdersForProduct(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)

	// ask a bot
	api.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request) {
		endpoints.ProceedAsk(w, r, storageStruct)
	}).Methods(http.MethodPost, http.MethodOptions)

	// train a bot
	api.HandleFunc("/train/{corpora}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.ProceedTrain(w, r)
	}).Methods(http.MethodPost, http.MethodOptions)

	// parse a feed
	api.HandleFunc("/feed/{corpora}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.ProceedFeed(w, r)
	}).Methods(http.MethodPost, http.MethodOptions)

	// visitors endpoint
	api.HandleFunc("/visitors", func(w http.ResponseWriter, r *http.Request) {
		endpoints.Visitors(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/visitors-offline", func(w http.ResponseWriter, r *http.Request) {
		endpoints.VisitorsOffline(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)
	// order endpoint
	api.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		endpoints.Orders(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)
	// token generator endpoint
	api.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		endpoints.Token(w, r)
	}).Methods(http.MethodGet, http.MethodOptions)

	// products
	api.HandleFunc("/products/{storeId}/{limit}/{offset}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetProducts(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/products/needs/{storeId}/{limit}/{offset}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetProductsToStore(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/products/needs/{storeId}/{productCode}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetProductToStore(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)

	// products warehouse
	api.HandleFunc("/warehouse/products/{storeId}/{limit}/{offset}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetProductsWarehouse(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/warehouse/product/{code}/{storeId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetProductWarehouse(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/warehouse/product", func(w http.ResponseWriter, r *http.Request) {
		endpoints.CreateProduct(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/warehouse/product", func(w http.ResponseWriter, r *http.Request) {
		endpoints.UpdateProduct(w, r, repository)
	}).Methods(http.MethodPut, http.MethodOptions)

	// suppliers
	api.HandleFunc("/supplier/{storeId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetSuppliers(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/supplier/detail/{supplierId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetSupplierDetail(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/supplier", func(w http.ResponseWriter, r *http.Request) {
		endpoints.CreateSupplier(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/supplier", func(w http.ResponseWriter, r *http.Request) {
		endpoints.UpdateSupplier(w, r, repository)
	}).Methods(http.MethodPut, http.MethodOptions)
	api.HandleFunc("/supplier/{supplierId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.DeleteSupplier(w, r, repository)
	}).Methods(http.MethodDelete, http.MethodOptions)
	api.HandleFunc("/supplier/order", func(w http.ResponseWriter, r *http.Request) {
		endpoints.CreateSupplierOrder(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)

	// invoices
	api.HandleFunc("/invoice/{storeId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetInvoices(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/invoice", func(w http.ResponseWriter, r *http.Request) {
		endpoints.CreateInvoice(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/invoice", func(w http.ResponseWriter, r *http.Request) {
		endpoints.UpdateInvoice(w, r, repository)
	}).Methods(http.MethodPut, http.MethodOptions)
	api.HandleFunc("/invoice/import", func(w http.ResponseWriter, r *http.Request) {
		endpoints.ImportInvoices(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/invoice/query/{storeId}/{from}/{to}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetInvoicesQuery(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/invoice/{invoiceId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.DeleteInvoice(w, r, repository)
	}).Methods(http.MethodDelete, http.MethodOptions)

	// stats
	api.HandleFunc("/stats/{storeId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetStats(w, r, repository, influx)
	}).Methods(http.MethodGet, http.MethodOptions)

	// pdf report
	api.HandleFunc("/reports/pdf/{storeId}/{date}/{lang}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetPDF(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)

	// excel store import
	api.HandleFunc("/warehouse/import/{storeId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.ImportWarehouse(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)

	// excel report
	api.HandleFunc("/reports/excel/{storeId}/{date}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetExcel(w, r, repository)
	}).Methods(http.MethodGet, http.MethodOptions)

	// excel report
	api.HandleFunc("/reports/excelfull/{storeId}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.GetExcelFull(w, r, repository, influx)
	}).Methods(http.MethodGet, http.MethodOptions)

	// share report
	api.HandleFunc("/reports/share/{storeId}/{date}/{lang}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.ShareReport(w, r, repository)
	}).Methods(http.MethodPost, http.MethodOptions)

	// serve the js lib to show data
	exePath, _ := filepath.Abs("./")
	s := http.StripPrefix("/public/", http.FileServer(http.Dir(exePath+"/lib/")))
	r.PathPrefix("/public/").Handler(s)

	s2 := http.StripPrefix("/files/", http.FileServer(http.Dir(exePath+"/temp/")))
	r.PathPrefix("/public/").Handler(s2)

	// add swagger
	if os.Getenv("API_DOC") == "true" {
		r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	}

	log.Fatal(http.ListenAndServe(":8888", r))
}
