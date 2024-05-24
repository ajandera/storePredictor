package invoices

import "time"

type Invoices struct {
	Id         string
	StoreRefer string
	DueDate    time.Time
	Amount     string
	Currency   string
}
