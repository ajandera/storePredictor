package endpoints

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"main/utils"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/google/uuid"
)

// Token function to generate token for tracking without cookies

// GET
// @tags Token
// @Summary Endpoint to generate unique token for tracking
// @Description Endpoint to generate unique token for tracking
// @Accept  json
// @Produce  json
// @Success 200 {object} string "OK"
// @Router /token  [get]
func Token(w http.ResponseWriter, r *http.Request) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	var requestTime int64
	var requestTTL int64

	requestTime, _ = strconv.ParseInt(r.Header.Get("X-Request-Time"), 10, 32)
	if requestTime == 0 {
		requestTime = time.Now().Unix()
	}

	requestTTL, _ = strconv.ParseInt(r.Header.Get("X-Time-To-Live"), 10, 32)
	if requestTTL == 0 {
		requestTTL = 24 * 3600
	}

	uniqueCode := uuid.New().String()

	response := simplejson.New()
	response.Set("time", requestTime)
	response.Set("expire", requestTTL)
	response.Set("code", uniqueCode)

	payload, err := response.MarshalJSON()
	if err != nil {
		log.Println(err.Error())
	}

	etag := md5.Sum([]byte(strconv.Itoa(int(requestTime))))
	diff := math.Max(float64((requestTime+requestTTL)-time.Now().Unix()), 0)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "max-age="+strconv.Itoa(int(diff)))
	w.Header().Set("ETag", hex.EncodeToString(etag[:]))
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
