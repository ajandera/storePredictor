package utils_test

import (
	"encoding/json"
	"main/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode"

	"github.com/ajandera/sp_model/rdbsClientInfo"
	"github.com/stretchr/testify/assert"
)

func TestIsValidUUID(t *testing.T) {
	assert.True(t, utils.IsValidUUID("11111111-1111-4111-8111-111111111111"))
	assert.True(t, utils.IsValidUUID("AAAAAAAA-AAAA-4AAA-BAAA-AAAAAAAAAAAA"))
	assert.False(t, utils.IsValidUUID(""))
	assert.False(t, utils.IsValidUUID("not-a-uuid"))
	assert.False(t, utils.IsValidUUID("11111111-1111-1111-8111-111111111111")) // not v4
	assert.False(t, utils.IsValidUUID("11111111-1111-4111-8111-11111111111"))  // too short
	assert.False(t, utils.IsValidUUID("11111111-1111-4111-8111-111111111111'; DROP TABLE x;--"))
}

func TestRandomToken(t *testing.T) {
	a, b := utils.RandomToken(16), utils.RandomToken(16)
	assert.Len(t, a, 32) // hex encoded
	assert.NotEqual(t, a, b)
	assert.Empty(t, utils.RandomToken(0))
}

func TestStringWithCharset(t *testing.T) {
	s := utils.StringWithCharset(50, "abc")
	assert.Len(t, s, 50)
	assert.Empty(t, strings.Trim(s, "abc"))
}

func TestGeneratePassword(t *testing.T) {
	for i := 0; i < 50; i++ {
		pw := utils.GeneratePassword(16, 2, 3, 4)
		assert.Len(t, pw, 16)

		var special, digits, upper int
		for _, c := range pw {
			switch {
			case unicode.IsDigit(c):
				digits++
			case unicode.IsUpper(c):
				upper++
			case strings.ContainsRune("!@#$%&*", c):
				special++
			}
		}
		assert.GreaterOrEqual(t, special, 2)
		assert.GreaterOrEqual(t, digits, 3)
		assert.GreaterOrEqual(t, upper, 4)
	}
}

func TestWriteJSONError(t *testing.T) {
	rec := httptest.NewRecorder()
	utils.WriteJSONError(rec, http.StatusForbidden, "nope")

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	var body map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, false, body["success"])
	assert.Equal(t, "nope", body["error"])
}

func TestPublicAccountStripsSecrets(t *testing.T) {
	acc := rdbsClientInfo.Accounts{Password: "hash", RestoreToken: "tok"}
	pub := utils.PublicAccount(acc)
	assert.Empty(t, pub.Password)
	assert.Empty(t, pub.RestoreToken)
	assert.Equal(t, "hash", acc.Password, "input must not be mutated")

	list := utils.PublicAccounts([]rdbsClientInfo.Accounts{acc, acc})
	assert.Len(t, list, 2)
	for _, a := range list {
		assert.Empty(t, a.Password)
		assert.Empty(t, a.RestoreToken)
	}
	assert.NotNil(t, utils.PublicAccounts(nil))
}

func TestReadUserIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "1.1.1.1:1234"
	assert.Equal(t, "1.1.1.1:1234", utils.ReadUserIP(req))

	req.Header.Set("X-Forwarded-For", "2.2.2.2")
	assert.Equal(t, "2.2.2.2", utils.ReadUserIP(req))

	req.Header.Set("X-Real-Ip", "3.3.3.3")
	assert.Equal(t, "3.3.3.3", utils.ReadUserIP(req))
}

func TestBearerToken(t *testing.T) {
	cases := map[string]string{
		"Bearer abc":    "abc",
		"bearer abc":    "abc",
		"  Bearer abc ": "abc",
		"Bearer ":       "",
		"Basic abc":     "",
		"abc":           "",
		"":              "",
	}
	for header, want := range cases {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", header)
		assert.Equal(t, want, utils.BearerToken(req), "header %q", header)
	}
}

func TestSetupCORS(t *testing.T) {
	rec := httptest.NewRecorder()
	var w http.ResponseWriter = rec
	utils.SetupCORS(&w, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Headers"), "Authorization")
}

func TestValidateJWTConfig(t *testing.T) {
	t.Setenv("JWT_SECRET", "short")
	assert.Error(t, utils.ValidateJWTConfig())
	t.Setenv("JWT_SECRET", strings.Repeat("x", 32))
	assert.NoError(t, utils.ValidateJWTConfig())
}

func TestRestoreToken(t *testing.T) {
	t.Setenv("JWT_SECRET", strings.Repeat("x", 32))
	id := "11111111-1111-4111-8111-111111111111"
	a, err := utils.GenerateRestoreToken(id)
	assert.NoError(t, err)
	b, _ := utils.GenerateRestoreToken(id)
	assert.NotEqual(t, a, b, "jti must make tokens unique")

	got, err := utils.ParseToken(a, utils.TokenTypeRestore)
	assert.NoError(t, err)
	assert.Equal(t, id, got)

	_, err = utils.ParseToken(a, utils.TokenTypeAccess)
	assert.ErrorIs(t, err, utils.ErrInvalidToken)
}
