package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"main/repo"
	a "main/structs/account"
	"math/big"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strings"

	"github.com/ajandera/sp_model/rdbsClientInfo"

	"github.com/bitly/go-simplejson"
)

// character sets for password generator
var (
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*"
	numberSet      = "0123456789"
	allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
)

// tokenBody struct to handle refresh token
type tokenReqBody struct {
	RefreshToken string `json:"refresh_token"`
}

func Auth(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	// Declare a new Visitor struct.
	var account a.Account

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print(err.Error())
		return
	}

	if len(account.Email) == 0 || len(account.Password) == 0 {
		WriteJSONError(w, http.StatusUnauthorized, "Email and password not match.")
		return
	}

	accountValid := m.Auth(account.Email, account.Password)
	if !IsValidUUID(accountValid.Id.String()) {
		WriteJSONError(w, http.StatusUnauthorized, "Email and password not match.")
		return
	}

	token, err := GenerateJWT(accountValid.Name, accountValid.Id.String())
	if err != nil {
		log.Print(err)
		WriteJSONError(w, http.StatusInternalServerError, "Unable to generate token.")
		return
	}

	response := simplejson.New()
	response.Set("success", true)
	response.Set("jwt", token)
	writeJSON(w, http.StatusOK, response)
}

// WriteJSONError writes {"success": false, "error": message} with the given status code.
func WriteJSONError(w http.ResponseWriter, status int, message string) {
	response := simplejson.New()
	response.Set("success", false)
	response.Set("error", message)
	writeJSON(w, status, response)
}

func writeJSON(w http.ResponseWriter, status int, response *simplejson.Json) {
	payload, err := response.MarshalJSON()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(payload)
}

// PublicAccount removes secrets (password hash, restore token) from the account
// before it is sent to the client.
func PublicAccount(account rdbsClientInfo.Accounts) rdbsClientInfo.Accounts {
	account.Password = ""
	account.RestoreToken = ""
	return account
}

// PublicAccounts removes secrets from list of accounts.
func PublicAccounts(accounts []rdbsClientInfo.Accounts) []rdbsClientInfo.Accounts {
	result := make([]rdbsClientInfo.Accounts, 0, len(accounts))
	for _, account := range accounts {
		result = append(result, PublicAccount(account))
	}
	return result
}

func Logout(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	response := simplejson.New()
	response.Set("success", true)

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Print(err.Error())

	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func Refresh(w http.ResponseWriter, r *http.Request, m repo.Repository) {
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	var tokenReq = tokenReqBody{}
	err := json.NewDecoder(r.Body).Decode(&tokenReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print(err)
		return
	}

	// only a valid, signed, non-expired refresh token is accepted
	accountId, err := ParseToken(tokenReq.RefreshToken, TokenTypeRefresh)
	if err != nil {
		WriteJSONError(w, http.StatusUnauthorized, "Refresh token is not valid.")
		return
	}

	// the account must still exist
	account := m.GetAccountById(accountId)
	if !IsValidUUID(account.Id.String()) {
		WriteJSONError(w, http.StatusUnauthorized, "Refresh token is not valid.")
		return
	}

	newTokenPair, err := GenerateJWT(account.Name, account.Id.String())
	if err != nil {
		log.Print(err)
		WriteJSONError(w, http.StatusInternalServerError, "Unable to generate token.")
		return
	}

	response := simplejson.New()
	response.Set("success", true)
	response.Set("account", PublicAccount(account))
	response.Set("jwt", newTokenPair)
	writeJSON(w, http.StatusOK, response)
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

// secureIntn returns a cryptographically secure random number in [0, n).
func secureIntn(n int) int {
	v, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		panic(err) // crypto/rand never fails on supported platforms
	}
	return int(v.Int64())
}

// RandomToken returns URL-safe random token with given number of random bytes (hex encoded).
func RandomToken(bytesLength int) string {
	b := make([]byte, bytesLength)
	if _, err := rand.Read(b); err != nil {
		panic(err) // crypto/rand never fails on supported platforms
	}
	return fmt.Sprintf("%x", b)
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[secureIntn(len(charset))]
	}
	return string(b)
}

func SetupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func GeneratePassword(passwordLength, minSpecialChar, minNum, minUpperCase int) string {
	var password strings.Builder

	//Set special character
	for i := 0; i < minSpecialChar; i++ {
		random := secureIntn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := secureIntn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := secureIntn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := secureIntn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	// Fisher-Yates shuffle
	for i := len(inRune) - 1; i > 0; i-- {
		j := secureIntn(i + 1)
		inRune[i], inRune[j] = inRune[j], inRune[i]
	}
	return string(inRune)
}

func SendEmailWithTemplate(email string, subject string, templateName string, token string, attachments []string) {
	// Sender data.
	from := os.Getenv("SMTP_USERNAME") // "storepredictor@storepredictor.com"

	// prepare template
	t, err := template.ParseFiles("templates/" + templateName + ".html")
	if err != nil {
		log.Println("email template:", err)
		return
	}
	var body bytes.Buffer

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = email
	headers["Subject"] = subject

	for k, v := range headers {
		body.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
	}

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte("Subject: " + subject + " \n" + mimeHeaders + "\n\n"))

	err = t.Execute(&body, struct {
		Email string
		Link  string
	}{
		Email: email,
		Link:  "https://my.storepredictor.com/restore?token=" + token,
	})
	if err != nil {
		log.Println("email template:", err)
		return
	}

	// send attachment if needed
	for _, file := range attachments {
		data, err := os.ReadFile(file)
		if err != nil {
			log.Println("email attachment:", err)
			return
		}
		b := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
		base64.StdEncoding.Encode(b, data)
		body.Write(b)
	}

	deliverEmail(email, body.Bytes())
}

func SendEmailWithoutTemplate(email string, subject string, htmlString string) {
	// Sender data.
	from := os.Getenv("SMTP_USERNAME") // "storepredictor@storepredictor.com"

	var body bytes.Buffer

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = email
	headers["Subject"] = subject

	for k, v := range headers {
		body.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
	}

	sub := "Subject: " + subject + "\n"
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	content := "<html><body>" + htmlString + "</body></html>"
	body.Write([]byte(sub + mimeHeaders + content))

	deliverEmail(email, body.Bytes())
}

// deliverEmail sends the prepared message. Errors are logged and never stop the API.
// SMTP_URL=mailhog uses a plain local SMTP server (development), otherwise an
// implicit TLS connection (port 465) with PLAIN auth is used.
func deliverEmail(recipient string, message []byte) {
	smtpHost := os.Getenv("SMTP_URL")  // "smtp.seznam.cz"
	smtpPort := os.Getenv("SMTP_PORT") // "465"
	from := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" {
		log.Println("SMTP_URL is not set, email to", recipient, "not sent")
		return
	}

	var c *smtp.Client
	var err error
	if smtpHost == "mailhog" {
		// Connect to the local development SMTP server.
		c, err = smtp.Dial(smtpHost + ":" + smtpPort)
		if err != nil {
			log.Println("smtp:", err)
			return
		}
		if from == "" {
			from = "noreply@storepredictor.com"
		}
	} else {
		// TLS config
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         smtpHost,
		}

		// Here is the key, you need to call tls.Dial instead of smtp.Dial
		// for smtp servers running on 465 that require an ssl connection
		// from the very beginning (no starttls)
		conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsconfig)
		if err != nil {
			log.Println("smtp:", err)
			return
		}
		c, err = smtp.NewClient(conn, smtpHost)
		if err != nil {
			log.Println("smtp:", err)
			conn.Close()
			return
		}
		// Auth
		if err = c.Auth(smtp.PlainAuth("", from, password, smtpHost)); err != nil {
			log.Println("smtp:", err)
			c.Close()
			return
		}
	}
	defer c.Close()

	// To && From
	if err = c.Mail(from); err != nil {
		log.Println("smtp:", err)
		return
	}
	if err = c.Rcpt(recipient); err != nil {
		log.Println("smtp:", err)
		return
	}

	// Data
	wr, err := c.Data()
	if err != nil {
		log.Println("smtp:", err)
		return
	}
	if _, err = wr.Write(message); err != nil {
		log.Println("smtp:", err)
		return
	}
	if err = wr.Close(); err != nil {
		log.Println("smtp:", err)
		return
	}
	if err = c.Quit(); err != nil {
		log.Println("smtp:", err)
	}
}

// ReadUserIP function to read User IP from request
func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
