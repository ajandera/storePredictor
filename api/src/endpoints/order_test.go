package endpoints_test

import (
	"main/endpoints"
	"main/repo"
	ord "main/structs/order"
	"net/http"
	"testing"

	"github.com/ajandera/sp_model/rdbsClientInfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderAdd(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	f := newFakeRepo()
	caller := f.addAccount("a", "").Id.String()
	h := endpoints.AddOrder
	send := func(c string, body interface{}) int {
		return call(t, f, h, "/o", http.MethodPost, "/o", c, body).Code
	}
	assert.Equal(t, http.StatusUnauthorized, send("", ord.Order{}))
	assert.Equal(t, http.StatusBadRequest, send(caller, "{bad"))
	assert.Equal(t, 0, f.called("EditAccount"))

	rec := call(t, f, h, "/o", http.MethodPost, "/o", caller, ord.Order{Name: "N", Email: "e@x.y", PlanRefer: "plan", Password: "pw"})
	require.Equal(t, http.StatusCreated, rec.Code)
	args := f.lastCall("EditAccount")
	assert.Equal(t, caller, args[0], "the account of the token is edited, not one from the body")
	assert.Equal(t, "N", args[1])
	assert.Equal(t, "plan", args[10])
	assert.NotContains(t, rec.Body.String(), "secret-", "secrets must not leak")
	assert.Equal(t, http.StatusOK, call(t, f, h, "/o", http.MethodOptions, "/o", "", nil).Code)
}

func TestOrderAccountOrder(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	f := newFakeRepo()
	owner := f.addAccount("owner", "")
	stranger := f.addAccount("stranger", "")
	storeId := newV4()
	f.OnIsPermitted = func(string, string) bool { return true }
	f.OnCreateOrder = func(string, string, string, float64, bool) string { return "order-1" }
	p := "/o/{accountId}/{storeId}"
	url := "/o/" + owner.Id.String() + "/" + storeId
	body := ord.Prediction{PlanRefer: "plan", Amount: 9.5, Paid: true}

	rec := call(t, f, endpoints.AccountOrder, p, http.MethodPost, url, owner.Id.String(), body)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "order-1")
	assert.Equal(t, []interface{}{owner.Id.String(), storeId, "plan", 9.5, true}, f.lastCall("CreateOrder"))

	assert.Equal(t, http.StatusUnauthorized, call(t, f, endpoints.AccountOrder, p, http.MethodPost, url, "", body).Code)
	assert.Equal(t, http.StatusBadRequest, call(t, f, endpoints.AccountOrder, p, http.MethodPost, url, owner.Id.String(), "{bad").Code)
	// order for an account of another tenant
	assert.Equal(t, http.StatusForbidden, call(t, f, endpoints.AccountOrder, p, http.MethodPost, url, stranger.Id.String(), body).Code)
	assert.Equal(t, 1, f.called("CreateOrder"))
	f.permitNone()
	assert.Equal(t, http.StatusForbidden, call(t, f, endpoints.AccountOrder, p, http.MethodPost, url, owner.Id.String(), body).Code)
	assert.Equal(t, http.StatusOK, call(t, f, endpoints.AccountOrder, p, http.MethodOptions, url, "", nil).Code)
}

func TestOrderUpgrade(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	f := newFakeRepo()
	owner := f.addAccount("owner", "")
	child := f.addAccount("child", owner.Id.String())
	p := "/u/{accountId}"
	do := func(caller, target string, body interface{}) int {
		return call(t, f, endpoints.UpgradeOrder, p, http.MethodPost, "/u/"+target, caller, body).Code
	}

	assert.Equal(t, http.StatusUnauthorized, do("", owner.Id.String(), ord.Upgrade{}))
	assert.Equal(t, http.StatusForbidden, do(child.Id.String(), owner.Id.String(), ord.Upgrade{}), "child cannot upgrade its owner")
	assert.Equal(t, http.StatusBadRequest, do(owner.Id.String(), owner.Id.String(), "{bad"))
	assert.Equal(t, 0, f.called("EditAccount"))

	for _, year := range []bool{false, true} {
		rec := call(t, f, endpoints.UpgradeOrder, p, http.MethodPost, "/u/"+child.Id.String(), owner.Id.String(),
			ord.Upgrade{PlanRefer: "pro", Year: year})
		require.Equal(t, http.StatusCreated, rec.Code)
		args := f.lastCall("EditAccount")
		assert.Equal(t, child.Id.String(), args[0])
		assert.Equal(t, "pro", args[10])
		assert.NotEmpty(t, args[9], "paidTo date is set")
	}
	assert.Equal(t, http.StatusOK, call(t, f, endpoints.UpgradeOrder, p, http.MethodOptions, "/u/x", "", nil).Code)
}

func TestOrderManualPredictionGuards(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	f, caller, storeId, _ := supplierFixture(t)
	p := "/m/{storeId}/{date}"
	run := func(c, store string) int {
		return call(t, f, func(w http.ResponseWriter, r *http.Request, m repo.Repository) {
			endpoints.ManualPrediction(w, r, m, nil)
		}, p, http.MethodGet, "/m/"+store+"/2030-01-01", c, nil).Code
	}
	assert.Equal(t, http.StatusUnauthorized, run("", storeId))
	assert.Equal(t, http.StatusForbidden, run(caller, newV4()))
	// no messaging publisher configured
	assert.Equal(t, http.StatusInternalServerError, run(caller, storeId))
	assert.Equal(t, http.StatusOK, call(t, f, func(w http.ResponseWriter, r *http.Request, m repo.Repository) {
		endpoints.ManualPrediction(w, r, m, nil)
	}, p, http.MethodOptions, "/m/x/y", "", nil).Code)
}

func TestOrderOnetimePredictionWithoutBroker(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	t.Setenv("RABBIT_USER", "u")
	t.Setenv("RABBIT_PASS", "p")
	t.Setenv("RABBIT_HOST", "127.0.0.1:1") // nothing listens here
	f, caller, storeId, _ := supplierFixture(t)
	orderId := newV4()
	f.OnGetOrderById = func(id string) rdbsClientInfo.Orders {
		if id == orderId {
			return rdbsClientInfo.Orders{StoreRefer: storeId}
		}
		return rdbsClientInfo.Orders{}
	}
	p := "/one/{orderId}"
	assert.Equal(t, http.StatusUnauthorized, call(t, f, endpoints.OnetimePrediction, p, http.MethodGet, "/one/"+orderId, "", nil).Code)
	assert.Equal(t, http.StatusForbidden, call(t, f, endpoints.OnetimePrediction, p, http.MethodGet, "/one/"+newV4(), caller, nil).Code)
	assert.Equal(t, http.StatusInternalServerError, call(t, f, endpoints.OnetimePrediction, p, http.MethodGet, "/one/"+orderId, caller, nil).Code)
	assert.Equal(t, http.StatusOK, call(t, f, endpoints.OnetimePrediction, p, http.MethodOptions, "/one/x", "", nil).Code)
}
