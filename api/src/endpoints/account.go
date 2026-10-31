package endpoints

import (
	"encoding/json"
	"log"
	"main/repo"
	a "main/structs/account"
	structs "main/structs/pw"
	"main/utils"
	"net/http"

	"github.com/ajandera/sp_model/rdbsClientInfo"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

// AccountMe function to handle API requests to logged user

// GET
// @tags Account
// @Summary Endpoint to send logged user account
// @Description Endpoint to send logged user account
// @Accept  json
// @Produce  json
// @Success 200 {object} a.Account
// @Router /account/me  [get]
func AccountMe(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if accountId, ok := authenticate(w, r); ok {
		response := simplejson.New()
		response.Set("success", true)
		response.Set("user", utils.PublicAccount(m.GetAccountById(accountId)))

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

// GetAccounts function to handle API requests to get accounts from database

// GET
// @tags Account
// @Summary Endpoint to send all accounts
// @Description Endpoint to send all accounts
// @Accept  json
// @Produce  json
// @Success 200 {object} []a.Account
// @Router /accounts  [get]
func GetAccounts(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if callerId, ok := authenticate(w, r); ok {
		// only accounts of the caller's tenant (main account + its child accounts)
		owner := tenantOwner(m.GetAccountById(callerId))
		accounts := append([]rdbsClientInfo.Accounts{m.GetAccountById(owner)}, m.GetChildAccountById(owner)...)

		response := simplejson.New()
		response.Set("success", true)
		response.Set("accounts", utils.PublicAccounts(accounts))

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

// SendNewPw function to handle API requests to forgot password

// GET
// @tags Password
// @Summary Endpoint to send new passwrod
// @Description Endpoint to send new password
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /forgot/{accountId}/{lang}  [get]
func SendNewPw(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	accountId := vars["accountId"]
	lang := vars["lang"]

	account := m.GetAccountById(accountId)
	if utils.IsValidUUID(account.Id.String()) {
		// signed, URL-safe restore token valid for 24 hours
		token, err := utils.GenerateRestoreToken(account.Id.String())
		if err != nil {
			log.Print(err)
			utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to generate token.")
			return
		}
		m.SetRestorePw(
			account.Id.String(),
			token)

		var attachments []string
		utils.SendEmailWithTemplate(account.Email, "Recover your password", lang+"/forgot", token, attachments)
	}

	response := simplejson.New()
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
func Accounts(w http.ResponseWriter, r *http.Request, m repo.Repository) {
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print(err)
		return
	}
	if len(account.Email) == 0 || len(account.Password) <= 6 {
		http.Error(w, "Email and password (min. 7 characters) are required.", http.StatusBadRequest)
		return
	}
	// public registration always creates a main account: the parent, role and plan
	// can't be chosen by the client (otherwise anybody could join a foreign tenant)
	account.Parent = ""
	account.Role = ""
	account.PaidTo = ""
	account.PlanRefer = ""

	response := simplejson.New()
	check := m.GetAccountByEmail(account.Email)
	if utils.IsValidUUID(check.Id.String()) {
		response.Set("success", false)
		response.Set("message", "Email already registered")
		payload, err := response.MarshalJSON()
		if err != nil {
			log.Print(err.Error())
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
			log.Print(err)
			return
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
func UpdateAccounts(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if callerId, ok := authenticate(w, r); ok {
		// Declare a new Visitor struct.
		var account a.Account

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&account)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Print(err)
			return
		}

		// plan and payment are never changed through this endpoint
		account.PaidTo = ""
		account.PlanRefer = ""

		switch {
		case account.ID == "":
			// new child account (user) of the logged main account
			caller := m.GetAccountById(callerId)
			if caller.Parent != "" {
				denyAccess(w)
				return
			}
			if len(account.Email) == 0 || len(account.Password) <= 6 {
				http.Error(w, "Email and password (min. 7 characters) are required.", http.StatusBadRequest)
				return
			}
			if utils.IsValidUUID(m.GetAccountByEmail(account.Email).Id.String()) {
				utils.WriteJSONError(w, http.StatusConflict, "Email already registered")
				return
			}
			account.Parent = callerId
		case account.ID == callerId:
			// own profile: role and parent can't be changed by the account itself
			account.Role = ""
			account.Parent = ""
		case canManageAccount(m, callerId, account.ID):
			// child account managed by the caller, it can't be moved to other parent
			account.Parent = ""
		default:
			denyAccess(w)
			return
		}

		response := simplejson.New()
		m.EditAccount(account.ID, account.Name, account.Email, account.Street, account.City, account.Zip, account.CountryCode,
			account.CompanyNumber, account.VatNumber, account.PaidTo, account.PlanRefer, account.Role,
			account.Parent, account.Password, account.Newsletter)
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

// DeleteAccounts function to handle API requests to remove account from database

// DELETE
// @tags Account
// @Summary Endpoint to delete account
// @Description Endpoint to delete account
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /accounts/{accountId}  [delete]
func DeleteAccounts(w http.ResponseWriter, r *http.Request, m repo.Repository) {
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
			http.Error(w, "Account id not  send.", http.StatusBadRequest)
			return
		}
		if !canManageAccount(m, callerId, accountId) {
			denyAccess(w)
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
func UpdatePassword(w http.ResponseWriter, r *http.Request, m repo.Repository) {
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
		log.Print(err)
		return
	}

	// the token must be a valid (signed, not expired) restore token and it must be
	// the last token issued for the account (it is removed after the first use)
	accountId, err := utils.ParseToken(refreshToken, utils.TokenTypeRestore)
	if err != nil || m.GetAccountById(accountId).RestoreToken != refreshToken {
		utils.WriteJSONError(w, http.StatusBadRequest, "Restore token is not valid or expired.")
		return
	}
	if len(t.Password) <= 6 {
		utils.WriteJSONError(w, http.StatusBadRequest, "Password must have at least 7 characters.")
		return
	}

	response := simplejson.New()
	m.UpdatePw(refreshToken, t.Password)
	response.Set("success", true)

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Print(err)
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
func GetAccountOrders(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	storeId := vars["storeId"]
	accountId := vars["accountId"]

	if isAuthorized(w, r, m, storeId) == true {
		callerId, _ := utils.ParseToken(utils.BearerToken(r), utils.TokenTypeAccess)
		if !canViewAccount(m, callerId, accountId) {
			denyAccess(w)
			return
		}

		response := simplejson.New()
		response.Set("success", true)
		response.Set("orders", m.GetAccountOrders(accountId))

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

// GetAccountByEmail function to handle API requests to get account by email

// GET
// @tags Account
// @Summary Endpoint to get account by email
// @Description Endpoint to get account by email. Public endpoint, it returns only the account id
// @Description (zero UUID when email is not registered), never the account data.
// @Accept  json
// @Produce  json
// @Success 200 {object} a.Account
// @Router /account/email/{email}  [get]
func GetAccountByEmail(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	email := vars["email"]

	response := simplejson.New()
	response.Set("success", true)
	// the endpoint is public (forgot password, user existence check), expose the id only
	response.Set("account", map[string]string{"Id": m.GetAccountByEmail(email).Id.String()})

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// GET
// @tags Account
// @Summary Endpoint to get account
// @Description Endpoint to get account
// @Accept  json
// @Produce  json
// @Success 200 {object} structs.Account
// @Router /accounts/{accountId}  [get]
func GetAccount(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	accountId := vars["accountId"]

	if callerId, ok := authenticate(w, r); ok {
		if !canViewAccount(m, callerId, accountId) {
			denyAccess(w)
			return
		}

		response := simplejson.New()
		response.Set("success", true)
		response.Set("account", utils.PublicAccount(m.GetAccountById(accountId)))

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

// GetChildAccount function to handle API requests to get children account for parent

// GET
// @tags Account
// @Summary Endpoint to get children account for parent
// @Description Endpoint to get children account for parent
// @Accept  json
// @Produce  json
// @Success 200 {object} structs.Account
// @Router /account/child/{id}  [get]
func GetChildAccount(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	callerId, ok := authenticate(w, r)
	if !ok {
		return
	}
	if !canViewAccount(m, callerId, id) {
		denyAccess(w)
		return
	}

	response := simplejson.New()
	response.Set("success", true)
	response.Set("accounts", utils.PublicAccounts(m.GetChildAccountById(id)))

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// GetPlans function to handle API requests to get all plans from database
func GetPlans(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	response := simplejson.New()
	response.Set("success", true)
	response.Set("plans", m.GetPlans())

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// GetPaidPlans function to handle API requests to get paid plans from database
func GetPaidPlans(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	response := simplejson.New()
	response.Set("success", true)
	response.Set("plans", m.GetPaidPlans())

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
