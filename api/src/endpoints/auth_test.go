package endpoints_test

import (
	"encoding/json"
	"main/endpoints"
	"main/repo"
	"main/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testSecret = "super-secret-test-jwt-key-with-sufficient-length-32-chars"

type handler func(http.ResponseWriter, *http.Request, repo.Repository)

// protected lists handlers that must reject a request before touching the database.
var protected = map[string]handler{
	"AccountMe":   endpoints.AccountMe,
	"GetAccounts": endpoints.GetAccounts,
}

func do(h handler, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	// fake repository: handlers must return before using it
	h(rec, req, newFakeRepo())
	return rec
}

func TestProtectedEndpointsRejectMissingToken(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	for name, h := range protected {
		t.Run(name, func(t *testing.T) {
			rec := do(h, httptest.NewRequest(http.MethodGet, "/", nil))
			assert.Equal(t, http.StatusUnauthorized, rec.Code)

			var body map[string]interface{}
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, false, body["success"])
			assert.Equal(t, "Not Authorized", body["error"])
		})
	}
}

func TestProtectedEndpointsRejectInvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	tokens, err := utils.GenerateJWT("Test", "11111111-1111-4111-8111-111111111111")
	assert.NoError(t, err)

	cases := map[string]string{
		"garbage":       "Bearer not-a-token",
		"refresh token": "Bearer " + tokens["refresh_token"], // wrong token type
		"no bearer":     tokens["access_token"],
	}
	for name, h := range protected {
		for cname, header := range cases {
			t.Run(name+"/"+cname, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", header)
				assert.Equal(t, http.StatusUnauthorized, do(h, req).Code)
			})
		}
	}
}

func TestProtectedEndpointsRejectWhenSecretMissing(t *testing.T) {
	t.Setenv("JWT_SECRET", "")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer abc")
	assert.Equal(t, http.StatusUnauthorized, do(endpoints.AccountMe, req).Code)
}

func TestOptionsPreflightSkipsAuth(t *testing.T) {
	for name, h := range protected {
		t.Run(name, func(t *testing.T) {
			rec := do(h, httptest.NewRequest(http.MethodOptions, "/", nil))
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
			assert.Empty(t, rec.Body.String())
		})
	}
}
