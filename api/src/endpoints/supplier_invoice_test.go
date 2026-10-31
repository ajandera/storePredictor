package endpoints_test

import (
	"main/endpoints"
	inv "main/structs/invoices"
	od "main/structs/openData"
	sp "main/structs/suppliers"
	"net/http"
	"testing"
	"time"

	"github.com/ajandera/sp_model/rdbsClientInfo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// supplierFixture: owner account with one store and one supplier in it.
func supplierFixture(t *testing.T) (f *fakeRepo, caller, storeId, supplierId string) {
	t.Helper()
	t.Setenv("JWT_SECRET", testSecret)
	f = newFakeRepo()
	caller = f.addAccount("a", "").Id.String()
	storeId, supplierId = newV4(), newV4()
	f.OnGetSupplier = func(id string) rdbsClientInfo.Suppliers {
		if id == supplierId {
			return rdbsClientInfo.Suppliers{StoreRefer: storeId}
		}
		return rdbsClientInfo.Suppliers{}
	}
	f.OnIsPermitted = func(_, s string) bool { return s == storeId }
	return
}

func TestSupplierGetList(t *testing.T) {
	f, caller, storeId, _ := supplierFixture(t)
	url := "/s/" + storeId
	p := "/s/{storeId}"
	assert.Equal(t, http.StatusOK, call(t, f, endpoints.GetSuppliers, p, http.MethodGet, url, caller, nil).Code)
	assert.Equal(t, http.StatusForbidden, call(t, f, endpoints.GetSuppliers, p, http.MethodGet, "/s/"+newV4(), caller, nil).Code)
	assert.Equal(t, http.StatusUnauthorized, call(t, f, endpoints.GetSuppliers, p, http.MethodGet, url, "", nil).Code)
}

func TestSupplierGetDetail(t *testing.T) {
	f, caller, _, supplierId := supplierFixture(t)
	p := "/d/{supplierId}"
	assert.Equal(t, http.StatusOK, call(t, f, endpoints.GetSupplierDetail, p, http.MethodGet, "/d/"+supplierId, caller, nil).Code)
	// unknown supplier -> no store -> not permitted
	assert.Equal(t, http.StatusForbidden, call(t, f, endpoints.GetSupplierDetail, p, http.MethodGet, "/d/"+newV4(), caller, nil).Code)
	assert.Equal(t, http.StatusForbidden, call(t, f, endpoints.GetSupplierDetail, p, http.MethodGet, "/d/bad", caller, nil).Code)
	assert.Equal(t, http.StatusUnauthorized, call(t, f, endpoints.GetSupplierDetail, p, http.MethodGet, "/d/"+supplierId, "", nil).Code)
}

func TestSupplierCreate(t *testing.T) {
	f, caller, storeId, _ := supplierFixture(t)
	h := endpoints.CreateSupplier
	call2 := func(c string, body interface{}) int {
		return call(t, f, h, "/s", http.MethodPost, "/s", c, body).Code
	}
	assert.Equal(t, http.StatusUnauthorized, call2("", sp.Suppliers{}))
	assert.Equal(t, http.StatusBadRequest, call2(caller, "{bad"))
	assert.Equal(t, http.StatusForbidden, call2(caller, sp.Suppliers{StoreRefer: newV4()}))
	assert.Equal(t, 0, f.called("CreateSupplier"))

	assert.Equal(t, http.StatusCreated, call2(caller, sp.Suppliers{Name: "ACME", StoreRefer: storeId, Email: "a@b.c"}))
	assert.Equal(t, "ACME", f.lastCall("CreateSupplier")[0])
}

func TestSupplierOrderSendsNothingWithoutSMTP(t *testing.T) {
	t.Setenv("SMTP_URL", "")
	f, caller, _, _ := supplierFixture(t)
	h := endpoints.CreateSupplierOrder
	assert.Equal(t, http.StatusUnauthorized, call(t, f, h, "/o", http.MethodPost, "/o", "", sp.Order{}).Code)
	assert.Equal(t, http.StatusBadRequest, call(t, f, h, "/o", http.MethodPost, "/o", caller, "{bad").Code)
	assert.Equal(t, http.StatusCreated, call(t, f, h, "/o", http.MethodPost, "/o", caller,
		sp.Order{Recipient: "x@y.z", Subject: "s", Template: "<b>hi</b>"}).Code)
}

func TestSupplierUpdate(t *testing.T) {
	f, caller, storeId, supplierId := supplierFixture(t)
	h := endpoints.UpdateSupplier
	put := func(c string, body interface{}) int {
		return call(t, f, h, "/s", http.MethodPut, "/s", c, body).Code
	}
	assert.Equal(t, http.StatusBadRequest, put(caller, "{bad"))
	assert.Equal(t, http.StatusUnauthorized, put("", sp.Suppliers{Id: supplierId}))
	assert.Equal(t, http.StatusForbidden, put(caller, sp.Suppliers{Id: newV4(), StoreRefer: newV4()}))
	assert.Equal(t, 0, f.called("UpdateSupplier"))

	// known supplier: store looked up from the db
	assert.Equal(t, http.StatusCreated, put(caller, sp.Suppliers{Id: supplierId, Name: "new"}))
	assert.Equal(t, supplierId, f.lastCall("UpdateSupplier")[0])
	// unknown supplier falls back to the store sent in body
	assert.Equal(t, http.StatusCreated, put(caller, sp.Suppliers{Id: newV4(), StoreRefer: storeId}))
	assert.Equal(t, 2, f.called("UpdateSupplier"))
}

func TestSupplierDelete(t *testing.T) {
	f, caller, _, supplierId := supplierFixture(t)
	h := endpoints.DeleteSupplier
	p := "/s/{supplierId:.*}"
	del := func(c, id string) int {
		return call(t, f, h, p, http.MethodDelete, "/s/"+id, c, nil).Code
	}
	assert.Equal(t, http.StatusUnauthorized, del("", supplierId))
	assert.Equal(t, http.StatusBadRequest, del(caller, ""))
	assert.Equal(t, http.StatusForbidden, del(caller, newV4()))
	assert.Equal(t, 0, f.called("DeleteSupplier"))
	assert.Equal(t, http.StatusOK, del(caller, supplierId))
	assert.Equal(t, []interface{}{supplierId}, f.lastCall("DeleteSupplier"))
}

// invoiceFixture: owner account with a store holding one invoice.
func invoiceFixture(t *testing.T) (f *fakeRepo, owner, child, stranger rdbsClientInfo.Accounts, storeId, invoiceId string) {
	t.Helper()
	t.Setenv("JWT_SECRET", testSecret)
	f = newFakeRepo()
	owner = f.addAccount("owner", "")
	child = f.addAccount("child", owner.Id.String())
	stranger = f.addAccount("stranger", "")
	storeId, invoiceId = newV4(), newV4()
	f.stores[owner.Id.String()] = []rdbsClientInfo.Stores{{Id: uuid.MustParse(storeId)}}
	f.OnGetInvoices = func(s string) []rdbsClientInfo.Invoices {
		if s == storeId {
			return []rdbsClientInfo.Invoices{{Id: uuid.MustParse(invoiceId)}}
		}
		return nil
	}
	f.OnIsPermitted = func(acc, s string) bool { return s == storeId && acc == owner.Id.String() }
	return
}

func TestInvoiceGetAndQuery(t *testing.T) {
	f, owner, _, stranger, storeId, _ := invoiceFixture(t)
	for _, c := range []struct {
		name, pattern, url string
		h                  repoHandler
	}{
		{"list", "/i/{storeId}", "/i/" + storeId, endpoints.GetInvoices},
		{"query", "/i/{storeId}/{from}/{to}", "/i/" + storeId + "/1/9", endpoints.GetInvoicesQuery},
	} {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, http.StatusOK, call(t, f, c.h, c.pattern, http.MethodGet, c.url, owner.Id.String(), nil).Code)
			assert.Equal(t, http.StatusForbidden, call(t, f, c.h, c.pattern, http.MethodGet, c.url, stranger.Id.String(), nil).Code)
			assert.Equal(t, http.StatusUnauthorized, call(t, f, c.h, c.pattern, http.MethodGet, c.url, "", nil).Code)
		})
	}
	args := f.lastCall("GetInvoicesFilter")
	assert.Equal(t, []interface{}{storeId, 1, 9}, args)
}

func TestInvoiceCreate(t *testing.T) {
	f, owner, _, stranger, storeId, _ := invoiceFixture(t)
	due := time.Date(2030, 1, 2, 0, 0, 0, 0, time.UTC)
	h := endpoints.CreateInvoice
	post := func(c string, body interface{}) int {
		return call(t, f, h, "/i", http.MethodPost, "/i", c, body).Code
	}
	assert.Equal(t, http.StatusUnauthorized, post("", inv.Invoices{}))
	assert.Equal(t, http.StatusBadRequest, post(owner.Id.String(), "{bad"))
	assert.Equal(t, http.StatusForbidden, post(stranger.Id.String(), inv.Invoices{StoreRefer: storeId}))
	assert.Equal(t, 0, f.called("CreateInvoice"))

	require.Equal(t, http.StatusCreated, post(owner.Id.String(),
		inv.Invoices{StoreRefer: storeId, Amount: "12.5", Currency: "EUR", DueDate: due}))
	args := f.lastCall("CreateInvoice")
	assert.Equal(t, 12.5, args[1])
	assert.Equal(t, "EUR", args[2])
	assert.Equal(t, storeId, args[3])
}

func TestInvoiceUpdateAndDeleteAreTenantScoped(t *testing.T) {
	f, owner, child, stranger, _, invoiceId := invoiceFixture(t)

	put := func(c string, id string) int {
		return call(t, f, endpoints.UpdateInvoice, "/i", http.MethodPut, "/i", c, inv.Invoices{Id: id, Amount: "3"}).Code
	}
	assert.Equal(t, http.StatusUnauthorized, put("", invoiceId))
	assert.Equal(t, http.StatusBadRequest,
		call(t, f, endpoints.UpdateInvoice, "/i", http.MethodPut, "/i", owner.Id.String(), "{bad").Code)
	assert.Equal(t, http.StatusForbidden, put(stranger.Id.String(), invoiceId))
	assert.Equal(t, http.StatusForbidden, put(owner.Id.String(), newV4()))
	assert.Equal(t, http.StatusForbidden, put(owner.Id.String(), "bad-id"))
	assert.Equal(t, 0, f.called("UpdateInvoice"))
	assert.Equal(t, http.StatusCreated, put(child.Id.String(), invoiceId)) // child shares tenant of owner
	assert.Equal(t, invoiceId, f.lastCall("UpdateInvoice")[0])

	del := func(c, id string) int {
		return call(t, f, endpoints.DeleteInvoice, "/i/{invoiceId:.*}", http.MethodDelete, "/i/"+id, c, nil).Code
	}
	assert.Equal(t, http.StatusUnauthorized, del("", invoiceId))
	assert.Equal(t, http.StatusBadRequest, del(owner.Id.String(), ""))
	assert.Equal(t, http.StatusForbidden, del(stranger.Id.String(), invoiceId))
	assert.Equal(t, 0, f.called("DeleteInvoice"))
	assert.Equal(t, http.StatusOK, del(owner.Id.String(), invoiceId))
	assert.Equal(t, []interface{}{invoiceId}, f.lastCall("DeleteInvoice"))
}

func TestOpenData(t *testing.T) {
	f, caller, storeId, _ := supplierFixture(t)
	p := "/od/{storeId:.*}"
	get := func(c, id string) int {
		return call(t, f, endpoints.GetOpenData, p, http.MethodGet, "/od/"+id, c, nil).Code
	}
	assert.Equal(t, http.StatusUnauthorized, get("", storeId))
	assert.Equal(t, http.StatusForbidden, get(caller, newV4()))
	assert.Equal(t, http.StatusOK, get(caller, storeId))

	add := func(c string, body interface{}) int {
		return call(t, f, endpoints.AddOpenData, "/od", http.MethodPost, "/od", c, body).Code
	}
	assert.Equal(t, http.StatusUnauthorized, add("", od.OpenData{}))
	assert.Equal(t, http.StatusBadRequest, add(caller, "{bad"))
	assert.Equal(t, http.StatusForbidden, add(caller, od.OpenData{StoreRefer: newV4()}))
	assert.Equal(t, http.StatusCreated, add(caller, od.OpenData{StoreRefer: storeId, StorePower: 2}))
	assert.Equal(t, 2.0, f.lastCall("CreateOpenData")[0])
}
