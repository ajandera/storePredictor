package structs

import "time"

type Stores struct {
	ID                         string
	CountryCode                string
	LastPrediction             time.Time
	Url                        string
	MaximalProductPrice        float64
	MinimalProductPrice        float64
	ActualStorePower           float64
	ActualCustomerSatisfaction float64
	PerceivedValue             float64
	ProductSell                int
	Code                       string
	AccountRefer               string
	Feed                       string
	Window                     int8
}
