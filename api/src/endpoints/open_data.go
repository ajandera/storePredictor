package endpoints

import (
	"encoding/json"
	"log"
	od "main/structs/openData"
	"main/utils"
	"net/http"

	model "github.com/ajandera/sp_model"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

// GET
// @tags OpenData
// @Summary Endpoint to get opendata for store
// @Description Endpoint to get opendata for store
// @Accept  json
// @Produce  json
// @Success 200 {object} []od.OpenData
// @Router /open-data/{storeId}  [get]
func GetOpenData(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	//if isAuthorized(w, r) == true {
	vars := mux.Vars(r)
	storeId := vars["storeId"]

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	if len(storeId) <= 0 {
		http.Error(w, "Store id not send", http.StatusBadRequest)
		return
	}

	response := simplejson.New()
	response.Set("success", true)
	response.Set("opendata", m.GetOpenData(storeId))

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Printf(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// POST
// @tags OpenData
// @Summary Endpoint to save open data for store
// @Description Endpoint to save open data for store
// @Param  od.OpenData body  od.OpenData true  "OpenData object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /open-data  [post]
func AddOpenData(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if isAuthorized(w, r, m, "") == true {
		// Declare a new Visitor struct.
		var openData od.OpenData

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&openData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf(err.Error())
			return
		}

		response := simplejson.New()
		m.CreateOpenData(openData.StorePower, openData.CustomerSatisfaction, openData.MaximalProductPrice,
			openData.MinimalProductPrice, openData.PerceivedValue, openData.StoreRefer)
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
