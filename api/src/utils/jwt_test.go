package utils_test

import (
	"main/utils"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestJWTGenerationAndValidation(t *testing.T) {
	os.Setenv("JWT_SECRET", "super-secret-test-jwt-key-with-sufficient-length-32-chars")

	tokens, err := utils.GenerateJWT("Test User", "11111111-1111-4111-8111-111111111111")
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens["access_token"])
	assert.NotEmpty(t, tokens["refresh_token"])

	// Valid access token
	accId, err := utils.ParseToken(tokens["access_token"], utils.TokenTypeAccess)
	assert.NoError(t, err)
	assert.Equal(t, "11111111-1111-4111-8111-111111111111", accId)

	// Access token used as refresh token must fail
	_, err = utils.ParseToken(tokens["access_token"], utils.TokenTypeRefresh)
	assert.Error(t, err)

	// Valid refresh token
	accId, err = utils.ParseToken(tokens["refresh_token"], utils.TokenTypeRefresh)
	assert.NoError(t, err)
	assert.Equal(t, "11111111-1111-4111-8111-111111111111", accId)

	// Garbage token must fail
	_, err = utils.ParseToken("garbage", utils.TokenTypeAccess)
	assert.Error(t, err)

	// Restore token
	restoreToken, err := utils.GenerateRestoreToken("11111111-1111-4111-8111-111111111111")
	assert.NoError(t, err)
	accId, err = utils.ParseToken(restoreToken, utils.TokenTypeRestore)
	assert.NoError(t, err)
	assert.Equal(t, "11111111-1111-4111-8111-111111111111", accId)
}

func TestExpiredToken(t *testing.T) {
	key := "super-secret-test-jwt-key-with-sufficient-length-32-chars"
	os.Setenv("JWT_SECRET", key)

	expired := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"client":     "User",
		"id":         "11111111-1111-4111-8111-111111111111",
		"typ":        utils.TokenTypeAccess,
		"exp":        time.Now().Add(-1 * time.Hour).Unix(),
	})
	tokenStr, _ := expired.SignedString([]byte(key))

	_, err := utils.ParseToken(tokenStr, utils.TokenTypeAccess)
	assert.Error(t, err)
}
