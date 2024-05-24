package endpoints

import (
	"encoding/json"
	"log"
	inv "main/structs/invoices"
	"main/utils"
	"net/http"
	"strconv"

	model "github.com/ajandera/sp_model"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

// GetInvoice function to handle API requests to get invoices

// GET
// @tags Invoices
// @Summary Endpoint to send logged user account
// @Description Endpoint to send logged user account
// @Accept  json
// @Produce  json
// @Success 200 {object} []inv.Invoices
// @Router /invoice/{storeId}  [get]
func GetInvoices(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]

	if isAuthorized(w, r, m, storeId) == true {
		response := simplejson.New()
		response.Set("success", true)
		response.Set("invoices", m.GetInvoices(storeId))

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

// CreateInvoice function to handle API requests to store invoice in database

// POST
// @tags Invoices
// @Summary Endpoint to store Invoices data
// @Description Endpoint to store Invoices data
// @Param  inv.Invoices body  inv.Invoices true  "Invoices object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /invoice  [post]
func CreateInvoice(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if isAuthorized(w, r, m, "") == true {
		// Declare a new Visitor struct.
		var invoice inv.Invoices

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&invoice)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf(err.Error())
			return
		}

		response := simplejson.New()
		amount, _ := strconv.ParseFloat(invoice.Amount, 64)
		m.CreateInvoice(invoice.DueDate, amount, invoice.Currency, invoice.StoreRefer)
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

// UpdateInvoices function to handle API requests to update Invoices in database

// PUT
// @tags Invoices
// @Summary Endpoint to update Invoices data
// @Description Endpoint to update Invoices data
// @Param  inv.Invoices body  inv.Invoices true  "Invoices object to update in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /invoice  [put]
func UpdateInvoice(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	// Declare a new Stores struct.
	var invoice inv.Invoices

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&invoice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf(err.Error())
		return
	}

	if isAuthorized(w, r, m, invoice.StoreRefer) == true {
		response := simplejson.New()
		amount, _ := strconv.ParseFloat(invoice.Amount, 64)
		m.UpdateInvoice(invoice.Id, invoice.DueDate, amount, invoice.Currency)
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

// GetInvoicesQuery function to handle API requests to get invoices

// GET
// @tags Invoices
// @Summary Endpoint to send logged user Invoice
// @Description Endpoint to send logged user Invoice
// @Accept  json
// @Produce  json
// @Success 200 {object} []inv.Invoices
// @Router /invoice/query/{storeId}/{from}/{to}  [get]
func GetInvoicesQuery(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	from, _ := strconv.Atoi(vars["from"])
	to, _ := strconv.Atoi(vars["to"])

	if isAuthorized(w, r, m, storeId) == true {
		response := simplejson.New()
		response.Set("success", true)
		response.Set("invoices", m.GetInvoicesFilter(storeId, from, to))

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

// DELETE
// @tags Invoices
// @Summary Endpoint to delete Invoice
// @Description Endpoint to delete Invoice
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /invoice/{invoiceId}  [delete]
func DeleteInvoice(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if isAuthorized(w, r, m, "") == true {
		vars := mux.Vars(r)
		supplierId := vars["invoiceId"]

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		if len(supplierId) <= 0 {
			http.Error(w, "Invoice Id id is not send", http.StatusBadRequest)
			return
		}

		response := simplejson.New()
		m.DeleteInvoice(supplierId)
		response.Set("success", true)

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Printf(err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}
