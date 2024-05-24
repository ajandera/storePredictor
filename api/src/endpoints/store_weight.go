package endpoints

import (
	"encoding/json"
	"log"
	sw "main/structs/storeWeights"
	"main/utils"
	"net/http"

	model "github.com/ajandera/sp_model"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

// GET StoreWeights
// @tags StoreWeights
// @Summary Endpoint to get store weights for store
// @Description Endpoint to get store weights for store
// @Accept  json
// @Produce  json
// @Success 200 {object} []sw.StoreWeights
// @Router /store-weights/{storeId}  [get]
func GetStoreWeights(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if isAuthorized(w, r, m, "") == true {
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
		response.Set("opendata", m.GetStoreWeights(storeId))

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
// @tags StoreWeights
// @Summary Endpoint to store StoreWeights data
// @Description Endpoint to StoreWeights data
// @Param  sw.StoreWeights body  sw.StoreWeights true  "StoreWeights object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /accounts/{lang}  [post]
func CreateStoreWeights(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if isAuthorized(w, r, m, "") == true {
		// Declare a new Visitor struct.
		var storeWeights sw.StoreWeights

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&storeWeights)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf(err.Error())
			return
		}

		response := simplejson.New()
		m.CreateStoreWeights(storeWeights.StoreRefer, storeWeights.Name, storeWeights.Beta, storeWeights.Gama,
			storeWeights.Delta, storeWeights.A, storeWeights.B, storeWeights.C, storeWeights.D, storeWeights.E,
			storeWeights.ProbabilityWeights, storeWeights.Shift, storeWeights.LongShift)
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

// PUT
// @tags StoreWeights
// @Summary Endpoint to update StoreWeights data
// @Description Endpoint to update StoreWeights data
// @Param  sw.StoreWeights body  sw.StoreWeights true  "StoreWeights object to update in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /store-weights  [put]
func UpdateStoreWeights(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if isAuthorized(w, r, m, "") == true {
		// Declare a new Visitor struct.
		var storeWeights sw.StoreWeights

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&storeWeights)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf(err.Error())
			return
		}

		response := simplejson.New()
		m.EditStoreWeights(storeWeights.StoreRefer, storeWeights.Name, storeWeights.Beta, storeWeights.Gama,
			storeWeights.Delta, storeWeights.A, storeWeights.B, storeWeights.C, storeWeights.D, storeWeights.E, storeWeights.ProbabilityWeights, storeWeights.Shift, storeWeights.LongShift)
		response.Set("success", true)

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Printf(err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	}
}
