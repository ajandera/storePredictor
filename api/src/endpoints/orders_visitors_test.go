package endpoints_test

import (
	"main/endpoints"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ajandera/sp_model/rdbsClientData"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// codeResolver makes CheckStoreCode succeed only for the given host.
func codeResolver(f *fakeRepo, host, storeId string) {
	f.OnCheckStoreCode = func(code, h string) string {
		if code == "SP-OK" && h == host {
			return storeId
		}
		return ""
	}
}

func post(t *testing.T, h repoHandler, f *fakeRepo, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	return call(t, f, h, "/x", http.MethodPost, "/x", "", body)
}

func TestOrderCreate(t *testing.T) {
	storeId := newV4()
	items := []rdbsClientData.Item{{ProductCode: "P", Quantity: 2}}

	t.Run("exact hostname", func(t *testing.T) {
		f := newFakeRepo()
		codeResolver(f, "shop.example.com", storeId)
		rec := post(t, endpoints.Orders, f, endpoints.Order{Code: "SP-OK", Url: "https://shop.example.com/cart",
			OrderId: "o1", Items: items, TotalPrice: 10, Currency: "EUR", Tag: "t"})
		require.Equal(t, http.StatusCreated, rec.Code)
		args := f.lastCall("SaveOrder")
		assert.Equal(t, 10.0, args[0])
		assert.Equal(t, storeId, args[2])
		assert.Equal(t, "o1", args[4])
		assert.Contains(t, rec.Body.String(), `"success":true`)
	})

	t.Run("falls back to the parent domain", func(t *testing.T) {
		f := newFakeRepo()
		codeResolver(f, "example.com", storeId)
		rec := post(t, endpoints.Orders, f, endpoints.Order{Code: "SP-OK", Url: "https://a.b.example.com/", OrderId: "o2"})
		require.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, storeId, f.lastCall("SaveOrder")[2])
	})

	t.Run("unknown store code is rejected without saving", func(t *testing.T) {
		f := newFakeRepo()
		rec := post(t, endpoints.Orders, f, endpoints.Order{Code: "bad", Url: "https://shop.example.com/"})
		assert.Equal(t, http.StatusCreated, rec.Code) // existing behaviour: body carries the error
		assert.Contains(t, rec.Body.String(), `"success":false`)
		assert.Contains(t, rec.Body.String(), "Not valid store code")
		assert.Equal(t, 0, f.called("SaveOrder"))
	})

	t.Run("single label host without a store is rejected", func(t *testing.T) {
		f := newFakeRepo()
		rec := post(t, endpoints.Orders, f, endpoints.Order{Code: "bad", Url: "http://localhost/"})
		assert.Contains(t, rec.Body.String(), `"success":false`)
		assert.Equal(t, 0, f.called("SaveOrder"))
	})

	t.Run("bad json", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, post(t, endpoints.Orders, newFakeRepo(), "{bad").Code)
	})

	t.Run("options preflight", func(t *testing.T) {
		rec := call(t, newFakeRepo(), endpoints.Orders, "/x", http.MethodOptions, "/x", "", nil)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Empty(t, rec.Body.String())
	})
}

func TestOrderShopify(t *testing.T) {
	storeId := newV4()
	f := newFakeRepo()
	f.OnGetStoreByUrl = func(host string) (string, error) {
		if host == "shop.example.com" {
			return storeId, nil
		}
		return "", nil
	}
	rec := post(t, endpoints.ShopifyOrders, f, endpoints.ShopifyOrder{
		Currency: "USD", CurrentTotalPrice: 30, OrderNumber: "n1",
		OrderStatusUrl: "https://shop.example.com:443/status",
		LineItems: []endpoints.ShopifyOrderItem{
			{Price: 10, Sku: "S1", Quantity: 1, Name: "one"},
			{Price: 20, Sku: "S2", Quantity: 2, Name: "two"},
		},
	})
	require.Equal(t, http.StatusCreated, rec.Code)
	args := f.lastCall("SaveOrder")
	assert.Equal(t, 30.0, args[0])
	assert.Equal(t, "USD", args[1])
	assert.Equal(t, storeId, args[2])
	got := args[3].([]rdbsClientData.Item)
	require.Len(t, got, 2)
	assert.Equal(t, "S2", got[1].ProductCode)
	assert.Equal(t, int8(2), got[1].Quantity)
	assert.Equal(t, "n1", args[4])

	assert.Equal(t, http.StatusBadRequest, post(t, endpoints.ShopifyOrders, newFakeRepo(), "{bad").Code)
	assert.Equal(t, http.StatusOK,
		call(t, f, endpoints.ShopifyOrders, "/x", http.MethodOptions, "/x", "", nil).Code)
}

func TestVisitorTracking(t *testing.T) {
	storeId := newV4()

	t.Run("exact hostname", func(t *testing.T) {
		f := newFakeRepo()
		codeResolver(f, "shop.example.com", storeId)
		rec := post(t, endpoints.Visitors, f, endpoints.Visitor{Code: "SP-OK", Url: "https://shop.example.com/p",
			Header: "ua", ProductCode: "P1", Tag: "t"})
		require.Equal(t, http.StatusCreated, rec.Code)
		args := f.lastCall("SaveVisitor")
		assert.Equal(t, storeId, args[1])
		assert.Equal(t, "https://shop.example.com/p", args[2])
		assert.Equal(t, "P1", args[4])
	})

	t.Run("parent domain", func(t *testing.T) {
		f := newFakeRepo()
		codeResolver(f, "example.com", storeId)
		rec := post(t, endpoints.Visitors, f, endpoints.Visitor{Code: "SP-OK", Url: "https://x.example.com/"})
		require.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, 1, f.called("SaveVisitor"))
	})

	t.Run("numeric product code uses the alias format", func(t *testing.T) {
		f := newFakeRepo()
		codeResolver(f, "shop.example.com", storeId)
		rec := post(t, endpoints.Visitors, f, `{"Code":"SP-OK","Url":"https://shop.example.com/","ProductCode":12}`)
		require.Equal(t, http.StatusCreated, rec.Code)
		assert.True(t, strings.HasPrefix(f.lastCall("SaveVisitor")[4].(string), "12"))
	})

	t.Run("unknown code is forbidden", func(t *testing.T) {
		f := newFakeRepo()
		rec := post(t, endpoints.Visitors, f, endpoints.Visitor{Code: "bad", Url: "https://shop.example.com/"})
		assert.Equal(t, http.StatusForbidden, rec.Code)
		assert.Equal(t, 0, f.called("SaveVisitor"))
	})

	t.Run("bad url", func(t *testing.T) {
		rec := post(t, endpoints.Visitors, newFakeRepo(), endpoints.Visitor{Code: "x", Url: "http://localhost"})
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid json in both formats", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, post(t, endpoints.Visitors, newFakeRepo(), "{bad").Code)
	})

	t.Run("options", func(t *testing.T) {
		assert.Equal(t, http.StatusOK,
			call(t, newFakeRepo(), endpoints.Visitors, "/x", http.MethodOptions, "/x", "", nil).Code)
	})
}

func TestVisitorOffline(t *testing.T) {
	storeId := newV4()
	f := newFakeRepo()
	f.OnCheckStoreCodeOffline = func(code, unit string) string {
		if code == "SP-OK" && unit == "u1" {
			return storeId
		}
		return ""
	}
	rec := post(t, endpoints.VisitorsOffline, f, endpoints.VisitorOffline{Code: "SP-OK", Unit: "u1", Info: "5"})
	require.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, []interface{}{"5", storeId}, f.lastCall("SaveVisitorOffline"))

	assert.Equal(t, http.StatusBadRequest,
		post(t, endpoints.VisitorsOffline, f, endpoints.VisitorOffline{Code: "nope", Unit: "u1"}).Code)
	assert.Equal(t, http.StatusBadRequest, post(t, endpoints.VisitorsOffline, f, "{bad").Code)
	assert.Equal(t, 1, f.called("SaveVisitorOffline"))
	assert.Equal(t, http.StatusOK,
		call(t, f, endpoints.VisitorsOffline, "/x", http.MethodOptions, "/x", "", nil).Code)
}
