package structs

type Xml struct {
	Shopitem []Shopitem `xml:"SHOPITEM"`
}

type Shopitem struct {
	Product      string  `xml:"PRODUCT"`
	Description  string  `xml:"DESCRIPTION"`
	Price        string  `xml:"PRICE"`
	Url          string  `xml:"URL"`
	Imgurl       string  `xml:"IMGURL"`
	PriceVat     float32 `xml:"PRICE_VAT"`
	Manufacturer string  `xml:"MANUFACTURER"`
	Param        []Param
}

type Param struct {
	ParamName string `xml:"PARAM_NAME"`
	Val       string `xml:"VAT"`
}
