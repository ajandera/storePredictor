package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ajandera/sp_model/rdbsClientInfo"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

// fakeState is the in-memory data behind the default hooks of fakeRepo.
type fakeState struct {
	accounts map[string]rdbsClientInfo.Accounts
	stores   map[string][]rdbsClientInfo.Stores // by account id
	deleted  []string
}

// newFakeRepo returns a fake with account/store lookups backed by in-memory maps.
// Every other method records the call and returns zero values unless a test sets its On<Method> hook.
func newFakeRepo() *fakeRepo {
	f := &fakeRepo{fakeState: fakeState{
		accounts: map[string]rdbsClientInfo.Accounts{},
		stores:   map[string][]rdbsClientInfo.Stores{},
	}}
	f.OnGetAccountById = func(id string) rdbsClientInfo.Accounts { return f.accounts[id] }
	f.OnGetChildAccountById = func(id string) []rdbsClientInfo.Accounts {
		var out []rdbsClientInfo.Accounts
		for _, a := range f.accounts {
			if a.Parent == id {
				out = append(out, a)
			}
		}
		return out
	}
	f.OnGetStoresByAccount = func(id string) []rdbsClientInfo.Stores { return f.stores[id] }
	f.OnDeleteAccount = func(id string) {
		f.deleted = append(f.deleted, id)
		delete(f.accounts, id)
	}
	return f
}

// addAccount registers an account with a fresh id; parent may be empty.
func (f *fakeRepo) addAccount(name, parent string) rdbsClientInfo.Accounts {
	acc := rdbsClientInfo.Accounts{
		Id:           uuid.MustParse(newV4()),
		Name:         name,
		Parent:       parent,
		Password:     "secret-hash",
		RestoreToken: "secret-restore",
	}
	f.accounts[acc.Id.String()] = acc
	return acc
}

func newV4() string {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return id.String()
}

// serve routes the request through a mux router so path variables are populated.
func serve(t *testing.T, pattern string, h func(http.ResponseWriter, *http.Request), req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	router := mux.NewRouter()
	router.HandleFunc(pattern, h)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

// bearer returns a request carrying a valid access token for the account.
func bearer(t *testing.T, method, target string, accountId string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(method, target, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken(t, accountId))
	return req
}

func accessToken(t *testing.T, accountId string) string {
	t.Helper()
	tokens, err := generateJWT(accountId)
	require.NoError(t, err)
	return tokens["access_token"]
}
