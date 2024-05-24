// Package api to handle api requests
package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/utils"
	"net/http"
	"net/url"
	"strings"

	model "github.com/ajandera/sp_model"
	"github.com/bitly/go-simplejson"
)

// Visitor struct to store data
type Visitor struct {
	Url         string
	Header      string
	Ip          string
	ProductCode string
	ProductName string
	Code        string
	Tag         string
}

// VisitorAlias Visitor alias strut for parse product code as float
type VisitorAlias struct {
	Url         string
	Header      string
	Ip          string
	ProductCode float64
	ProductName string
	Code        string
	Tag         string
}

// VisitorOffline struct to store data
type VisitorOffline struct {
	Info string
	Unit string
	Code string
}

// Visitors function to handle API request to store visitors in database

// POST
// @tags Visitors
// @Summary Endpoint to store visitors data for tracking
// @Description Endpoint to store visitors data for tracking
// @Param  Visitor body  Visitor true  "Visitor object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /visitors  [post]
func Visitors(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	// Declare a new Visitor struct.
	var v Visitor
	var va VisitorAlias

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	body, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		errVA := json.NewDecoder(r.Body).Decode(&va)
		if errVA != nil {
			http.Error(w, errVA.Error(), http.StatusBadRequest)
			log.Println(errVA.Error())
			return
		}
	}

	if va.Code != "" {
		v.Ip = va.Ip
		v.Code = va.Code
		v.Url = va.Url
		v.Header = va.Header
		v.ProductCode = fmt.Sprintf("%f", va.ProductCode)
		v.Tag = va.Tag
	}

	// check if SP code is valid for send domain
	u, err := url.Parse(v.Url)
	if err != nil {
		log.Println(err)
	}
	parts := strings.Split(u.Hostname(), ".")
	// for subdomain
	storeIdSubdomain := m.CheckStoreCode(v.Code, u.Hostname())
	response := simplejson.New()
	userIP := utils.ReadUserIP(r)
	if utils.IsValidUUID(storeIdSubdomain) == true {
		m.SaveVisitor(userIP, storeIdSubdomain, v.Url, v.Header, v.ProductCode, v.Tag)
		response.Set("success", true)
		payload, err := response.MarshalJSON()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	} else {
		if len(parts) > 1 {
			domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
			storeId := m.CheckStoreCode(v.Code, domain)
			if utils.IsValidUUID(storeId) == true {
				m.SaveVisitor(userIP, storeId, v.Url, v.Header, v.ProductCode, v.Tag)
				response.Set("success", true)
				w.WriteHeader(http.StatusCreated)
			} else {
				response.Set("success", false)
				w.WriteHeader(http.StatusForbidden)
				response.Set("error", "Not valid store code")
				log.Println("Bad url: " + domain)
			}

			payload, err := response.MarshalJSON()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Println(err.Error())
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(payload)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("Bad url: " + u.Hostname())
		}
	}
}

// Visitors function to handle API request to store visitors in database

// POST
// @tags Visitors Offline
// @Summary Endpoint to store visitors data for tracking
// @Description Endpoint to store visitors data for tracking
// @Param  VisitorOffline body VisitorOffline true  "VisitorOffline object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /visitors-offline  [post]
func VisitorsOffline(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	// Declare a new Visitor struct.
	var v VisitorOffline

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	body, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&v)

	// for subdomain
	storeId := m.CheckStoreCodeOffline(v.Code, v.Unit)
	response := simplejson.New()
	if utils.IsValidUUID(storeId) == true {
		m.SaveVisitorOffline(v.Info, storeId)
		response.Set("success", true)
		payload, err := response.MarshalJSON()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Bad unit: " + v.Unit)
	}
}
