package structs

type Plans struct {
	ID       string
	Price    float64
	Period   int
	Name     string
	Products int
	Enabled  bool
	Free     bool
	OneTime  bool
}
