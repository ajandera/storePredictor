package endpoints_test

import (
	"bytes"
	"main/endpoints"
	"main/repo"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ajandera/sp_model/rdbsClientData"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

// xlsx builds a spreadsheet whose Sheet1 holds the given rows.
func xlsx(t *testing.T, rows [][]interface{}) []byte {
	t.Helper()
	file := excelize.NewFile()
	for i, row := range rows {
		cell, err := excelize.CoordinatesToCellName(1, i+1)
		require.NoError(t, err)
		require.NoError(t, file.SetSheetRow("Sheet1", cell, &row))
	}
	buf, err := file.WriteToBuffer()
	require.NoError(t, err)
	return buf.Bytes()
}

func importRequest(t *testing.T, target, caller, field, filename string, content []byte) *http.Request {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	if field != "" {
		part, err := mw.CreateFormFile(field, filename)
		require.NoError(t, err)
		_, err = part.Write(content)
		require.NoError(t, err)
	}
	require.NoError(t, mw.Close())
	req := httptest.NewRequest(http.MethodPost, target, &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	if caller != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken(t, caller))
	}
	return req
}

func runImport(t *testing.T, h func(http.ResponseWriter, *http.Request, repo.Repository), f *fakeRepo, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	router := mux.NewRouter()
	router.HandleFunc("/import/{storeId}", func(w http.ResponseWriter, r *http.Request) { h(w, r, f) })
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestImportWarehouse(t *testing.T) {
	t.Chdir(t.TempDir()) // the handler stores the upload in the working directory
	f, caller, storeId, _ := supplierFixture(t)
	url := "/import/" + storeId
	h := endpoints.ImportWarehouse

	sheet := xlsx(t, [][]interface{}{
		{"P1", "Existing", 5},
		{"P2", "New", 7},
		{"short", "row"}, // fewer than 3 columns is skipped
	})
	f.OnGetProduct = func(code, _ string) rdbsClientData.Product {
		if code == "P1" {
			return rdbsClientData.Product{ProductCode: "P1"}
		}
		return rdbsClientData.Product{}
	}

	rec := runImport(t, h, f, importRequest(t, url, caller, "importFile", "w.xlsx", sheet))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	assert.Equal(t, []interface{}{"P1", "Existing", storeId, int8(5)}, f.lastCall("UpdateProduct"))
	assert.Equal(t, []interface{}{"P2", "New", int8(7), storeId}, f.lastCall("CreateProduct"))
	assert.Equal(t, 1, f.called("UpdateProduct"))
	assert.Equal(t, 1, f.called("CreateProduct"))

	t.Run("errors", func(t *testing.T) {
		t.Chdir(t.TempDir())
		assert.Equal(t, http.StatusUnauthorized,
			runImport(t, h, f, importRequest(t, url, "", "importFile", "w.xlsx", sheet)).Code)
		assert.Equal(t, http.StatusForbidden,
			runImport(t, h, f, importRequest(t, "/import/"+newV4(), caller, "importFile", "w.xlsx", sheet)).Code)
		assert.Equal(t, http.StatusBadRequest,
			runImport(t, h, f, importRequest(t, url, caller, "", "", nil)).Code)
		assert.Equal(t, http.StatusBadRequest,
			runImport(t, h, f, importRequest(t, url, caller, "importFile", "w.xlsx", []byte("not a spreadsheet"))).Code)
	})
}

func TestImportWarehouseSheetWithoutSheet1(t *testing.T) {
	t.Chdir(t.TempDir())
	f, caller, storeId, _ := supplierFixture(t)

	file := excelize.NewFile()
	_, err := file.NewSheet("Other")
	require.NoError(t, err)
	require.NoError(t, file.DeleteSheet("Sheet1"))
	buf, err := file.WriteToBuffer()
	require.NoError(t, err)

	rec := runImport(t, endpoints.ImportWarehouse, f, importRequest(t, "/import/"+storeId, caller, "importFile", "w.xlsx", buf.Bytes()))
	assert.Equal(t, 0, f.called("CreateProduct"))
	assert.NotEqual(t, http.StatusInternalServerError, rec.Code)
}

func TestImportInvoices(t *testing.T) {
	t.Chdir(t.TempDir())
	f, caller, storeId, _ := supplierFixture(t)
	url := "/import/" + storeId
	h := endpoints.ImportInvoices

	due := time.Date(2031, 3, 4, 0, 0, 0, 0, time.UTC)
	sheet := xlsx(t, [][]interface{}{
		{due.Format(time.RFC3339), 99.5, "EUR"},
		{"only", "two"},
	})
	rec := runImport(t, h, f, importRequest(t, url, caller, "importFile", "i.xlsx", sheet))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	require.Equal(t, 1, f.called("CreateInvoice"))
	args := f.lastCall("CreateInvoice")
	assert.True(t, due.Equal(args[0].(time.Time)))
	assert.Equal(t, 99.5, args[1])
	assert.Equal(t, "EUR", args[2])
	assert.Equal(t, storeId, args[3])

	assert.Equal(t, http.StatusUnauthorized, runImport(t, h, f, importRequest(t, url, "", "importFile", "i.xlsx", sheet)).Code)
	assert.Equal(t, http.StatusForbidden, runImport(t, h, f, importRequest(t, "/import/"+newV4(), caller, "importFile", "i.xlsx", sheet)).Code)
	assert.Equal(t, http.StatusBadRequest, runImport(t, h, f, importRequest(t, url, caller, "", "", nil)).Code)
	assert.Equal(t, http.StatusBadRequest, runImport(t, h, f, importRequest(t, url, caller, "importFile", "i.xlsx", []byte("junk"))).Code)
	assert.Equal(t, 1, f.called("CreateInvoice"))
}
