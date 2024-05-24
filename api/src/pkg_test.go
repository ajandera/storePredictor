package main_test

import (
	"fmt"
	"io"
	"log"
	endpoints "main/endpoints"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		endpoints.HealthCheck(w, r)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	response, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	asrt := assert.New(t)
	r := "{\"success\":true}"
	asrt.Equal(r, string(response), "Response is not true")
}

func TestOrders(t *testing.T) {

}

func TestToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		endpoints.Token(w, r)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	etag := res.Header.Get("Etag")
	fmt.Println(etag)
	if etag == "" {
		t.Errorf("Etag is empty")
	}
	asrt := assert.New(t)
	asrt.Equal(res.StatusCode, 200, "Request was not success.")
}

func TestVisitors(t *testing.T) {

}

func TestVisitorsOffline(t *testing.T) {

}
