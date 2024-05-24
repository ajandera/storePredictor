package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	structs "main/structs/share"
	"main/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	model "github.com/ajandera/sp_model"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"github.com/xuri/excelize/v2"
)

const (
	SheetNameGlobal   = "Global prediction"
	SheetNameProducts = "Prediction per products"
)

// ReportData struct to get reports data
type ReportData struct {
	Name     string
	Products []ReportProducts
}

type ReportProducts struct {
	Name        string
	Code        string
	Expected    int
	DateToOrder time.Time
	DateToNeed  time.Time
}

type Reports struct {
	Success bool
	Report  string
}

// InfluxResponseProduct struct to store response from influx client
type InfluxResponseProduct struct {
	Index string
	Val   string
	Date  string
	Code  string
}

// PDFService struct to generate pdf report
type PDFService struct{}

// NewPDFService function to handle pdf report creation
func NewPDFService() *PDFService {
	return &PDFService{}
}

// generatePDF function to generate pdf
func generatePDF(data *ReportData, lang string, filename string) error {
	var templ *template.Template
	var err error

	// use Go's default HTML template generation tools to generate your HTML
	exePath, _ := filepath.Abs("./")
	if templ, err = template.ParseFiles(exePath + "/templates/" + lang + "/pdfReport.html"); err != nil {
		log.Printf(err.Error())
		return err
	}

	// apply the parsed HTML template data and keep the result in a Buffer
	var body bytes.Buffer
	if err = templ.Execute(&body, data); err != nil {
		log.Printf(err.Error())
		return err
	}

	// initalize a wkhtmltopdf generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	// read the HTML page as a PDF page
	page := wkhtmltopdf.NewPageReader(strings.NewReader(body.String()))

	// enable this if the HTML file contains local references such as images, CSS, etc.
	page.EnableLocalFileAccess.Set(true)

	// add the page to your generator
	pdfg.AddPage(page)

	// manipulate page attributes as needed
	pdfg.MarginLeft.Set(10)
	pdfg.MarginRight.Set(10)
	pdfg.Dpi.Set(600)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)

	// magic
	err = pdfg.Create()
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	name := exePath + "/temp/" + filename
	errWrite := pdfg.WriteFile(name)
	if errWrite != nil {
		log.Printf(err.Error())
	}
	return nil
}

// GetPDF function generate pdf report

// GET
// @tags Reports
// @Summary Endpoint to generate pdf report
// @Description Endpoint to generate pdf report
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.ResReport
// @Router /reports/pdf/{storeId}/{date}/{lang}  [get]
func GetPDF(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	date := vars["date"]
	lang := vars["lang"]

	if isAuthorized(w, r, m, storeId) == true {
		response := simplejson.New()
		response.Set("success", true)
		filename := generatePDFReport(date, storeId, lang, m)
		response.Set("report", filename)
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

// generatePDFReport function to generate pdf report
func generatePDFReport(date string, storeId string, lang string, m model.Repository) string {
	// prepare data
	parsedDate, _ := time.Parse(time.RFC3339, date)
	var data ReportData
	store := m.GetStoreById(storeId)
	data.Name = store.Url

	for _, p := range m.GetProductsToStore(storeId, 10000, 0) {
		if parsedDate.Before(p.DateToNeed) {
			var row ReportProducts
			row.Name = p.ProductCode
			row.Code = p.ProductCode
			row.Expected = int(p.Quantity)
			row.DateToOrder = p.DateToOrder
			row.DateToNeed = p.DateToNeed
			data.Products = append(data.Products, row)
		}
	}

	filename := store.Url + "_report_" + date + ".pdf"
	err := generatePDF(&data, lang, filename)
	if err != nil {
		log.Printf(err.Error())
		return ""
	}
	return filename
}

// GetExcel function to generate excel report

// GET
// @tags Reports
// @Summary Endpoint to generate excel report
// @Description Endpoint generate excel report
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.ResReport
// @Router /reports/excel/{storeId}/{date}  [get]
func GetExcel(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	date := vars["date"]
	if isAuthorized(w, r, m, storeId) == true {
		response := simplejson.New()
		response.Set("success", true)
		filename := GenerateExcel(storeId, date, m)
		response.Set("report", filename)
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

// GenerateExcel function to generate excel
func GenerateExcel(storeId string, date string, m model.Repository) string {
	now := time.Now()
	f := excelize.NewFile()
	store := m.GetStoreById(storeId)
	filename := store.Url + "_report_" + now.Format("2006-01-02") + ".xlsx"

	f.SetCellValue("Sheet1", "A1", store.Url)
	f.SetCellValue("Sheet1", "A2", "Created from storepredictor.com")

	f.SetCellValue("Sheet1", "A4", "Product Code")
	f.SetCellValue("Sheet1", "B4", "Expected")
	f.SetCellValue("Sheet1", "C4", "Date To Order")
	f.SetCellValue("Sheet1", "D4", "Date To Need")
	check, _ := time.Parse(time.RFC3339, date)
	for i, p := range m.GetProductsToStore(storeId, 10000, 0) {
		if check.Before(p.DateToNeed) {
			f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+4), p.ProductCode)
			f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+4), p.Quantity)
			f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+4), p.DateToOrder.Format(time.RFC3339))
			f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+4), p.DateToNeed.Format(time.RFC3339))
		}
	}

	exePath, _ := filepath.Abs("./")
	filePath := exePath + "/temp/" + filename
	if err := f.SaveAs(filePath); err != nil {
		log.Println(err)
	}

	return filename
}

// GET Reports
// @tags Reports
// @Summary Endpoint to generate full excel report
// @Description Endpoint generate full excel report
// @Accept  json
// @Produce  json
// @Success 200 {object} Reports
// @Router /reports/excelfull/{storeId}  [get]
func GetExcelFull(w http.ResponseWriter, r *http.Request, m model.Repository, i model.Influx) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	//if isAuthorized(w, r, m, storeId) == true {
	response := simplejson.New()
	response.Set("success", true)
	filename := GenerateExcelFull(storeId, m, i)
	response.Set("report", filename)
	payload, err := response.MarshalJSON()
	if err != nil {
		log.Printf(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)

	//}
}

// GenerateExcelFull function to generate full excel with predicted data
func GenerateExcelFull(storeId string, m model.Repository, i model.Influx) string {
	now := time.Now()
	f := excelize.NewFile()
	index, _ := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)
	f.SetSheetName("Sheet1", SheetNameGlobal)

	store := m.GetStoreById(storeId)
	filename := store.Url + "_report_full_" + now.Format("2006-01-02") + ".xlsx"
	query := `from(bucket:"` + storeId + `") 
				|> range(start:-0mo, stop:2mo)  
				|> filter(fn: (r) => (r._measurement == "order" and (r._field == "value")))`
	data, errInf := i.GetInfluxQuery(query, os.Getenv("INFLUX_ORGANIZATION"))
	if errInf != nil {
		log.Printf(errInf.Error())
		return ""
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
	f.SetCellValue(SheetNameGlobal, "A1", store.Url)
	f.SetCellValue(SheetNameGlobal, "A2", "Created from storepredictor.com")
	f.SetCellValue(SheetNameGlobal, "A4", "Date")
	f.SetCellValue(SheetNameGlobal, "B4", "Predicted orders")

	for i, p := range influxResponse {
		f.SetCellValue(SheetNameGlobal, "A"+strconv.Itoa(i+4), p.Date)
		f.SetCellValue(SheetNameGlobal, "B"+strconv.Itoa(i+4), p.Val)
	}

	index2, _ := f.NewSheet("Sheet2")
	f.SetActiveSheet(index2)
	f.SetSheetName("Sheet2", SheetNameProducts)

	// per products
	iIdx := 1
	for _, prod := range m.GetProducts(map[string]interface{}{"store_id": storeId, "limit": 100000, "offset": 0}) {
		// load data from influx
		queryProducts := `from(bucket:"` + storeId + `") 
				|> range(start:-0mo, stop:2mo)  
				|> filter(fn: (r) => (r._measurement == "` + prod.ProductCode + `" and (r._field == "value")))`
		dataProducts, errInfProd := i.GetInfluxQuery(queryProducts, os.Getenv("INFLUX_ORGANIZATION"))
		if errInfProd != nil {
			log.Printf(errInf.Error())
			return ""
		}

		// Iterate over query response
		var influxResponseProd []InfluxResponseProduct
		for dataProducts.Next() {
			// Notice when group key has changed
			if dataProducts.TableChanged() {
				//fmt.Printf("table: %s\n", data.TableMetadata().String())
			}
			// Access data and store  d03
			if fmt.Sprint(dataProducts.Record().ValueByKey("daysToMeasurement")) == "d03" {
				log.Printf(prod.ProductCode)
				influxResponseProd = append(influxResponseProd, InfluxResponseProduct{
					fmt.Sprint(dataProducts.Record().ValueByKey("daysToMeasurement")),
					fmt.Sprint(dataProducts.Record().Value()),
					dataProducts.Record().Time().String(),
					dataProducts.Record().Measurement()})
			}
		}
		jIdx := 1
		for _, p := range influxResponseProd {
			idx, _ := excelize.CoordinatesToCellName(jIdx, iIdx)
			if jIdx == 1 {
				f.SetCellValue(SheetNameProducts, idx, p.Code)
			} else {
				f.SetCellValue(SheetNameProducts, idx, p.Date)
				iIdx++
				idxVal, _ := excelize.CoordinatesToCellName(jIdx, iIdx)
				f.SetCellValue(SheetNameProducts, idxVal, p.Index+" - "+p.Val)
			}
			jIdx++
		}
		iIdx++
	}

	exePath, _ := filepath.Abs("./")
	filePath := exePath + "/temp/" + filename
	if err := f.SaveAs(filePath); err != nil {
		log.Println(err)
	}

	return filename
}

// ShareReport function to handle API requests to logged user

// GET
// @tags Reports
// @Summary Endpoint to share reposrt
// @Description Endpoint to share reports
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /reports/share/{storeId}/{date}/{lang}  [get]
func ShareReport(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	date := vars["date"]
	lang := vars["lang"]
	exePath, _ := filepath.Abs("./")
	store := m.GetStoreById(storeId)

	if isAuthorized(w, r, m, storeId) == true {
		// Declare a new Share struct.
		var share structs.Share

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&share)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf(err.Error())
			return
		}

		var attachments []string
		if share.Excel == true {
			filename := GenerateExcel(storeId, date, m)
			attachments = append(attachments, exePath+"/temp/"+filename)
		}

		if share.Pdf == true {
			pdfFile := generatePDFReport(date, storeId, lang, m)
			attachments = append(attachments, exePath+"/temp/"+pdfFile)
		}

		utils.SendEmailWithTemplate(share.Email, "Report "+store.Url+" "+date, lang+"/shareReport", "", attachments)

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
