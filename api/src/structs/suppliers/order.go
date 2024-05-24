package suppliers

type Order struct {
	Template     string
	Recipient    string
	Subject      string
	ProductCodes []string
	Quantities   []int
}
