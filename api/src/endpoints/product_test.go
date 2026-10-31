package endpoints_test

import (
	"fmt"
	"main/endpoints"
	pr "main/structs/product"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// storeGetCases are the read endpoints guarded only by store permission.
var storeGetCases = []struct {
	name, pattern, urlSuffix, method string
	h                                repoHandler
}{
	{"products", "/p/{storeId}/{limit}/{offset}", "/5/10", "GetProducts", endpoints.GetProducts},
	{"warehouse", "/p/{storeId}/{limit}/{offset}", "/5/10", "GetProductsWarehouse", endpoints.GetProductsWarehouse},
	{"product warehouse", "/p/{storeId}/{code}", "/P1", "GetProduct", endpoints.GetProductWarehouse},
	{"products to store", "/p/{storeId}/{limit}/{offset}", "/5/10", "GetProductsToStore", endpoints.GetProductsToStore},
	{"product to store", "/p/{storeId}/{productCode}", "/P1", "GetProductToStore", endpoints.GetProductToStore},
	{"client orders", "/p/{storeId}/{limit}/{offset}", "/5/10", "GetOrders", endpoints.GetClientOrders},
	{"orders for product", "/p/{storeId}/{productCode}", "/P1", "GetOrdersWithProduct", endpoints.GetOrdersForProduct},
}

func TestProductAndOrderReadEndpointsEnforcePermission(t *testing.T) {
	for _, c := range storeGetCases {
		t.Run(c.name, func(t *testing.T) {
			f, caller, storeId, _ := supplierFixture(t)
			url := "/p/" + storeId + c.urlSuffix
			do := func(caller, url string) int {
				return call(t, f, c.h, c.pattern, http.MethodGet, url, caller, nil).Code
			}
			assert.Equal(t, http.StatusUnauthorized, do("", url))
			assert.Equal(t, http.StatusForbidden, do(caller, "/p/"+newV4()+c.urlSuffix))
			assert.Equal(t, 0, f.called(c.method))

			assert.Equal(t, http.StatusOK, do(caller, url))
			require.Equal(t, 1, f.called(c.method), "repository must be queried once")
			assert.Contains(t, fmt.Sprint(f.lastCall(c.method)), storeId)
			assert.Equal(t, http.StatusOK, call(t, f, c.h, c.pattern, http.MethodOptions, url, "", nil).Code)
		})
	}
}

func TestProductPagingArguments(t *testing.T) {
	f, caller, storeId, _ := supplierFixture(t)
	call(t, f, endpoints.GetProductsWarehouse, "/p/{storeId}/{limit}/{offset}", http.MethodGet,
		"/p/"+storeId+"/5/10", caller, nil)
	assert.Equal(t, []interface{}{storeId, 5, 10}, f.lastCall("GetProductsWarehouse"))

	call(t, f, endpoints.GetClientOrders, "/p/{storeId}/{limit}/{offset}", http.MethodGet,
		"/p/"+storeId+"/3/4", caller, nil)
	args := f.lastCall("GetOrders")
	assert.Equal(t, map[string]interface{}{"store_id": storeId}, args[0])
	assert.Equal(t, 3, args[1])
	assert.Equal(t, 4, args[2])
}

func TestProductCreateAndUpdate(t *testing.T) {
	for _, c := range []struct {
		name   string
		h      repoHandler
		method string
		db     string
	}{
		{"create", endpoints.CreateProduct, http.MethodPost, "CreateProduct"},
		{"update", endpoints.UpdateProduct, http.MethodPut, "UpdateProduct"},
	} {
		t.Run(c.name, func(t *testing.T) {
			f, caller, storeId, _ := supplierFixture(t)
			send := func(caller string, body interface{}) int {
				return call(t, f, c.h, "/p", c.method, "/p", caller, body).Code
			}
			body := pr.Product{ProductCode: "P1", Name: "n", Quantity: 3, StoreId: storeId}

			assert.Equal(t, http.StatusUnauthorized, send("", body))
			assert.Equal(t, http.StatusBadRequest, send(caller, "{bad"))
			assert.Equal(t, http.StatusForbidden, send(caller, pr.Product{StoreId: newV4()}))
			assert.Equal(t, 0, f.called(c.db))

			assert.Equal(t, http.StatusCreated, send(caller, body))
			require.Equal(t, 1, f.called(c.db))
			assert.Contains(t, f.lastCall(c.db), "P1")
			assert.Contains(t, f.lastCall(c.db), storeId)
		})
	}
}
