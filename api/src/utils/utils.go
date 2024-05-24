package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	a "main/structs/account"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"time"

	model "github.com/ajandera/sp_model"
	"github.com/ajandera/sp_model/rdbsClientInfo"

	"github.com/bitly/go-simplejson"
	"github.com/golang-jwt/jwt/v4"
)

// variables for JWT token
var (
	mySigningKey   = []byte("DFGDFGhcsadkjhfwe+ě+23123asldxjhsdljfh1234234")
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

func Auth(w http.ResponseWriter, r *http.Request, m model.Repository) {
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
		log.Printf(err.Error())
		return
	}

	response := simplejson.New()

	accountValid := m.Auth(account.Email, account.Password)
	token, err := GenerateJWT(accountValid.Name, accountValid.Id.String())
	if IsValidUUID(accountValid.Id.String()) == true && err == nil {
		response.Set("success", true)
		response.Set("jwt", token)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
		response.Set("success", false)
		response.Set("error", "Email and password not match.")
	}

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Printf(err.Error())

	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func Logout(w http.ResponseWriter, r *http.Request, m model.Repository) {
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	response := simplejson.New()
	response.Set("success", true)

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Printf(err.Error())

	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func Refresh(w http.ResponseWriter, r *http.Request, m model.Repository) {
	SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	var tokenReq = tokenReqBody{}
	err := json.NewDecoder(r.Body).Decode(&tokenReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf(err.Error())

		return
	}

	// Parse takes the token string and a function for looking up the key.
	// The latter is especially useful if you use multiple keys for your application.
	// The standard is to use 'kid' in the head of the token to identify
	// which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenReq.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return mySigningKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		var account rdbsClientInfo.Accounts = m.GetAccountById(fmt.Sprint(claims["id"]))
		if account.Id.String() != "" {
			newTokenPair, err := GenerateJWT(account.Name, account.Id.String())
			if err != nil {

				return
			}

			response := simplejson.New()
			if IsValidUUID(account.Id.String()) == true && err == nil {
				response.Set("success", true)
				response.Set("account", account)
				response.Set("jwt", newTokenPair)
			} else {
				response.Set("success", false)
				response.Set("error", "Refresh token is not valid.")
			}

			payload, err := response.MarshalJSON()
			if err != nil {
				log.Printf(err.Error())

			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(payload)
		}
	}
}

func GenerateJWT(name string, id string) (map[string]string, error) {

	// access token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = name
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		log.Printf(err.Error())

	}

	// refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["id"] = id
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	refreshTokenString, err := refreshToken.SignedString(mySigningKey)

	if err != nil {
		log.Printf(err.Error())

	}

	return map[string]string{
		"access_token":  tokenString,
		"refresh_token": refreshTokenString,
	}, nil
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func StringWithCharset(length int, charset string) string {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
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
		random := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func SendEmailWithTemplate(email string, subject string, templateName string, token string, attachments []string) {
	smtpHost := os.Getenv("SMTP_URL")  // "smtp.seznam.cz"
	smtpPort := os.Getenv("SMTP_PORT") // "465"

	// Sender data.
	from := os.Getenv("SMTP_USERNAME")     // "storepredictor@storepredictor.com"
	password := os.Getenv("SMTP_PASSWORD") //"storepredictor123"

	// prepare template
	t, _ := template.ParseFiles("templates/" + templateName + ".html")
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
	body.Write([]byte(fmt.Sprintf("Subject: "+subject+" \n%s\n\n", mimeHeaders)))

	err := t.Execute(&body, struct {
		Email string
		Link  string
	}{
		Email: email,
		Link:  "https://my.storepredictor.com/restore?token=" + token,
	})

	// send attachment if needed
	for _, file := range attachments {
		data := readFile(file)
		b := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
		base64.StdEncoding.Encode(b, data)
		body.Write(b)
	}

	if err != nil {
		log.Printf(err.Error())

		return
	} else {
		if smtpHost == "mailhog" {
			// Connect to the remote SMTP server.
			c, err := smtp.Dial(smtpHost + ":" + smtpPort)
			if err != nil {
				log.Fatal(err)
			}
			defer c.Close()
			// Set the sender and recipient.
			c.Mail("email")
			c.Rcpt("test@storepredictor.com")
			// Send the email body.
			wc, err := c.Data()
			if err != nil {
				log.Fatal(err)
			}
			defer wc.Close()
			if _, err = body.WriteTo(wc); err != nil {
				log.Fatal(err)
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
			servername := smtpHost + ":" + smtpPort
			conn, err := tls.Dial("tcp", servername, tlsconfig)
			if err != nil {
				log.Println(err)

			}

			c, err := smtp.NewClient(conn, smtpHost)
			if err != nil {
				log.Println(err)

			}

			// Auth
			auth := smtp.PlainAuth("", from, password, smtpHost)
			if err = c.Auth(auth); err != nil {
				log.Println(err)

			}

			// To && From
			if err = c.Mail(from); err != nil {
				log.Println(err)

			}

			if err = c.Rcpt(email); err != nil {
				log.Println(err)

			}

			// Data
			wr, err := c.Data()
			if err != nil {
				log.Println(err)

			}

			_, err = wr.Write(body.Bytes())
			if err != nil {
				log.Println(err)

			}

			err = wr.Close()
			if err != nil {
				log.Println(err)

			}
			errC := c.Quit()
			if errC != nil {
				log.Println(err.Error())

			}
		}
	}
}

func readFile(fileName string) []byte {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func SendEmailWithoutTemplate(email string, subject string, htmlString string) {
	// smtp server configuration.
	smtpHost := os.Getenv("SMTP_URL")  // "smtp.seznam.cz"
	smtpPort := os.Getenv("SMTP_PORT") // "465"

	// Sender data.
	from := os.Getenv("SMTP_USERNAME")     // "storepredictor@storepredictor.com"
	password := os.Getenv("SMTP_PASSWORD") //"storepredictor123"

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

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	servername := smtpHost + ":" + smtpPort
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Println(err)

	}

	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Println(err)

	}

	// Auth
	auth := smtp.PlainAuth("", from, password, smtpHost)
	if err = c.Auth(auth); err != nil {
		log.Println(err)

	}

	// To && From
	if err = c.Mail(from); err != nil {
		log.Println(err)

	}

	if err = c.Rcpt(email); err != nil {
		log.Println(err)

	}

	// Data
	wr, err := c.Data()
	if err != nil {
		log.Println(err)

	}

	_, err = wr.Write(body.Bytes())
	if err != nil {
		log.Println(err)

	}

	err = wr.Close()
	if err != nil {
		log.Println(err)

	}
	errC := c.Quit()
	if errC != nil {
		log.Println(err.Error())

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
