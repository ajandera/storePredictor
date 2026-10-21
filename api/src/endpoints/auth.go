package endpoints

import (
	"main/utils"
	"net/http"

	model "github.com/ajandera/sp_model"
	"github.com/ajandera/sp_model/rdbsClientInfo"
)

// authenticate verifies the access token from the Authorization header and
// returns the id of the logged account. On failure it writes 401 and returns false.
func authenticate(w http.ResponseWriter, r *http.Request) (string, bool) {
	accountId, err := utils.ParseToken(utils.BearerToken(r), utils.TokenTypeAccess)
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "Not Authorized")
		return "", false
	}
	return accountId, true
}

// isAuthorized checks the access token and that the store belongs to the logged account
// (or to its parent account). It writes 401/403 and returns false when access is denied.
func isAuthorized(w http.ResponseWriter, r *http.Request, m model.Repository, storeId string) bool {
	accountId, ok := authenticate(w, r)
	if !ok {
		return false
	}
	if !utils.IsValidUUID(storeId) || !m.IsPermitted(accountId, storeId) {
		utils.WriteJSONError(w, http.StatusForbidden, "Not Permitted")
		return false
	}
	return true
}

// tenantOwner returns id of the main account the given account belongs to
// (child accounts share the data of their parent).
func tenantOwner(account rdbsClientInfo.Accounts) string {
	if account.Parent != "" {
		return account.Parent
	}
	return account.Id.String()
}

// canViewAccount returns true when target account belongs to the same tenant as the caller.
func canViewAccount(m model.Repository, callerId string, targetId string) bool {
	if !utils.IsValidUUID(targetId) {
		return false
	}
	if callerId == targetId {
		return true
	}
	caller := m.GetAccountById(callerId)
	target := m.GetAccountById(targetId)
	if !utils.IsValidUUID(caller.Id.String()) || !utils.IsValidUUID(target.Id.String()) {
		return false
	}
	return tenantOwner(caller) == tenantOwner(target)
}

// canManageAccount returns true when the caller may change or delete the target account:
// it is the caller itself or a child account created by the caller.
func canManageAccount(m model.Repository, callerId string, targetId string) bool {
	if !utils.IsValidUUID(targetId) {
		return false
	}
	if callerId == targetId {
		return true
	}
	target := m.GetAccountById(targetId)
	return utils.IsValidUUID(target.Id.String()) && target.Parent == callerId
}

// denyAccess writes 403 response.
func denyAccess(w http.ResponseWriter) {
	utils.WriteJSONError(w, http.StatusForbidden, "Not Permitted")
}

// invoiceBelongsToCaller returns true when the invoice belongs to one of the stores
// of the caller's tenant (the model has no lookup of invoice by id).
func invoiceBelongsToCaller(m model.Repository, callerId string, invoiceId string) bool {
	if !utils.IsValidUUID(invoiceId) {
		return false
	}
	owner := tenantOwner(m.GetAccountById(callerId))
	for _, store := range m.GetStoresByAccount(owner) {
		for _, invoice := range m.GetInvoices(store.Id.String()) {
			if invoice.Id.String() == invoiceId {
				return true
			}
		}
	}
	return false
}

// supplierStore returns id of the store the supplier belongs to.
func supplierStore(m model.Repository, supplierId string) string {
	if !utils.IsValidUUID(supplierId) {
		return ""
	}
	return m.GetSupplier(supplierId).StoreRefer
}
