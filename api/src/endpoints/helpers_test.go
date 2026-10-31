package endpoints_test

import (
	"bytes"
	"encoding/json"
	"main/repo"
	"main/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

// repoHandler is the signature shared by all repository backed handlers.
type repoHandler func(http.ResponseWriter, *http.Request, repo.Repository)

// call routes a request through a mux router (so path variables work).
// caller is the account id the access token is issued for; empty means no Authorization header.
// body may be nil, a string/[]byte (sent as is) or any value (sent as JSON).
func call(t *testing.T, f *fakeRepo, h repoHandler, pattern, method, target, caller string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	switch b := body.(type) {
	case nil:
		reader = bytes.NewReader(nil)
	case string:
		reader = bytes.NewReader([]byte(b))
	case []byte:
		reader = bytes.NewReader(b)
	default:
		raw, err := json.Marshal(b)
		require.NoError(t, err)
		reader = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, target, reader)
	if caller != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken(t, caller))
	}
	router := mux.NewRouter()
	router.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) { h(w, r, f) })
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

// permitAll makes every store belong to the caller.
func (f *fakeRepo) permitAll() {
	f.OnIsPermitted = func(string, string) bool { return true }
}

// permitNone denies every store.
func (f *fakeRepo) permitNone() {
	f.OnIsPermitted = func(string, string) bool { return false }
}

func generateRestore(accountId string) (string, error) {
	return utils.GenerateRestoreToken(accountId)
}
