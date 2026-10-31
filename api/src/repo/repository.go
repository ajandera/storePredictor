// Package repo defines the data access contract used by the HTTP handlers.
// It is satisfied by model.Repository from sp_model and allows handlers to be
// tested with in-memory fakes instead of a real database.
package repo

import (
	model "github.com/ajandera/sp_model"
	"time"

	"github.com/ajandera/sp_model/rdbsClientData"
	"github.com/ajandera/sp_model/rdbsClientInfo"
	"gorm.io/gorm"
)

// Repository lists the data operations the API handlers rely on.
type Repository interface {
	Auth(email string, password string) rdbsClientInfo.Accounts
	CheckStoreCode(code string, url string) string
	CheckStoreCodeOffline(code string, url string) string
	CreateAccount(email string, password string, newsletter bool) *gorm.DB
	CreateInvoice(dueDate time.Time, amount float64, currency string, storeRefer string) *gorm.DB
	CreateOpenData(storePower float64, customerSatisfaction float64, maximalProductPrice float64, minimalProductPrice float64, perceivedValue float64, storeRefer string) *gorm.DB
	CreateOrder(accountRefer string, storeRefer string, planRefer string, amount float64, paid bool) string
	CreateProduct(productCode string, name string, quantity int8, storeId string) rdbsClientData.Products
	CreateStore(countryCode string, url string, code string, accountRefer string, offline bool, shoptetId string, shoptetToken string, feed string, window int8) *gorm.DB
	CreateStoreWeights(storeRefer string, name string, beta float64, gama float64, delta float64, a float64, b float64, c float64, d float64, e float64, probabilityWeights string, shift int, longShift int) *gorm.DB
	CreateSupplier(name string, street string, city string, zip string, country string, email string, phone string, person string, storeRefer string, template string, subject string) *gorm.DB
	DeleteAccount(id string)
	DeleteInvoice(id string)
	DeleteStore(id string)
	DeleteSupplier(id string)
	EditAccount(id string, name string, email string, street string, city string, zip string, countryCode string, companyNumber string, vatNumber string, paidTo string, planRefer string, role string, parent string, password string, newsletter bool) rdbsClientInfo.Accounts
	EditStore(id string, countryCode string, url string, maximalProductPrice float64, minimalProductPrice float64, actualStorePower float64, actualCustomerSatisfaction float64, perceivedValue float64, productSell int, offline bool, feed string, window int8) rdbsClientInfo.Stores
	EditStoreWeights(storeRefer string, name string, beta float64, gama float64, delta float64, a float64, b float64, c float64, d float64, e float64, probabilityWeights string, shift int, longShift int) rdbsClientInfo.StoreWeights
	GetAccountByEmail(email string) rdbsClientInfo.Accounts
	GetAccountById(accountId string) rdbsClientInfo.Accounts
	GetAccountOrders(accountId string) []rdbsClientInfo.Orders
	GetChildAccountById(accountId string) []rdbsClientInfo.Accounts
	GetInvoices(storeId string) []rdbsClientInfo.Invoices
	GetInvoicesFilter(storeId string, from int, to int) []rdbsClientInfo.Invoices
	GetNumberOrders(storeId string) float64
	GetOpenData(storeRefer string) []rdbsClientInfo.OpenData
	GetOrderById(id string) rdbsClientInfo.Orders
	GetOrders(condition map[string]interface{}, limit int, offset int) []rdbsClientData.Orders
	GetOrdersWithProduct(productCode string, storeId string) []rdbsClientData.Orders
	GetPaidPlans() []rdbsClientInfo.Plan
	GetPlans() []rdbsClientInfo.Plan
	GetProduct(productCode string, storeId string) rdbsClientData.Product
	GetProductToStore(productCode string, storeId string) rdbsClientData.ProductToStore
	GetProducts(condition map[string]interface{}) []rdbsClientData.TopSellProduct
	GetProductsToStore(storeId string, limit int, offset int) []rdbsClientData.ProductsToStore
	GetProductsWarehouse(storeId string, limit int, offset int) []rdbsClientData.Product
	GetStoreById(storeId string) rdbsClientInfo.Stores
	GetStoreByUrl(url string) (string, error)
	GetStoreWeights(storeId string) rdbsClientInfo.StoreWeights
	GetStoresByAccount(accountId string) []rdbsClientInfo.Stores
	GetSumOrder(storeId string) float64
	GetSumVisitors(storeId string) float64
	GetSupplier(supplierId string) rdbsClientInfo.Suppliers
	GetSuppliers(storeId string) []rdbsClientInfo.Suppliers
	GetVisitorsForPredictionView(from string, to string, store string) []rdbsClientData.VisitorsByDay
	IsPermitted(accountId string, storeId string) bool
	SaveOrder(amount float64, currency string, storeId string, orderItems []rdbsClientData.Item, orderId string, tag string) *gorm.DB
	SaveVisitor(ip string, storeId string, url string, header string, productCode string, tag string) *gorm.DB
	SaveVisitorOffline(info string, storeId string) *gorm.DB
	SetRestorePw(id string, token string) rdbsClientInfo.Accounts
	UpdateInvoice(id string, dueDate time.Time, amount float64, currency string) rdbsClientInfo.Invoices
	UpdateProduct(productCode string, name string, storeId string, quantity int8) *gorm.DB
	UpdatePw(token string, password string) rdbsClientInfo.Accounts
	UpdateSupplier(id string, name string, street string, city string, zip string, country string, email string, phone string, person string, template string, subject string) rdbsClientInfo.Suppliers
}

// Compile-time guarantee that the real model repository satisfies the contract.
var _ Repository = model.Repository{}
