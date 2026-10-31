package endpoints

import (
	"encoding/json"
	"log"
	"main/repo"
	s "main/structs/stores"
	"main/utils"
	"net/http"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

// GET
// @tags Stores
// @Summary Endpoint to get stores per accountId
// @Description Endpoint to get stores per accountId
// @Accept  json
// @Produce  json
// @Success 200 {object} []s.Stores
// @Router /stores/{accountId}  [get]
func GetStoresByAccount(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if callerId, ok := authenticate(w, r); ok {
		vars := mux.Vars(r)
		accountId := vars["accountId"]
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		if len(accountId) <= 0 {
			http.Error(w, "Missing account Id.", http.StatusBadRequest)
			return
		}
		if !canViewAccount(m, callerId, accountId) {
			denyAccess(w)
			return
		}

		response := simplejson.New()
		response.Set("success", true)
		response.Set("stores", m.GetStoresByAccount(accountId))

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Print(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}

// POST
// @tags Stores
// @Summary Endpoint to store Stores data
// @Description Endpoint to store Stores data
// @Param  s.Stores body  s.Stores true  "Stores object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /stores  [post]
func Stores(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if callerId, ok := authenticate(w, r); ok {
		// Declare a new Visitor struct.
		var store s.Stores

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&store)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Print(err)
			return
		}

		// stores always belong to the main account of the caller's tenant
		store.AccountRefer = tenantOwner(m.GetAccountById(callerId))

		response := simplejson.New()
		const charset = "abcdefghijklmnopqrstuvwxyz" +
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		code := "SP-" + utils.StringWithCharset(5, charset)
		m.CreateStore(store.CountryCode, store.Url, code, store.AccountRefer, false, "", "", store.Feed, store.Window)
		response.Set("success", true)

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Print(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	}
}

// PUT
// @tags Stores
// @Summary Endpoint to update Stores data
// @Description Endpoint to update Stores data
// @Param  s.Stores body  s.Stores true  "Stores object to update in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /stores  [put]
func UpdateStores(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	// Declare a new Stores struct.
	var store s.Stores

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&store)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print(err.Error())
		return
	}

	if isAuthorized(w, r, m, store.ID) == true {
		response := simplejson.New()
		m.EditStore(store.ID, store.CountryCode, store.Url, store.MaximalProductPrice, store.MinimalProductPrice,
			store.ActualStorePower, store.ActualCustomerSatisfaction, store.PerceivedValue, store.ProductSell, false, store.Feed, store.Window)
		response.Set("success", true)

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Print(err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	}
}

// DELETE
// @tags Stores
// @Summary Endpoint to delete Stores
// @Description Endpoint to delete Stores
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /stores/{storeId}  [delete]
func DeleteStores(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]

	if isAuthorized(w, r, m, storeId) == true {
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		if len(storeId) <= 0 {
			http.Error(w, "Store id is not set.", http.StatusBadRequest)
			return
		}

		response := simplejson.New()
		m.DeleteStore(storeId)
		response.Set("success", true)

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Print(err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}

// GET
// @tags Stores
// @Summary Endpoint to get stores
// @Description Endpoint to get stores
// @Accept  json
// @Produce  json
// @Success 200 {object} []s.Stores
// @Router /stores  [get]
func GetStores(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	utils.SetupCORS(&w, r)
	log.Println((*r).Method)
	if (*r).Method == "OPTIONS" {
		return
	}
	if callerId, ok := authenticate(w, r); ok {
		// only stores of the caller's tenant
		response := simplejson.New()
		response.Set("success", true)
		response.Set("stores", m.GetStoresByAccount(tenantOwner(m.GetAccountById(callerId))))

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Print(err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}
