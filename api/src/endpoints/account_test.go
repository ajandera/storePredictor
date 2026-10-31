package endpoints_test

import (
	"encoding/json"
	"main/endpoints"
	"main/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateJWT(accountId string) (map[string]string, error) {
	return utils.GenerateJWT("Test", accountId)
}

func decode(t *testing.T, rec *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body), rec.Body.String())
	return body
}

func TestAccountMeReturnsCallerWithoutSecrets(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	f := newFakeRepo()
	me := f.addAccount("me", "")

	rec := httptest.NewRecorder()
	endpoints.AccountMe(rec, bearer(t, http.MethodGet, "/account/me", me.Id.String()), f)

	require.Equal(t, http.StatusOK, rec.Code)
	body := decode(t, rec)
	assert.Equal(t, true, body["success"])
	user := body["user"].(map[string]interface{})
	assert.Equal(t, "me", user["Name"])
	assert.Empty(t, user["Password"])
	assert.Empty(t, user["RestoreToken"])
	assert.NotContains(t, rec.Body.String(), "secret-")
}

func TestGetAccountPermissions(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	f := newFakeRepo()
	owner := f.addAccount("owner", "")
	child := f.addAccount("child", owner.Id.String())
	sibling := f.addAccount("sibling", owner.Id.String())
	stranger := f.addAccount("stranger", "")

	get := func(caller, target string) *httptest.ResponseRecorder {
		return serve(t, "/accounts/{accountId}", func(w http.ResponseWriter, r *http.Request) {
			endpoints.GetAccount(w, r, f)
		}, bearer(t, http.MethodGet, "/accounts/"+target, caller))
	}

	cases := []struct {
		name           string
		caller, target string
		want           int
	}{
		{"self", owner.Id.String(), owner.Id.String(), http.StatusOK},
		{"owner views child", owner.Id.String(), child.Id.String(), http.StatusOK},
		{"child views owner", child.Id.String(), owner.Id.String(), http.StatusOK},
		{"child views sibling", child.Id.String(), sibling.Id.String(), http.StatusOK},
		{"stranger views owner", stranger.Id.String(), owner.Id.String(), http.StatusForbidden},
		{"owner views stranger", owner.Id.String(), stranger.Id.String(), http.StatusForbidden},
		{"unknown account", owner.Id.String(), newV4(), http.StatusForbidden},
		{"invalid id", owner.Id.String(), "not-a-uuid", http.StatusForbidden},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rec := get(c.caller, c.target)
			assert.Equal(t, c.want, rec.Code)
			assert.NotContains(t, rec.Body.String(), "secret-")
			if c.want == http.StatusForbidden {
				assert.Equal(t, "Not Permitted", decode(t, rec)["error"])
			}
		})
	}
}

func TestGetAccountsReturnsOnlyOwnTenant(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	f := newFakeRepo()
	owner := f.addAccount("owner", "")
	child := f.addAccount("child", owner.Id.String())
	f.addAccount("stranger", "")

	rec := httptest.NewRecorder()
	// a child account sees its tenant: the owner and all owner's children
	endpoints.GetAccounts(rec, bearer(t, http.MethodGet, "/accounts", child.Id.String()), f)

	require.Equal(t, http.StatusOK, rec.Code)
	accounts := decode(t, rec)["accounts"].([]interface{})
	var names []string
	for _, a := range accounts {
		names = append(names, a.(map[string]interface{})["Name"].(string))
	}
	assert.ElementsMatch(t, []string{"owner", "child"}, names)
	assert.NotContains(t, rec.Body.String(), "secret-")
}

func TestDeleteAccountsPermissions(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)

	del := func(f *fakeRepo, caller, target string) *httptest.ResponseRecorder {
		return serve(t, "/accounts/{accountId}", func(w http.ResponseWriter, r *http.Request) {
			endpoints.DeleteAccounts(w, r, f)
		}, bearer(t, http.MethodDelete, "/accounts/"+target, caller))
	}

	t.Run("owner deletes own child", func(t *testing.T) {
		f := newFakeRepo()
		owner := f.addAccount("owner", "")
		child := f.addAccount("child", owner.Id.String())

		rec := del(f, owner.Id.String(), child.Id.String())
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, []string{child.Id.String()}, f.deleted)
	})

	t.Run("account deletes itself", func(t *testing.T) {
		f := newFakeRepo()
		owner := f.addAccount("owner", "")

		assert.Equal(t, http.StatusOK, del(f, owner.Id.String(), owner.Id.String()).Code)
		assert.Equal(t, []string{owner.Id.String()}, f.deleted)
	})

	t.Run("child cannot delete its owner", func(t *testing.T) {
		f := newFakeRepo()
		owner := f.addAccount("owner", "")
		child := f.addAccount("child", owner.Id.String())

		assert.Equal(t, http.StatusForbidden, del(f, child.Id.String(), owner.Id.String()).Code)
		assert.Empty(t, f.deleted)
	})

	t.Run("child cannot delete sibling", func(t *testing.T) {
		f := newFakeRepo()
		owner := f.addAccount("owner", "")
		child := f.addAccount("child", owner.Id.String())
		sibling := f.addAccount("sibling", owner.Id.String())

		assert.Equal(t, http.StatusForbidden, del(f, child.Id.String(), sibling.Id.String()).Code)
		assert.Empty(t, f.deleted)
	})

	t.Run("stranger cannot delete foreign account", func(t *testing.T) {
		f := newFakeRepo()
		owner := f.addAccount("owner", "")
		stranger := f.addAccount("stranger", "")

		assert.Equal(t, http.StatusForbidden, del(f, stranger.Id.String(), owner.Id.String()).Code)
		assert.Empty(t, f.deleted)
	})

	t.Run("unauthenticated request deletes nothing", func(t *testing.T) {
		f := newFakeRepo()
		owner := f.addAccount("owner", "")

		req := httptest.NewRequest(http.MethodDelete, "/accounts/"+owner.Id.String(), nil)
		rec := serve(t, "/accounts/{accountId}", func(w http.ResponseWriter, r *http.Request) {
			endpoints.DeleteAccounts(w, r, f)
		}, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Empty(t, f.deleted)
	})
}

func TestGetChildAccountPermissions(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	f := newFakeRepo()
	owner := f.addAccount("owner", "")
	f.addAccount("child", owner.Id.String())
	stranger := f.addAccount("stranger", "")

	get := func(caller string) *httptest.ResponseRecorder {
		return serve(t, "/account/child/{id}", func(w http.ResponseWriter, r *http.Request) {
			endpoints.GetChildAccount(w, r, f)
		}, bearer(t, http.MethodGet, "/account/child/"+owner.Id.String(), caller))
	}

	rec := get(owner.Id.String())
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Len(t, decode(t, rec)["accounts"], 1)
	assert.NotContains(t, rec.Body.String(), "secret-")

	assert.Equal(t, http.StatusForbidden, get(stranger.Id.String()).Code)
}

func TestGetStoresByAccountPermissions(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	f := newFakeRepo()
	owner := f.addAccount("owner", "")
	stranger := f.addAccount("stranger", "")

	get := func(caller string) *httptest.ResponseRecorder {
		return serve(t, "/stores/{accountId}", func(w http.ResponseWriter, r *http.Request) {
			endpoints.GetStoresByAccount(w, r, f)
		}, bearer(t, http.MethodGet, "/stores/"+owner.Id.String(), caller))
	}

	rec := get(owner.Id.String())
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, true, decode(t, rec)["success"])

	assert.Equal(t, http.StatusForbidden, get(stranger.Id.String()).Code)
}
