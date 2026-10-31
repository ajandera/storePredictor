package endpoints_test

// Code generated from sp_model.Repository by the test helper generator; hooks are
// set per test through the On<Method> fields. Unset methods record the call and
// return zero values.

import (
	"time"

	"github.com/ajandera/sp_model/rdbsClientData"
	"github.com/ajandera/sp_model/rdbsClientInfo"
	"gorm.io/gorm"
)

var _ = time.Now

type fakeRepo struct {
	fakeState
	calls                          []fakeCall
	OnAuth                         func(string, string) rdbsClientInfo.Accounts
	OnCheckStoreCode               func(string, string) string
	OnCheckStoreCodeOffline        func(string, string) string
	OnCreateAccount                func(string, string, bool) *gorm.DB
	OnCreateInvoice                func(time.Time, float64, string, string) *gorm.DB
	OnCreateOpenData               func(float64, float64, float64, float64, float64, string) *gorm.DB
	OnCreateOrder                  func(string, string, string, float64, bool) string
	OnCreateProduct                func(string, string, int8, string) rdbsClientData.Products
	OnCreateStore                  func(string, string, string, string, bool, string, string, string, int8) *gorm.DB
	OnCreateStoreWeights           func(string, string, float64, float64, float64, float64, float64, float64, float64, float64, string, int, int) *gorm.DB
	OnCreateSupplier               func(string, string, string, string, string, string, string, string, string, string, string) *gorm.DB
	OnDeleteAccount                func(string)
	OnDeleteInvoice                func(string)
	OnDeleteStore                  func(string)
	OnDeleteSupplier               func(string)
	OnEditAccount                  func(string, string, string, string, string, string, string, string, string, string, string, string, string, string, bool) rdbsClientInfo.Accounts
	OnEditStore                    func(string, string, string, float64, float64, float64, float64, float64, int, bool, string, int8) rdbsClientInfo.Stores
	OnEditStoreWeights             func(string, string, float64, float64, float64, float64, float64, float64, float64, float64, string, int, int) rdbsClientInfo.StoreWeights
	OnGetAccountByEmail            func(string) rdbsClientInfo.Accounts
	OnGetAccountById               func(string) rdbsClientInfo.Accounts
	OnGetAccountOrders             func(string) []rdbsClientInfo.Orders
	OnGetChildAccountById          func(string) []rdbsClientInfo.Accounts
	OnGetInvoices                  func(string) []rdbsClientInfo.Invoices
	OnGetInvoicesFilter            func(string, int, int) []rdbsClientInfo.Invoices
	OnGetNumberOrders              func(string) float64
	OnGetOpenData                  func(string) []rdbsClientInfo.OpenData
	OnGetOrderById                 func(string) rdbsClientInfo.Orders
	OnGetOrders                    func(map[string]interface{}, int, int) []rdbsClientData.Orders
	OnGetOrdersWithProduct         func(string, string) []rdbsClientData.Orders
	OnGetPaidPlans                 func() []rdbsClientInfo.Plan
	OnGetPlans                     func() []rdbsClientInfo.Plan
	OnGetProduct                   func(string, string) rdbsClientData.Product
	OnGetProductToStore            func(string, string) rdbsClientData.ProductToStore
	OnGetProducts                  func(map[string]interface{}) []rdbsClientData.TopSellProduct
	OnGetProductsToStore           func(string, int, int) []rdbsClientData.ProductsToStore
	OnGetProductsWarehouse         func(string, int, int) []rdbsClientData.Product
	OnGetStoreById                 func(string) rdbsClientInfo.Stores
	OnGetStoreByUrl                func(string) (string, error)
	OnGetStoreWeights              func(string) rdbsClientInfo.StoreWeights
	OnGetStoresByAccount           func(string) []rdbsClientInfo.Stores
	OnGetSumOrder                  func(string) float64
	OnGetSumVisitors               func(string) float64
	OnGetSupplier                  func(string) rdbsClientInfo.Suppliers
	OnGetSuppliers                 func(string) []rdbsClientInfo.Suppliers
	OnGetVisitorsForPredictionView func(string, string, string) []rdbsClientData.VisitorsByDay
	OnIsPermitted                  func(string, string) bool
	OnSaveOrder                    func(float64, string, string, []rdbsClientData.Item, string, string) *gorm.DB
	OnSaveVisitor                  func(string, string, string, string, string, string) *gorm.DB
	OnSaveVisitorOffline           func(string, string) *gorm.DB
	OnSetRestorePw                 func(string, string) rdbsClientInfo.Accounts
	OnUpdateInvoice                func(string, time.Time, float64, string) rdbsClientInfo.Invoices
	OnUpdateProduct                func(string, string, string, int8) *gorm.DB
	OnUpdatePw                     func(string, string) rdbsClientInfo.Accounts
	OnUpdateSupplier               func(string, string, string, string, string, string, string, string, string, string, string) rdbsClientInfo.Suppliers
}

type fakeCall struct {
	Method string
	Args   []interface{}
}

func (f *fakeRepo) record(method string, args []interface{}) {
	f.calls = append(f.calls, fakeCall{Method: method, Args: args})
}

// called returns how many times the method was invoked.
func (f *fakeRepo) called(method string) int {
	n := 0
	for _, c := range f.calls {
		if c.Method == method {
			n++
		}
	}
	return n
}

// lastCall returns the arguments of the latest invocation of the method.
func (f *fakeRepo) lastCall(method string) []interface{} {
	for i := len(f.calls) - 1; i >= 0; i-- {
		if f.calls[i].Method == method {
			return f.calls[i].Args
		}
	}
	return nil
}

var _ = rdbsClientData.Item{}
var _ = rdbsClientInfo.Accounts{}
var _ *gorm.DB

func (f *fakeRepo) Auth(a0 string, a1 string) rdbsClientInfo.Accounts {
	f.record("Auth", []interface{}{a0, a1})
	if f.OnAuth != nil {
		return f.OnAuth(a0, a1)
	}
	var z0 rdbsClientInfo.Accounts
	return z0
}

func (f *fakeRepo) CheckStoreCode(a0 string, a1 string) string {
	f.record("CheckStoreCode", []interface{}{a0, a1})
	if f.OnCheckStoreCode != nil {
		return f.OnCheckStoreCode(a0, a1)
	}
	var z0 string
	return z0
}

func (f *fakeRepo) CheckStoreCodeOffline(a0 string, a1 string) string {
	f.record("CheckStoreCodeOffline", []interface{}{a0, a1})
	if f.OnCheckStoreCodeOffline != nil {
		return f.OnCheckStoreCodeOffline(a0, a1)
	}
	var z0 string
	return z0
}

func (f *fakeRepo) CreateAccount(a0 string, a1 string, a2 bool) *gorm.DB {
	f.record("CreateAccount", []interface{}{a0, a1, a2})
	if f.OnCreateAccount != nil {
		return f.OnCreateAccount(a0, a1, a2)
	}
	var z0 *gorm.DB
	return z0
}

func (f *fakeRepo) CreateInvoice(a0 time.Time, a1 float64, a2 string, a3 string) *gorm.DB {
	f.record("CreateInvoice", []interface{}{a0, a1, a2, a3})
	if f.OnCreateInvoice != nil {
		return f.OnCreateInvoice(a0, a1, a2, a3)
	}
	var z0 *gorm.DB
	return z0
}

func (f *fakeRepo) CreateOpenData(a0 float64, a1 float64, a2 float64, a3 float64, a4 float64, a5 string) *gorm.DB {
	f.record("CreateOpenData", []interface{}{a0, a1, a2, a3, a4, a5})
	if f.OnCreateOpenData != nil {
		return f.OnCreateOpenData(a0, a1, a2, a3, a4, a5)
	}
	var z0 *gorm.DB
	return z0
}

func (f *fakeRepo) CreateOrder(a0 string, a1 string, a2 string, a3 float64, a4 bool) string {
	f.record("CreateOrder", []interface{}{a0, a1, a2, a3, a4})
	if f.OnCreateOrder != nil {
		return f.OnCreateOrder(a0, a1, a2, a3, a4)
	}
	var z0 string
	return z0
}

func (f *fakeRepo) CreateProduct(a0 string, a1 string, a2 int8, a3 string) rdbsClientData.Products {
	f.record("CreateProduct", []interface{}{a0, a1, a2, a3})
	if f.OnCreateProduct != nil {
		return f.OnCreateProduct(a0, a1, a2, a3)
	}
	var z0 rdbsClientData.Products
	return z0
}

func (f *fakeRepo) CreateStore(a0 string, a1 string, a2 string, a3 string, a4 bool, a5 string, a6 string, a7 string, a8 int8) *gorm.DB {
	f.record("CreateStore", []interface{}{a0, a1, a2, a3, a4, a5, a6, a7, a8})
	if f.OnCreateStore != nil {
		return f.OnCreateStore(a0, a1, a2, a3, a4, a5, a6, a7, a8)
	}
	var z0 *gorm.DB
	return z0
}

func (f *fakeRepo) CreateStoreWeights(a0 string, a1 string, a2 float64, a3 float64, a4 float64, a5 float64, a6 float64, a7 float64, a8 float64, a9 float64, a10 string, a11 int, a12 int) *gorm.DB {
	f.record("CreateStoreWeights", []interface{}{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12})
	if f.OnCreateStoreWeights != nil {
		return f.OnCreateStoreWeights(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12)
	}
	var z0 *gorm.DB
	return z0
}

func (f *fakeRepo) CreateSupplier(a0 string, a1 string, a2 string, a3 string, a4 string, a5 string, a6 string, a7 string, a8 string, a9 string, a10 string) *gorm.DB {
	f.record("CreateSupplier", []interface{}{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10})
	if f.OnCreateSupplier != nil {
		return f.OnCreateSupplier(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
	}
	var z0 *gorm.DB
	return z0
}

func (f *fakeRepo) DeleteAccount(a0 string) {
	f.record("DeleteAccount", []interface{}{a0})
	if f.OnDeleteAccount != nil {
		f.OnDeleteAccount(a0)
		return
	}
}

func (f *fakeRepo) DeleteInvoice(a0 string) {
	f.record("DeleteInvoice", []interface{}{a0})
	if f.OnDeleteInvoice != nil {
		f.OnDeleteInvoice(a0)
		return
	}
}

func (f *fakeRepo) DeleteStore(a0 string) {
	f.record("DeleteStore", []interface{}{a0})
	if f.OnDeleteStore != nil {
		f.OnDeleteStore(a0)
		return
	}
}

func (f *fakeRepo) DeleteSupplier(a0 string) {
	f.record("DeleteSupplier", []interface{}{a0})
	if f.OnDeleteSupplier != nil {
		f.OnDeleteSupplier(a0)
		return
	}
}

func (f *fakeRepo) EditAccount(a0 string, a1 string, a2 string, a3 string, a4 string, a5 string, a6 string, a7 string, a8 string, a9 string, a10 string, a11 string, a12 string, a13 string, a14 bool) rdbsClientInfo.Accounts {
	f.record("EditAccount", []interface{}{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14})
	if f.OnEditAccount != nil {
		return f.OnEditAccount(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14)
	}
	var z0 rdbsClientInfo.Accounts
	return z0
}

func (f *fakeRepo) EditStore(a0 string, a1 string, a2 string, a3 float64, a4 float64, a5 float64, a6 float64, a7 float64, a8 int, a9 bool, a10 string, a11 int8) rdbsClientInfo.Stores {
	f.record("EditStore", []interface{}{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11})
	if f.OnEditStore != nil {
		return f.OnEditStore(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11)
	}
	var z0 rdbsClientInfo.Stores
	return z0
}

func (f *fakeRepo) EditStoreWeights(a0 string, a1 string, a2 float64, a3 float64, a4 float64, a5 float64, a6 float64, a7 float64, a8 float64, a9 float64, a10 string, a11 int, a12 int) rdbsClientInfo.StoreWeights {
	f.record("EditStoreWeights", []interface{}{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12})
	if f.OnEditStoreWeights != nil {
		return f.OnEditStoreWeights(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12)
	}
	var z0 rdbsClientInfo.StoreWeights
	return z0
}

func (f *fakeRepo) GetAccountByEmail(a0 string) rdbsClientInfo.Accounts {
	f.record("GetAccountByEmail", []interface{}{a0})
	if f.OnGetAccountByEmail != nil {
		return f.OnGetAccountByEmail(a0)
	}
	var z0 rdbsClientInfo.Accounts
	return z0
}

func (f *fakeRepo) GetAccountById(a0 string) rdbsClientInfo.Accounts {
	f.record("GetAccountById", []interface{}{a0})
	if f.OnGetAccountById != nil {
		return f.OnGetAccountById(a0)
	}
	var z0 rdbsClientInfo.Accounts
	return z0
}

func (f *fakeRepo) GetAccountOrders(a0 string) []rdbsClientInfo.Orders {
	f.record("GetAccountOrders", []interface{}{a0})
	if f.OnGetAccountOrders != nil {
		return f.OnGetAccountOrders(a0)
	}
	var z0 []rdbsClientInfo.Orders
	return z0
}

func (f *fakeRepo) GetChildAccountById(a0 string) []rdbsClientInfo.Accounts {
	f.record("GetChildAccountById", []interface{}{a0})
	if f.OnGetChildAccountById != nil {
		return f.OnGetChildAccountById(a0)
	}
	var z0 []rdbsClientInfo.Accounts
	return z0
}

func (f *fakeRepo) GetInvoices(a0 string) []rdbsClientInfo.Invoices {
	f.record("GetInvoices", []interface{}{a0})
	if f.OnGetInvoices != nil {
		return f.OnGetInvoices(a0)
	}
	var z0 []rdbsClientInfo.Invoices
	return z0
}

func (f *fakeRepo) GetInvoicesFilter(a0 string, a1 int, a2 int) []rdbsClientInfo.Invoices {
	f.record("GetInvoicesFilter", []interface{}{a0, a1, a2})
	if f.OnGetInvoicesFilter != nil {
		return f.OnGetInvoicesFilter(a0, a1, a2)
	}
	var z0 []rdbsClientInfo.Invoices
	return z0
}

func (f *fakeRepo) GetNumberOrders(a0 string) float64 {
	f.record("GetNumberOrders", []interface{}{a0})
	if f.OnGetNumberOrders != nil {
		return f.OnGetNumberOrders(a0)
	}
	var z0 float64
	return z0
}

func (f *fakeRepo) GetOpenData(a0 string) []rdbsClientInfo.OpenData {
	f.record("GetOpenData", []interface{}{a0})
	if f.OnGetOpenData != nil {
		return f.OnGetOpenData(a0)
	}
	var z0 []rdbsClientInfo.OpenData
	return z0
}

func (f *fakeRepo) GetOrderById(a0 string) rdbsClientInfo.Orders {
	f.record("GetOrderById", []interface{}{a0})
	if f.OnGetOrderById != nil {
		return f.OnGetOrderById(a0)
	}
	var z0 rdbsClientInfo.Orders
	return z0
}

func (f *fakeRepo) GetOrders(a0 map[string]interface{}, a1 int, a2 int) []rdbsClientData.Orders {
	f.record("GetOrders", []interface{}{a0, a1, a2})
	if f.OnGetOrders != nil {
		return f.OnGetOrders(a0, a1, a2)
	}
	var z0 []rdbsClientData.Orders
	return z0
}

func (f *fakeRepo) GetOrdersWithProduct(a0 string, a1 string) []rdbsClientData.Orders {
	f.record("GetOrdersWithProduct", []interface{}{a0, a1})
	if f.OnGetOrdersWithProduct != nil {
		return f.OnGetOrdersWithProduct(a0, a1)
	}
	var z0 []rdbsClientData.Orders
	return z0
}

func (f *fakeRepo) GetPaidPlans() []rdbsClientInfo.Plan {
	f.record("GetPaidPlans", nil)
	if f.OnGetPaidPlans != nil {
		return f.OnGetPaidPlans()
	}
	var z0 []rdbsClientInfo.Plan
	return z0
}

func (f *fakeRepo) GetPlans() []rdbsClientInfo.Plan {
	f.record("GetPlans", nil)
	if f.OnGetPlans != nil {
		return f.OnGetPlans()
	}
	var z0 []rdbsClientInfo.Plan
	return z0
}

func (f *fakeRepo) GetProduct(a0 string, a1 string) rdbsClientData.Product {
	f.record("GetProduct", []interface{}{a0, a1})
	if f.OnGetProduct != nil {
		return f.OnGetProduct(a0, a1)
	}
	var z0 rdbsClientData.Product
	return z0
}

func (f *fakeRepo) GetProductToStore(a0 string, a1 string) rdbsClientData.ProductToStore {
	f.record("GetProductToStore", []interface{}{a0, a1})
	if f.OnGetProductToStore != nil {
		return f.OnGetProductToStore(a0, a1)
	}
	var z0 rdbsClientData.ProductToStore
	return z0
}

func (f *fakeRepo) GetProducts(a0 map[string]interface{}) []rdbsClientData.TopSellProduct {
	f.record("GetProducts", []interface{}{a0})
	if f.OnGetProducts != nil {
		return f.OnGetProducts(a0)
	}
	var z0 []rdbsClientData.TopSellProduct
	return z0
}

func (f *fakeRepo) GetProductsToStore(a0 string, a1 int, a2 int) []rdbsClientData.ProductsToStore {
	f.record("GetProductsToStore", []interface{}{a0, a1, a2})
	if f.OnGetProductsToStore != nil {
		return f.OnGetProductsToStore(a0, a1, a2)
	}
	var z0 []rdbsClientData.ProductsToStore
	return z0
}

func (f *fakeRepo) GetProductsWarehouse(a0 string, a1 int, a2 int) []rdbsClientData.Product {
	f.record("GetProductsWarehouse", []interface{}{a0, a1, a2})
	if f.OnGetProductsWarehouse != nil {
		return f.OnGetProductsWarehouse(a0, a1, a2)
	}
	var z0 []rdbsClientData.Product
	return z0
}

func (f *fakeRepo) GetStoreById(a0 string) rdbsClientInfo.Stores {
	f.record("GetStoreById", []interface{}{a0})
	if f.OnGetStoreById != nil {
		return f.OnGetStoreById(a0)
	}
	var z0 rdbsClientInfo.Stores
	return z0
}

func (f *fakeRepo) GetStoreByUrl(a0 string) (string, error) {
	f.record("GetStoreByUrl", []interface{}{a0})
	if f.OnGetStoreByUrl != nil {
		return f.OnGetStoreByUrl(a0)
	}
	var z0 string
	var z1 error
	return z0, z1
}

func (f *fakeRepo) GetStoreWeights(a0 string) rdbsClientInfo.StoreWeights {
	f.record("GetStoreWeights", []interface{}{a0})
	if f.OnGetStoreWeights != nil {
		return f.OnGetStoreWeights(a0)
	}
	var z0 rdbsClientInfo.StoreWeights
	return z0
}

func (f *fakeRepo) GetStoresByAccount(a0 string) []rdbsClientInfo.Stores {
	f.record("GetStoresByAccount", []interface{}{a0})
	if f.OnGetStoresByAccount != nil {
		return f.OnGetStoresByAccount(a0)
	}
	var z0 []rdbsClientInfo.Stores
	return z0
}

func (f *fakeRepo) GetSumOrder(a0 string) float64 {
	f.record("GetSumOrder", []interface{}{a0})
	if f.OnGetSumOrder != nil {
		return f.OnGetSumOrder(a0)
	}
	var z0 float64
	return z0
}

func (f *fakeRepo) GetSumVisitors(a0 string) float64 {
	f.record("GetSumVisitors", []interface{}{a0})
	if f.OnGetSumVisitors != nil {
		return f.OnGetSumVisitors(a0)
	}
	var z0 float64
	return z0
}

func (f *fakeRepo) GetSupplier(a0 string) rdbsClientInfo.Suppliers {
	f.record("GetSupplier", []interface{}{a0})
	if f.OnGetSupplier != nil {
		return f.OnGetSupplier(a0)
	}
	var z0 rdbsClientInfo.Suppliers
	return z0
}

func (f *fakeRepo) GetSuppliers(a0 string) []rdbsClientInfo.Suppliers {
	f.record("GetSuppliers", []interface{}{a0})
	if f.OnGetSuppliers != nil {
		return f.OnGetSuppliers(a0)
	}
	var z0 []rdbsClientInfo.Suppliers
	return z0
}

func (f *fakeRepo) GetVisitorsForPredictionView(a0 string, a1 string, a2 string) []rdbsClientData.VisitorsByDay {
	f.record("GetVisitorsForPredictionView", []interface{}{a0, a1, a2})
	if f.OnGetVisitorsForPredictionView != nil {
		return f.OnGetVisitorsForPredictionView(a0, a1, a2)
	}
	var z0 []rdbsClientData.VisitorsByDay
	return z0
}

func (f *fakeRepo) IsPermitted(a0 string, a1 string) bool {
	f.record("IsPermitted", []interface{}{a0, a1})
	if f.OnIsPermitted != nil {
		return f.OnIsPermitted(a0, a1)
	}
	var z0 bool
	return z0
}

func (f *fakeRepo) SaveOrder(a0 float64, a1 string, a2 string, a3 []rdbsClientData.Item, a4 string, a5 string) *gorm.DB {
	f.record("SaveOrder", []interface{}{a0, a1, a2, a3, a4, a5})
	if f.OnSaveOrder != nil {
		return f.OnSaveOrder(a0, a1, a2, a3, a4, a5)
	}
	var z0 *gorm.DB
	return z0
}

func (f *fakeRepo) SaveVisitor(a0 string, a1 string, a2 string, a3 string, a4 string, a5 string) *gorm.DB {
	f.record("SaveVisitor", []interface{}{a0, a1, a2, a3, a4, a5})
	if f.OnSaveVisitor != nil {
		return f.OnSaveVisitor(a0, a1, a2, a3, a4, a5)
	}
	var z0 *gorm.DB
	return z0
}

func (f *fakeRepo) SaveVisitorOffline(a0 string, a1 string) *gorm.DB {
	f.record("SaveVisitorOffline", []interface{}{a0, a1})
	if f.OnSaveVisitorOffline != nil {
		return f.OnSaveVisitorOffline(a0, a1)
	}
	var z0 *gorm.DB
	return z0
}

func (f *fakeRepo) SetRestorePw(a0 string, a1 string) rdbsClientInfo.Accounts {
	f.record("SetRestorePw", []interface{}{a0, a1})
	if f.OnSetRestorePw != nil {
		return f.OnSetRestorePw(a0, a1)
	}
	var z0 rdbsClientInfo.Accounts
	return z0
}

func (f *fakeRepo) UpdateInvoice(a0 string, a1 time.Time, a2 float64, a3 string) rdbsClientInfo.Invoices {
	f.record("UpdateInvoice", []interface{}{a0, a1, a2, a3})
	if f.OnUpdateInvoice != nil {
		return f.OnUpdateInvoice(a0, a1, a2, a3)
	}
	var z0 rdbsClientInfo.Invoices
	return z0
}

func (f *fakeRepo) UpdateProduct(a0 string, a1 string, a2 string, a3 int8) *gorm.DB {
	f.record("UpdateProduct", []interface{}{a0, a1, a2, a3})
	if f.OnUpdateProduct != nil {
		return f.OnUpdateProduct(a0, a1, a2, a3)
	}
	var z0 *gorm.DB
	return z0
}

func (f *fakeRepo) UpdatePw(a0 string, a1 string) rdbsClientInfo.Accounts {
	f.record("UpdatePw", []interface{}{a0, a1})
	if f.OnUpdatePw != nil {
		return f.OnUpdatePw(a0, a1)
	}
	var z0 rdbsClientInfo.Accounts
	return z0
}

func (f *fakeRepo) UpdateSupplier(a0 string, a1 string, a2 string, a3 string, a4 string, a5 string, a6 string, a7 string, a8 string, a9 string, a10 string) rdbsClientInfo.Suppliers {
	f.record("UpdateSupplier", []interface{}{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10})
	if f.OnUpdateSupplier != nil {
		return f.OnUpdateSupplier(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
	}
	var z0 rdbsClientInfo.Suppliers
	return z0
}
