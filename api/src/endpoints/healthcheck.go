package endpoints

import (
	"log"
	"main/utils"
	"net/http"

	"github.com/bitly/go-simplejson"
)

// GET
// @tags Healthcheck
// @Summary Endpoint to check function
// @Description Endpoint to check function
// @Accept  json
// @Produce  json
// @Success 200 {object} utils.Res
// @Router /healthcheck  [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.SetupCORS(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

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
