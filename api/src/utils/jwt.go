package utils

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Token types stored in the "typ" claim so an access token cannot be used
// as a refresh token and vice versa.
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
	TokenTypeRestore = "restore"

	accessTokenTTL  = 30 * time.Minute
	refreshTokenTTL = 24 * time.Hour
	restoreTokenTTL = 24 * time.Hour

	// minSecretLength is the minimal accepted length of JWT_SECRET (256 bits for HS256).
	minSecretLength = 32
)

// ErrInvalidToken is returned for every token that must not be trusted.
var ErrInvalidToken = errors.New("invalid token")

// signingKey returns the HMAC key used to sign and verify JWT tokens.
// The key is read from the JWT_SECRET environment variable.
func signingKey() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < minSecretLength {
		return nil, fmt.Errorf("JWT_SECRET must be set and at least %d characters long", minSecretLength)
	}
	return []byte(secret), nil
}

// ValidateJWTConfig checks that the JWT configuration is usable. Call it on startup.
func ValidateJWTConfig() error {
	_, err := signingKey()
	return err
}

// GenerateJWT creates an access/refresh token pair for the account.
func GenerateJWT(name string, id string) (map[string]string, error) {
	key, err := signingKey()
	if err != nil {
		return nil, err
	}
	now := time.Now()

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"client":     name,
		"id":         id,
		"typ":        TokenTypeAccess,
		"iat":        now.Unix(),
		"exp":        now.Add(accessTokenTTL).Unix(),
	})
	accessString, err := access.SignedString(key)
	if err != nil {
		return nil, err
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"typ": TokenTypeRefresh,
		"iat": now.Unix(),
		"exp": now.Add(refreshTokenTTL).Unix(),
	})
	refreshString, err := refresh.SignedString(key)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  accessString,
		"refresh_token": refreshString,
	}, nil
}

// GenerateRestoreToken creates signed, URL-safe password restore token valid for 24 hours.
// The random "jti" claim makes every token unique.
func GenerateRestoreToken(id string) (string, error) {
	key, err := signingKey()
	if err != nil {
		return "", err
	}
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"typ": TokenTypeRestore,
		"jti": RandomToken(16),
		"iat": now.Unix(),
		"exp": now.Add(restoreTokenTTL).Unix(),
	})
	return token.SignedString(key)
}

// ParseToken verifies the signature, expiry and type of the token and returns
// the account id stored in it. Any failure results in ErrInvalidToken.
func ParseToken(tokenString string, expectedType string) (string, error) {
	key, err := signingKey()
	if err != nil {
		return "", err
	}
	if tokenString == "" {
		return "", ErrInvalidToken
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	token, err := parser.Parse(tokenString, func(*jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil || token == nil || !token.Valid {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidToken
	}
	// exp is mandatory, MapClaims.Valid only checks it when present
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return "", ErrInvalidToken
	}
	if typ, _ := claims["typ"].(string); typ != expectedType {
		return "", ErrInvalidToken
	}
	id, _ := claims["id"].(string)
	if !IsValidUUID(id) {
		return "", ErrInvalidToken
	}
	return id, nil
}

// BearerToken extracts the token from the "Authorization: Bearer <token>" header.
func BearerToken(r *http.Request) string {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if len(header) > 7 && strings.EqualFold(header[:7], "Bearer ") {
		return strings.TrimSpace(header[7:])
	}
	return ""
}
