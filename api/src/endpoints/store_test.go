package endpoints_test

import (
	"main/endpoints"
	sw "main/structs/storeWeights"
	st "main/structs/stores"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCreate(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	f := newFakeRepo()
	owner := f.addAccount("owner", "")
	child := f.addAccount("child", owner.Id.String())

	// child creates the store under the owner (tenant) account
	rec := call(t, f, endpoints.Stores, "/stores", http.MethodPost, "/stores", child.Id.String(),
		st.Stores{CountryCode: "CZ", Url: "shop.cz", Feed: "f", Window: 7})
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, 1, f.called("CreateStore"))
	args := f.lastCall("CreateStore")
	assert.Equal(t, "CZ", args[0])
	assert.Equal(t, "shop.cz", args[1])
	assert.True(t, strings.HasPrefix(args[2].(string), "SP-"))
	assert.Equal(t, owner.Id.String(), args[3])

	assert.Equal(t, http.StatusBadRequest,
		call(t, f, endpoints.Stores, "/stores", http.MethodPost, "/stores", owner.Id.String(), "{bad").Code)
	assert.Equal(t, http.StatusUnauthorized,
		call(t, f, endpoints.Stores, "/stores", http.MethodPost, "/stores", "", st.Stores{}).Code)
	assert.Equal(t, 1, f.called("CreateStore"))
}

func TestStoreUpdate(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	f := newFakeRepo()
	acc := f.addAccount("a", "")
	storeId := newV4()
	body := st.Stores{ID: storeId, CountryCode: "SK", Url: "x.sk"}

	f.permitNone()
	assert.Equal(t, http.StatusForbidden,
		call(t, f, endpoints.UpdateStores, "/stores", http.MethodPut, "/stores", acc.Id.String(), body).Code)
	assert.Equal(t, 0, f.called("EditStore"))

	f.permitAll()
	assert.Equal(t, http.StatusCreated,
		call(t, f, endpoints.UpdateStores, "/stores", http.MethodPut, "/stores", acc.Id.String(), body).Code)
	assert.Equal(t, storeId, f.lastCall("EditStore")[0])

	assert.Equal(t, http.StatusBadRequest,
		call(t, f, endpoints.UpdateStores, "/stores", http.MethodPut, "/stores", acc.Id.String(), "nope").Code)
	assert.Equal(t, http.StatusUnauthorized,
		call(t, f, endpoints.UpdateStores, "/stores", http.MethodPut, "/stores", "", body).Code)
}

func TestStoreDelete(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	f := newFakeRepo()
	acc := f.addAccount("a", "")
	storeId := newV4()
	url := "/stores/" + storeId

	f.permitNone()
	assert.Equal(t, http.StatusForbidden,
		call(t, f, endpoints.DeleteStores, "/stores/{storeId}", http.MethodDelete, url, acc.Id.String(), nil).Code)
	assert.Equal(t, 0, f.called("DeleteStore"))

	f.permitAll()
	assert.Equal(t, http.StatusOK,
		call(t, f, endpoints.DeleteStores, "/stores/{storeId}", http.MethodDelete, url, acc.Id.String(), nil).Code)
	assert.Equal(t, []interface{}{storeId}, f.lastCall("DeleteStore"))

	assert.Equal(t, http.StatusUnauthorized,
		call(t, f, endpoints.DeleteStores, "/stores/{storeId}", http.MethodDelete, url, "", nil).Code)
}

func TestStoreGetStoresUsesTenant(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	f := newFakeRepo()
	owner := f.addAccount("owner", "")
	child := f.addAccount("child", owner.Id.String())

	rec := call(t, f, endpoints.GetStores, "/stores", http.MethodGet, "/stores", child.Id.String(), nil)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, []interface{}{owner.Id.String()}, f.lastCall("GetStoresByAccount"))

	assert.Equal(t, http.StatusUnauthorized,
		call(t, f, endpoints.GetStores, "/stores", http.MethodGet, "/stores", "", nil).Code)
}

func TestStoreGetStoresByAccountEdgeCases(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	f := newFakeRepo()
	acc := f.addAccount("a", "")

	assert.Equal(t, http.StatusUnauthorized,
		call(t, f, endpoints.GetStoresByAccount, "/stores/{accountId}", http.MethodGet, "/stores/"+acc.Id.String(), "", nil).Code)
	// empty path variable
	assert.Equal(t, http.StatusBadRequest,
		call(t, f, endpoints.GetStoresByAccount, "/stores/{accountId:.*}", http.MethodGet, "/stores/", acc.Id.String(), nil).Code)
}

func TestStoreWeights(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	f := newFakeRepo()
	acc := f.addAccount("a", "")
	storeId := newV4()
	body := sw.StoreWeights{StoreRefer: storeId, Name: "n", Beta: 1}

	// read
	f.permitNone()
	assert.Equal(t, http.StatusForbidden,
		call(t, f, endpoints.GetStoreWeights, "/w/{storeId}", http.MethodGet, "/w/"+storeId, acc.Id.String(), nil).Code)
	f.permitAll()
	assert.Equal(t, http.StatusOK,
		call(t, f, endpoints.GetStoreWeights, "/w/{storeId}", http.MethodGet, "/w/"+storeId, acc.Id.String(), nil).Code)
	assert.Equal(t, http.StatusUnauthorized,
		call(t, f, endpoints.GetStoreWeights, "/w/{storeId}", http.MethodGet, "/w/"+storeId, "", nil).Code)

	// create / update share validation rules
	for _, c := range []struct {
		name    string
		h       repoHandler
		method  string
		status  int
		dbCalls string
	}{
		{"create", endpoints.CreateStoreWeights, http.MethodPost, http.StatusCreated, "CreateStoreWeights"},
		{"update", endpoints.UpdateStoreWeights, http.MethodPut, http.StatusCreated, "EditStoreWeights"},
	} {
		t.Run(c.name, func(t *testing.T) {
			f := newFakeRepo()
			acc := f.addAccount("a", "")

			assert.Equal(t, http.StatusUnauthorized, call(t, f, c.h, "/w", c.method, "/w", "", body).Code)
			assert.Equal(t, http.StatusBadRequest, call(t, f, c.h, "/w", c.method, "/w", acc.Id.String(), "{bad").Code)

			f.permitNone()
			assert.Equal(t, http.StatusForbidden, call(t, f, c.h, "/w", c.method, "/w", acc.Id.String(), body).Code)
			assert.Equal(t, 0, f.called(c.dbCalls))

			f.permitAll()
			assert.Equal(t, c.status, call(t, f, c.h, "/w", c.method, "/w", acc.Id.String(), body).Code)
			assert.Equal(t, storeId, f.lastCall(c.dbCalls)[0])
		})
	}
}
