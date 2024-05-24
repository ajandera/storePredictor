package structs

import "time"

type ProductToStore struct {
	ProductCode string
	StoreId     string
	Quantity    int8
	DateToNeed  time.Time
	DateToOrder time.Time
}
