package structs

import "time"

type Account struct {
	ID                     string
	Name                   string
	Email                  string
	Street                 string
	City                   string
	Zip                    string
	CountryCode            string
	CompanyNumber          string
	VatNumber              string
	PaidTo                 string
	PlanRefer              string
	Role                   string
	Parent                 string
	Password               string
	Newsletter             bool
	NewsletterConfirmation time.Time
}
