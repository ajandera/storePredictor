package endpoints

import (
	"fmt"
	"log"
	"main/utils"
	"net/http"
	"os"
	"strconv"

	model "github.com/ajandera/sp_model"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

// InfluxResponse struct to store response from influx client
type InfluxResponse struct {
	Index string
	Val   string
	Date  string
}

// GetPredictedVisitors function to handle API requests to get predicted visitors

// GET
// @tags Prediction
// @Summary Endpoint to get predicted visitors
// @Description Endpoint to get predicted visitors
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.ResPrediction
// @Router /prediction/visitors/{storeId}/{from}/{to}  [get]
func GetPredictedVisitors(w http.ResponseWriter, r *http.Request, m model.Repository, i model.Influx) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	from := vars["from"]
	to := vars["to"]

	if isAuthorized(w, r, m, storeId) == true {
		// load data from influx
		query := `from(bucket:"` + storeId + `") 
				|> range(start:-` + from + `mo, stop:` + to + `mo)
				|> filter(fn: (r) => (r._measurement == "visitors" and r._field == "value"))`
		data, err := i.GetInfluxData(query, os.Getenv("INFLUX_ORGANIZATION"))
		if err != nil {
			log.Printf(err.Error())
			return
		}
		// prepare data to response
		response := simplejson.New()
		response.Set("success", true)
		response.Set("prediction", data)

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

// GetPredictedVisitorsQuery function to handle API requests to get predicted visitors raw data

// GET
// @tags Prediction
// @Summary Endpoint to get predicted visitors raw data
// @Description Endpoint to get predicted visitors raw data
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.ResPrediction
// @Router /prediction/query/visitors/{storeId}/{from}/{to}  [get]
func GetPredictedVisitorsQuery(w http.ResponseWriter, r *http.Request, m model.Repository, i model.Influx) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	from := vars["from"]
	to := vars["to"]

	if isAuthorized(w, r, m, storeId) == true {
		visitors := m.GetVisitorsForPredictionView(from, to, storeId)
		response := simplejson.New()
		response.Set("success", true)
		response.Set("visitors", visitors)

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Printf("Marshall response: " + err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}

// GetPredictedOrders function to handle API requests to get predicted orders

// GET
// @tags Prediction
// @Summary Endpoint to get predicted orders
// @Description Endpoint to get predicted orders
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.ResPrediction
// @Router /prediction/orders/{storeId}/{from}/{to}  [get]
func GetPredictedOrders(w http.ResponseWriter, r *http.Request, m model.Repository, i model.Influx) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	from := vars["from"]
	to := vars["to"]

	if isAuthorized(w, r, m, storeId) == true {
		// load data from influx
		query := `from(bucket:"` + storeId + `") 
				|> range(start:-` + from + `mo, stop:` + to + `mo)  
				|> filter(fn: (r) => (r._measurement == "order" and (r._field == "value" or r._field == "saoa")))`
		data, errInf := i.GetInfluxData(query, os.Getenv("INFLUX_ORGANIZATION"))
		if errInf != nil {
			log.Printf(errInf.Error())
			return
		}
		// prepare data to response
		response := simplejson.New()
		response.Set("success", true)
		response.Set("prediction", data)

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

// GetPredictedOrdersQuery function to handle API requests to get predicted orders raw data

// GET
// @tags Prediction
// @Summary Endpoint to get predicted orders raw data
// @Description Endpoint to get predicted orders raw data
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.ResPrediction
// @Router /prediction/query/orders/{storeId}/{from}/{to}  [get]
func GetPredictedOrdersQuery(w http.ResponseWriter, r *http.Request, m model.Repository, i model.Influx) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	from := vars["from"]
	to := vars["to"]

	if isAuthorized(w, r, m, storeId) == true {
		// load data from influx
		query := `from(bucket:"` + storeId + `") 
				|> range(start:-` + from + `mo, stop:` + to + `mo)  
				|> filter(fn: (r) => (r._measurement == "order" and (r._field == "value" or r._field == "saoa")))`
		data, errInf := i.GetInfluxQuery(query, os.Getenv("INFLUX_ORGANIZATION"))
		if errInf != nil {
			log.Printf(errInf.Error())
			return
		}

		// Iterate over query response
		var influxResponse []InfluxResponse
		var saoa float64 = 0.0
		var counter float64 = 0.0
		for data.Next() {
			// Notice when group key has changed
			if data.TableChanged() {
				//fmt.Printf("table: %s\n", data.TableMetadata().String())
			}
			// Access data
			if data.Record().Field() == "value" {
				influxResponse = append(influxResponse, InfluxResponse{
					fmt.Sprint(data.Record().ValueByKey("daysToMeasurement")),
					fmt.Sprint(data.Record().Value()),
					data.Record().Time().String()})
			} else {
				avg, e := strconv.ParseFloat(fmt.Sprint(data.Record().Value()), 64)
				if e != nil {
					fmt.Println(e.Error())
				}
				saoa += avg
				counter++
			}
		}

		// prepare data to response
		response := simplejson.New()
		response.Set("success", true)
		response.Set("prediction", influxResponse)
		response.Set("saoa", saoa/counter)

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

// GetPredictedProducts function to handle API requests to get predicted products

// GET
// @tags Prediction
// @Summary Endpoint to get predicted products
// @Description Endpoint to get predicted products
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.ResPrediction
// @Router /prediction/products/{productCode}/{storeId}/{from}/{to}  [get]
func GetPredictedProducts(w http.ResponseWriter, r *http.Request, m model.Repository, i model.Influx) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	productCode := vars["productCode"]
	from := vars["from"]
	to := vars["to"]

	if isAuthorized(w, r, m, storeId) == true {
		// load data from influx
		query := `from(bucket:"` + storeId + `") 
				|> range(start:-` + from + `mo, stop:` + to + `mo)  
				|> filter(fn: (r) => (r._measurement == "` + productCode + `" and (r._field == "value")))`
		data, errInf := i.GetInfluxData(query, os.Getenv("INFLUX_ORGANIZATION"))
		if errInf != nil {
			log.Printf(errInf.Error())
			return
		}
		// prepare data to response
		response := simplejson.New()
		response.Set("success", true)
		response.Set("prediction", data)

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

// GetPredictedProductsQuery function to handle API requests to get product prediction raw data

// GET
// @tags Prediction
// @Summary Endpoint to to get product prediction raw data
// @Description Endpoint to get product prediction raw data
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.ResPrediction
// @Router /prediction/query/products/{productCode}/{storeId}/{from}/{to}  [get]
func GetPredictedProductsQuery(w http.ResponseWriter, r *http.Request, m model.Repository, i model.Influx) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	productCode := vars["productCode"]
	from := vars["from"]
	to := vars["to"]

	if isAuthorized(w, r, m, storeId) == true {
		// load data from influx
		query := `from(bucket:"` + storeId + `") 
				|> range(start:-` + from + `mo, stop:` + to + `mo)  
				|> filter(fn: (r) => (r._measurement == "` + productCode + `" and (r._field == "value")))`
		data, errInf := i.GetInfluxQuery(query, os.Getenv("INFLUX_ORGANIZATION"))
		if errInf != nil {
			log.Printf(errInf.Error())
			return
		}

		// Iterate over query response
		var influxResponse []InfluxResponse
		for data.Next() {
			// Notice when group key has changed
			if data.TableChanged() {
				//fmt.Printf("table: %s\n", data.TableMetadata().String())
			}
			// Access data
			influxResponse = append(influxResponse, InfluxResponse{
				fmt.Sprint(data.Record().ValueByKey("daysToMeasurement")),
				fmt.Sprint(data.Record().Value()),
				data.Record().Time().String()})
		}

		// prepare data to response
		response := simplejson.New()
		response.Set("success", true)
		response.Set("prediction", influxResponse)

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

// GetStats function to handle API requests to get stats

// GET
// @tags Prediction
// @Summary Endpoint to get stats
// @Description Endpoint to get stats
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.ResPrediction
// @Router /stats/{storeId}  [get]
func GetStats(w http.ResponseWriter, r *http.Request, m model.Repository, i model.Influx) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	if isAuthorized(w, r, m, storeId) == true {
		visitors := m.GetSumVisitors(storeId)
		orders := m.GetNumberOrders(storeId)
		ordersSum := m.GetSumOrder(storeId)
		var cr float64
		var lead float64

		if visitors > 0 {
			cr = (orders / visitors) * 100.00
			lead = ordersSum / visitors
		} else {
			cr = 0.00
			lead = 0.00
		}

		// prepare data to response
		response := simplejson.New()
		response.Set("success", true)
		response.Set("cr", cr)
		response.Set("lead", lead)
		response.Set("r2", GetPredictionR2(storeId, m, i))

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

// GetPredictionR2 function to calculate r-square

func GetPredictionR2(storeId string, m model.Repository, i model.Influx) float64 {
	var result float64
	sum := 0.00
	sumAvg := 0.00
	from := "3"
	to := "0"
	var influxResponse []InfluxResponse
	// load data from influx
	query := `from(bucket:"` + storeId + `") |> range(start:-` + from + `mo, stop:` + to + `mo)  |> filter(fn: (r) => (r._measurement == "order" and r._field == "value"))`
	data, _ := i.GetInfluxQuery(query, os.Getenv("INFLUX_ORGANIZATION"))
	// Iterate over query response
	var sumForAvg float64 = 0.0
	var counter float64 = 0.0
	for data.Next() {
		// Notice when group key has changed
		if data.TableChanged() {
			//fmt.Printf("table: %s\n", data.TableMetadata().String())
		}
		// Access data

		influxResponse = append(influxResponse, InfluxResponse{
			fmt.Sprint(data.Record().ValueByKey("daysToMeasurement")),
			fmt.Sprint(data.Record().Value()),
			data.Record().Time().String()})
	}

	for _, v := range influxResponse {
		var o float64
		if v.Index == "d0" {
			o, _ = strconv.ParseFloat(v.Val, 64)
			sumForAvg += o
			counter++
		}
	}
	avg := sumForAvg / counter
	for _, v := range influxResponse {
		if v.Index == "d0" {
			var orders float64
			var predictOrder float64
			orders, _ = strconv.ParseFloat(v.Val, 64)
			for _, v2 := range influxResponse {
				if v2.Date == v2.Date {
					predictOrder, _ = strconv.ParseFloat(v2.Val, 64)
				}
			}
			sum += (orders - predictOrder) * (orders - predictOrder)
			sumAvg += (orders - avg) * (orders - avg)
		}
	}

	// calculate r2
	result = 1 - (sum / sumAvg)
	return result
}
