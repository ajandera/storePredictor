package endpoints

import (
	"fmt"
	"io"
	"log"
	"main/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	model "github.com/ajandera/sp_model"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"github.com/xuri/excelize/v2"
)

// ImportWarehouse function to handle API requests to store Product in database

// POST
// @tags Import
// @Summary Endpoint to store Product data
// @Description Endpoint to store Product data
// @Param  rdbsClientData.Product body  rdbsClientData.Product true  "Product object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /warehouse/import/{storeId}  [post]
func ImportWarehouse(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	if isAuthorized(w, r, m, storeId) == true {
		file, handle, err := r.FormFile("importFile")

		if err != nil {
			log.Printf("From file: " + err.Error())
		}

		defer file.Close()

		// Create file
		dst, err := os.Create(handle.Filename)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Copy the uploaded file to the created file on the filesystem
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := excelize.OpenFile(handle.Filename)

		if err != nil {
			log.Printf(err.Error())
		}

		// Get all the rows in the Sheet1.
		rows, err := f.GetRows("Sheet1")
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, row := range rows {
			qty, _ := strconv.ParseInt(row[2], 0, 8)
			// check if exist
			product := m.GetProduct(row[0], storeId)
			if product.ProductCode != "" {
				m.UpdateProduct(row[0], row[1], storeId, int8(qty))
			} else {
				m.CreateProduct(row[0], row[1], int8(qty), storeId)
			}
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
}

// ImportInvoices function to handle API requests to store Invoice in database

// POST
// @tags Import
// @Summary Endpoint to store Invoices data
// @Description Endpoint to store Invoices data
// @Param  rdbsClientInfo.Invoices body rdbsClientInfo.Invoices true  "Invoices object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /invoice/import  [post]
func ImportInvoices(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	if isAuthorized(w, r, m, storeId) == true {
		file, handle, err := r.FormFile("importFile")

		if err != nil {
			log.Printf("From file: " + err.Error())
		}

		defer file.Close()

		// Create file
		dst, err := os.Create(handle.Filename)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Copy the uploaded file to the created file on the filesystem
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := excelize.OpenFile(handle.Filename)

		if err != nil {
			log.Printf(err.Error())
		}

		// Get all the rows in the Sheet1.
		rows, err := f.GetRows("Sheet1")
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, row := range rows {
			amount, _ := strconv.ParseFloat(row[1], 64)
			// DD-MM-YYYY
			date, _ := time.Parse(time.RFC3339, row[0])
			m.CreateInvoice(date, amount, row[2], storeId)
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
}
