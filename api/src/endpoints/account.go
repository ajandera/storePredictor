package endpoints

import (
	"encoding/json"
	"fmt"
	"log"
	a "main/structs/account"
	structs "main/structs/pw"
	"main/utils"
	"math/rand"
	"net/http"
	"strings"
	"time"

	model "github.com/ajandera/sp_model"
	"github.com/ajandera/sp_model/rdbsClientInfo"

	"github.com/bitly/go-simplejson"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

var mySigningKey = []byte("DFGDFGhcsadkjhfwe+ě+23123asldxjhsdljfh1234234")

// AccountMe function to handle API requests to logged user

// GET
// @tags Account
// @Summary Endpoint to send logged user account
// @Description Endpoint to send logged user account
// @Accept  json
// @Produce  json
// @Success 200 {object} a.Account
// @Router /account/me  [get]
func AccountMe(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if isAuthorized(w, r, m, "") == true {
		var sendToken = strings.Replace(r.Header["Authorization"][0], "Bearer ", "", 1)
		token, err := jwt.Parse(sendToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("JWT token not pass")
			}
			return mySigningKey, nil
		})
		claims := token.Claims.(jwt.MapClaims)
		response := simplejson.New()
		response.Set("success", true)
		response.Set("user", m.GetAccountById(fmt.Sprintf("%v", claims["id"])))

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

// GetAccounts function to handle API requests to get accounts from database

// GET
// @tags Account
// @Summary Endpoint to send all accounts
// @Description Endpoint to send all accounts
// @Accept  json
// @Produce  json
// @Success 200 {object} []a.Account
// @Router /accounts  [get]
func GetAccounts(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if isAuthorized(w, r, m, "") == true {
		response := simplejson.New()
		response.Set("success", true)
		response.Set("accounts", m.GetAccounts())

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

// SendNewPw function to handle API requests to forgot password

// GET
// @tags Password
// @Summary Endpoint to send new passwrod
// @Description Endpoint to send new password
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /forgot/{accountId}/{lang}  [get]
func SendNewPw(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	rand.Seed(time.Now().Unix())
	minSpecialChar := 1
	minNum := 1
	minUpperCase := 1
	passwordLength := 8
	token := utils.GeneratePassword(passwordLength, minSpecialChar, minNum, minUpperCase)

	vars := mux.Vars(r)
	accountId := vars["accountId"]
	lang := vars["lang"]

	account := m.GetAccountById(accountId)
	m.SetRestorePw(
		account.Id.String(),
		token)

	var attachments []string
	utils.SendEmailWithTemplate(account.Email, "Recover your password", lang+"/forgot", token, attachments)

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

// Accounts function to handle API requests to store account in database

// POST
// @tags Account
// @Summary Endpoint to store account data
// @Description Endpoint to store account data
// @Param  a.Account body  a.Account true  "Account object to store in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /accounts/{lang}  [post]
func Accounts(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	lang := vars["lang"]

	// Declare a new Account struct.
	var account a.Account

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil || len(account.Email) == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf(err.Error())
		return
	}

	response := simplejson.New()
	check := m.GetAccountByEmail(account.Email)
	if utils.IsValidUUID(check.Id.String()) {
		response.Set("success", false)
		response.Set("message", "Email already registered")
		payload, err := response.MarshalJSON()
		if err != nil {
			log.Printf(err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)

	} else {
		// register account
		var newAccount rdbsClientInfo.Accounts
		m.CreateAccount(account.Email, account.Password, false).Scan(&newAccount)
		// fill all
		m.EditAccount(newAccount.Id.String(), account.Name, account.Email, account.Street, account.City, account.Zip, account.CountryCode,
			account.CompanyNumber, account.VatNumber, account.PaidTo, account.PlanRefer, account.Role,
			account.Parent, account.Password, false)
		response.Set("success", true)

		var attachments []string
		utils.SendEmailWithTemplate(account.Email, "Welcome at storePredictor", lang+"/register", "", attachments)

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Fatalf(err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	}
}

// UpdateAccounts function to handle API requests to update account in database

// PUT
// @tags Account
// @Summary Endpoint to store account data
// @Description Endpoint to update account data
// @Param  a.Account body  a.Account true  "Account object to update in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /accounts  [put]
func UpdateAccounts(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if isAuthorized(w, r, m, "") == true {
		// Declare a new Visitor struct.
		var account a.Account

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&account)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf(err.Error())
			return
		}

		response := simplejson.New()
		m.EditAccount(account.ID, account.Name, account.Email, account.Street, account.City, account.Zip, account.CountryCode,
			account.CompanyNumber, account.VatNumber, account.PaidTo, account.PlanRefer, account.Role,
			account.Parent, account.Password, account.Newsletter)
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

// DeleteAccounts function to handle API requests to remove account from database

// DELETE
// @tags Account
// @Summary Endpoint to delete account
// @Description Endpoint to delete account
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /accounts/{accountId}  [delete]
func DeleteAccounts(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if isAuthorized(w, r, m, "") == true {
		vars := mux.Vars(r)
		accountId := vars["accountId"]

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		if len(accountId) <= 0 {
			http.Error(w, "Account id not  send.", http.StatusBadRequest)
			return
		}

		response := simplejson.New()
		m.DeleteAccount(accountId)
		response.Set("success", true)

		payload, err := response.MarshalJSON()
		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}

// UpdatePassword function to handle API requests to update password in database

// PUT
// @tags Password
// @Summary Endpoint to update password
// @Description Endpoint to update password
// @Param  structs.Pw body  structs.Pw true  "Pw object to update in database"
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /restore/{token}  [put]
func UpdatePassword(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	refreshToken := vars["token"]

	// Declare a new Password struct.
	var t structs.Pw

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf(err.Error())
		return
	}

	response := simplejson.New()
	m.UpdatePw(refreshToken, t.Password)
	response.Set("success", true)

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Printf(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(payload)
}

// GetAccountOrders function to handle API requests to get orders per account from database

// GET
// @tags Account
// @Summary Endpoint to get account orders
// @Description Endpoint to get account orders
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /account/order/{accountId}/{storeId}  [get]
func GetAccountOrders(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	accountId := vars["accountId"]

	if isAuthorized(w, r, m, storeId) == true {
		response := simplejson.New()
		response.Set("success", true)
		response.Set("orders", m.GetAccountOrders(accountId))

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

// GetAccountByEmail function to handle API requests to get account by email

// GET
// @tags Account
// @Summary Endpoint to get account by email
// @Description Endpoint to get account by email
// @Accept  json
// @Produce  json
// @Success 200 {object} a.Account
// @Router /account/email/{email}  [get]
func GetAccountByEmail(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	email := vars["email"]

	response := simplejson.New()
	response.Set("success", true)
	response.Set("account", m.GetAccountByEmail(email))

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Printf(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// isAuthorized function to check jwt token in request
func isAuthorized(w http.ResponseWriter, r *http.Request, m model.Repository, storeId string) bool {
	if r.Header["Authorization"] != nil {
		var sendToken = strings.Replace(r.Header["Authorization"][0], "Bearer ", "", 1)
		token, err := jwt.Parse(sendToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("JWT token not pass")
			}
			return mySigningKey, nil
		})

		// check if store belongs to account
		if utils.IsValidUUID(storeId) {
			claims := token.Claims.(jwt.MapClaims)
			if m.IsPermitted(fmt.Sprintf("%v", claims["id"]), storeId) == false {
				http.Error(w, "Not Permitted", http.StatusForbidden)
				return false
			}
		}

		if err != nil && token != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf(err.Error())
			return false
		}
		return true
	} else {
		http.Error(w, "Not Authorized", http.StatusForbidden)
		return false
	}
}

// GET
// @tags Account
// @Summary Endpoint to get account
// @Description Endpoint to get account
// @Accept  json
// @Produce  json
// @Success 200 {object} structs.Account
// @Router /accounts/{accountId}  [get]
func GetAccount(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	accountId := vars["accountId"]

	if isAuthorized(w, r, m, "") == true {
		response := simplejson.New()
		response.Set("success", true)
		response.Set("account", m.GetAccountById(accountId))

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

// GetChildAccount function to handle API requests to get children account for parent

// GET
// @tags Account
// @Summary Endpoint to get children account for parent
// @Description Endpoint to get children account for parent
// @Accept  json
// @Produce  json
// @Success 200 {object} structs.Account
// @Router /account/child/{id}  [get]
func GetChildAccount(w http.ResponseWriter, r *http.Request, m model.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	response := simplejson.New()
	response.Set("success", true)
	response.Set("accounts", m.GetChildAccountById(id))

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Printf(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
