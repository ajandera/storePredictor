package endpoints_test

import (
	"main/endpoints"
	a "main/structs/account"
	pw "main/structs/pw"
	"net/http"
	"testing"

	"github.com/ajandera/sp_model/rdbsClientInfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountSendNewPw(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	t.Setenv("SMTP_URL", "") // no mail is sent
	f := newFakeRepo()
	acc := f.addAccount("a", "")
	p := "/forgot/{lang}/{accountId}"

	rec := call(t, f, endpoints.SendNewPw, p, http.MethodGet, "/forgot/en/"+acc.Id.String(), "", nil)
	require.Equal(t, http.StatusOK, rec.Code)
	args := f.lastCall("SetRestorePw")
	assert.Equal(t, acc.Id.String(), args[0])
	assert.NotEmpty(t, args[1], "a restore token is stored")

	// unknown account answers success without storing anything (does not reveal existence)
	f2 := newFakeRepo()
	rec = call(t, f2, endpoints.SendNewPw, p, http.MethodGet, "/forgot/en/"+newV4(), "", nil)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 0, f2.called("SetRestorePw"))

	t.Setenv("JWT_SECRET", "") // token cannot be generated
	assert.Equal(t, http.StatusInternalServerError,
		call(t, f, endpoints.SendNewPw, p, http.MethodGet, "/forgot/en/"+acc.Id.String(), "", nil).Code)
	assert.Equal(t, http.StatusOK, call(t, f, endpoints.SendNewPw, p, http.MethodOptions, "/forgot/en/x", "", nil).Code)
}

func TestAccountRegistrationValidation(t *testing.T) {
	f := newFakeRepo()
	existing := f.addAccount("old", "")
	existing.Email = "taken@example.com"
	f.OnGetAccountByEmail = func(email string) rdbsClientInfo.Accounts {
		if email == existing.Email {
			return existing
		}
		return rdbsClientInfo.Accounts{}
	}
	reg := func(body interface{}) *http.Response {
		rec := call(t, f, endpoints.Accounts, "/reg/{lang}", http.MethodPost, "/reg/en", "", body)
		return rec.Result()
	}

	assert.Equal(t, http.StatusBadRequest, reg("{bad").StatusCode)
	assert.Equal(t, http.StatusBadRequest, reg(a.Account{Email: "", Password: "longenough"}).StatusCode)
	assert.Equal(t, http.StatusBadRequest, reg(a.Account{Email: "x@y.z", Password: "short"}).StatusCode, "7+ chars needed")

	rec := call(t, f, endpoints.Accounts, "/reg/{lang}", http.MethodPost, "/reg/en", "",
		a.Account{Email: existing.Email, Password: "longenough"})
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Email already registered")
	assert.Equal(t, 0, f.called("CreateAccount"))
	assert.Equal(t, http.StatusOK, call(t, f, endpoints.Accounts, "/reg/{lang}", http.MethodOptions, "/reg/en", "", nil).Code)
}

func TestAccountUpdate(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	setup := func() (*fakeRepo, rdbsClientInfo.Accounts, rdbsClientInfo.Accounts, rdbsClientInfo.Accounts, rdbsClientInfo.Accounts) {
		f := newFakeRepo()
		owner := f.addAccount("owner", "")
		child := f.addAccount("child", owner.Id.String())
		stranger := f.addAccount("stranger", "")
		taken := f.addAccount("taken", "")
		taken.Email = "taken@example.com"
		f.accounts[taken.Id.String()] = taken
		f.OnGetAccountByEmail = func(email string) rdbsClientInfo.Accounts {
			if email == taken.Email {
				return taken
			}
			return rdbsClientInfo.Accounts{}
		}
		return f, owner, child, stranger, taken
	}
	put := func(f *fakeRepo, caller string, body interface{}) int {
		return call(t, f, endpoints.UpdateAccounts, "/accounts", http.MethodPut, "/accounts", caller, body).Code
	}

	t.Run("auth and input", func(t *testing.T) {
		f, owner, _, _, _ := setup()
		assert.Equal(t, http.StatusUnauthorized, put(f, "", a.Account{}))
		assert.Equal(t, http.StatusBadRequest, put(f, owner.Id.String(), "{bad"))
		assert.Equal(t, 0, f.called("EditAccount"))
	})

	t.Run("owner creates a child account", func(t *testing.T) {
		f, owner, _, _, _ := setup()
		assert.Equal(t, http.StatusCreated, put(f, owner.Id.String(),
			a.Account{Email: "new@example.com", Password: "longenough", Name: "kid", Role: "viewer"}))
		args := f.lastCall("EditAccount")
		assert.Equal(t, owner.Id.String(), args[12], "parent is forced to the caller")
	})

	t.Run("a child cannot create accounts", func(t *testing.T) {
		f, _, child, _, _ := setup()
		assert.Equal(t, http.StatusForbidden, put(f, child.Id.String(), a.Account{Email: "n@e.com", Password: "longenough"}))
		assert.Equal(t, 0, f.called("EditAccount"))
	})

	t.Run("new account validation", func(t *testing.T) {
		f, owner, _, _, taken := setup()
		assert.Equal(t, http.StatusBadRequest, put(f, owner.Id.String(), a.Account{Email: "n@e.com", Password: "short"}))
		assert.Equal(t, http.StatusConflict, put(f, owner.Id.String(), a.Account{Email: taken.Email, Password: "longenough"}))
		assert.Equal(t, 0, f.called("EditAccount"))
	})

	t.Run("self update cannot escalate role, parent, plan or payment", func(t *testing.T) {
		f, owner, _, _, _ := setup()
		require.Equal(t, http.StatusCreated, put(f, owner.Id.String(), a.Account{
			ID: owner.Id.String(), Name: "me", Role: "admin", Parent: newV4(), PaidTo: "2099-01-01", PlanRefer: "gold",
		}))
		args := f.lastCall("EditAccount")
		assert.Equal(t, "me", args[1])
		assert.Equal(t, "", args[9], "PaidTo")
		assert.Equal(t, "", args[10], "PlanRefer")
		assert.Equal(t, "", args[11], "Role")
		assert.Equal(t, "", args[12], "Parent")
	})

	t.Run("owner edits own child, but not strangers", func(t *testing.T) {
		f, owner, child, stranger, _ := setup()
		assert.Equal(t, http.StatusCreated, put(f, owner.Id.String(), a.Account{ID: child.Id.String(), Name: "k"}))
		assert.Equal(t, http.StatusForbidden, put(f, owner.Id.String(), a.Account{ID: stranger.Id.String(), Name: "k"}))
		assert.Equal(t, http.StatusForbidden, put(f, child.Id.String(), a.Account{ID: owner.Id.String(), Name: "k"}))
		assert.Equal(t, 1, f.called("EditAccount"))
	})
}

func TestAccountUpdatePassword(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	f := newFakeRepo()
	acc := f.addAccount("a", "")
	token, err := generateRestore(acc.Id.String())
	require.NoError(t, err)
	acc.RestoreToken = token
	f.accounts[acc.Id.String()] = acc

	p := "/restore/{token}"
	send := func(token string, body interface{}) int {
		return call(t, f, endpoints.UpdatePassword, p, http.MethodPut, "/restore/"+token, "", body).Code
	}

	assert.Equal(t, http.StatusBadRequest, send(token, "{bad"))
	assert.Equal(t, http.StatusBadRequest, send("garbage", pw.Pw{Password: "longenough"}))
	assert.Equal(t, http.StatusBadRequest, send(token, pw.Pw{Password: "short"}))
	assert.Equal(t, 0, f.called("UpdatePw"))

	assert.Equal(t, http.StatusCreated, send(token, pw.Pw{Password: "longenough"}))
	assert.Equal(t, []interface{}{token, "longenough"}, f.lastCall("UpdatePw"))

	// a valid token that is not the one stored for the account (already used/replaced)
	other, err := generateRestore(acc.Id.String())
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, send(other, pw.Pw{Password: "longenough"}))
	assert.Equal(t, 1, f.called("UpdatePw"))
	assert.Equal(t, http.StatusOK, call(t, f, endpoints.UpdatePassword, p, http.MethodOptions, "/restore/x", "", nil).Code)
}

func TestAccountGetAccountOrders(t *testing.T) {
	f, caller, storeId, _ := supplierFixture(t)
	me := caller
	stranger := f.addAccount("stranger", "").Id.String()
	p := "/ao/{storeId}/{accountId}"
	do := func(c, store, acc string) int {
		return call(t, f, endpoints.GetAccountOrders, p, http.MethodGet, "/ao/"+store+"/"+acc, c, nil).Code
	}
	assert.Equal(t, http.StatusUnauthorized, do("", storeId, me))
	assert.Equal(t, http.StatusForbidden, do(me, newV4(), me))
	assert.Equal(t, http.StatusForbidden, do(me, storeId, stranger), "other tenant's account")
	assert.Equal(t, http.StatusOK, do(me, storeId, me))
	assert.Equal(t, []interface{}{me}, f.lastCall("GetAccountOrders"))
	assert.Equal(t, http.StatusOK, call(t, f, endpoints.GetAccountOrders, p, http.MethodOptions, "/ao/x/y", "", nil).Code)
}

func TestAccountGetByEmail(t *testing.T) {
	f := newFakeRepo()
	acc := f.addAccount("a", "")
	f.OnGetAccountByEmail = func(email string) rdbsClientInfo.Accounts {
		if email == "a@b.c" {
			return acc
		}
		return rdbsClientInfo.Accounts{}
	}
	rec := call(t, f, endpoints.GetAccountByEmail, "/e/{email}", http.MethodGet, "/e/a@b.c", "", nil)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), acc.Id.String())
	assert.NotContains(t, rec.Body.String(), "secret-")
	assert.Equal(t, http.StatusOK, call(t, f, endpoints.GetAccountByEmail, "/e/{email}", http.MethodOptions, "/e/x", "", nil).Code)
}

func TestAccountPlans(t *testing.T) {
	f := newFakeRepo()
	for name, c := range map[string]struct {
		h      repoHandler
		method string
	}{
		"plans":      {endpoints.GetPlans, "GetPlans"},
		"paid plans": {endpoints.GetPaidPlans, "GetPaidPlans"},
	} {
		t.Run(name, func(t *testing.T) {
			rec := call(t, f, c.h, "/plans", http.MethodGet, "/plans", "", nil)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, 1, f.called(c.method))
			assert.Equal(t, http.StatusOK, call(t, f, c.h, "/plans", http.MethodOptions, "/plans", "", nil).Code)
			assert.Equal(t, 1, f.called(c.method), "preflight does not query")
		})
	}
}

func TestAccountOptionsPreflightForAuthenticatedHandlers(t *testing.T) {
	f := newFakeRepo()
	for name, h := range map[string]repoHandler{
		"UpdateAccounts":   endpoints.UpdateAccounts,
		"DeleteAccounts":   endpoints.DeleteAccounts,
		"GetAccount":       endpoints.GetAccount,
		"GetChildAccount":  endpoints.GetChildAccount,
		"GetStoresByAccnt": endpoints.GetStoresByAccount,
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, http.StatusOK, call(t, f, h, "/x", http.MethodOptions, "/x", "", nil).Code)
		})
	}
}
