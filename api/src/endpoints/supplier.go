package endpoints

import (
	"encoding/json"
	"log"
	sp "main/structs/suppliers"
	"main/utils"
	"net/http"

	model "github.com/ajandera/sp_model"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

// GetSuppliers function to handle API requests to get suppliers

// GET
// @tags Suppliers
// @Summary Endpoint to send logged user account
// @Description Endpoint to send logged user account
// @Accept  json
// @Produce  json
// @Success 200 {object} []sp.Suppliers
// @Router /supplier/{storeId}  [get]
func GetSuppliers(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]

	if isAuthorized(w, r, m, storeId) == true {
		response := simplejson.New()
		response.Set("success", true)
		response.Set("suppliers", m.GetSuppliers(storeId))

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

// GetSupplierDetail function to handle API requests to get supplier detail

// GET
// @tags Suppliers
// @Summary Endpoint to get supplier detail
// @Description Endpoint to get supplier detail
// @Accept  json
// @Produce  json
// @Success 200 {object} sp.Suppliers
// @Router /supplier/detail/{supplierId}  [get]
func GetSupplierDetail(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	supplierId := vars["supplierId"]
	storeId := vars["storeId"]

	if isAuthorized(w, r, m, storeId) == true {
		response := simplejson.New()
		response.Set("success", true)
		response.Set("supplier", m.GetSupplier(supplierId))

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

// CreateSupplier function to handle API requests to store supplier in database

// POST
// @tags Suppliers
// @Summary Endpoint to store Supplier data
// @Description Endpoint to store Supplier data
// @Param  sp.Suppliers body sp.Suppliers true  "Supplier object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /supplier  [post]
func CreateSupplier(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if isAuthorized(w, r, m, "") == true {
		// Declare a new Visitor struct.
		var supplier sp.Suppliers

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&supplier)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf(err.Error())
			return
		}

		response := simplejson.New()
		m.CreateSupplier(supplier.Name, supplier.Street, supplier.City, supplier.Zip, supplier.Country, supplier.Email,
			supplier.Phone, supplier.Person, supplier.StoreRefer, supplier.Template, supplier.Subject)
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

// CfreateSupplierOrder function to handle API requests to send order to supplier

// POST
// @tags Suppliers
// @Summary Endpoint to store send supplier order
// @Description Endpoint to send supplier order
// @Param  sp.Suppliers body  sp.Suppliers true  "Supplier object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /supplier/order  [post]
func CreateSupplierOrder(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if isAuthorized(w, r, m, "") == true {
		// Declare a new Visitor struct.
		var order sp.Order

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf(err.Error())
			return
		}

		response := simplejson.New()
		utils.SendEmailWithoutTemplate(order.Recipient, order.Subject, order.Template)
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

// UpdateSupplier function to handle API requests to update supplier in database

// PUT
// @tags Suppliers
// @Summary Endpoint to update Supplier data
// @Description Endpoint to update Supplier data
// @Param  sp.Suppliers body  sp.Suppliers true  "Supplier object to update in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /supplier  [put]
func UpdateSupplier(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	// Declare a new Stores struct.
	var supplier sp.Suppliers

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&supplier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf(err.Error())
		return
	}

	if isAuthorized(w, r, m, supplier.StoreRefer) == true {
		response := simplejson.New()
		m.UpdateSupplier(supplier.Id, supplier.Name, supplier.Street, supplier.City, supplier.Zip, supplier.Country, supplier.Email,
			supplier.Phone, supplier.Person, supplier.Template, supplier.Subject)
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

// DeleteSupplier function to handle API requests to remove supplier from database

// DELETE
// @tags Suppliers
// @Summary Endpoint to delete Supplier
// @Description Endpoint to delete Supplier
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /supplier/{supplierId}  [delete]
func DeleteSupplier(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if isAuthorized(w, r, m, "") == true {
		vars := mux.Vars(r)
		supplierId := vars["supplierId"]

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		if len(supplierId) <= 0 {
			http.Error(w, "Supplier Id id is not send", http.StatusBadRequest)
			return
		}

		response := simplejson.New()
		m.DeleteSupplier(supplierId)
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
