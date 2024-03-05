package dto

type DtoFpmZv60 struct {
	ID            int    `json:"id"`
	BillingNumber string `json:"billingNumber"`
	BillingDate   string `json:"billingDate"`
	Payer         string `json:"payer"`
	Name          string `json:"name"`
	Npwp          string `json:"npwp"`
	FpNumber      string `json:"fpNumber"`
	Item          string `json:"item"`
	Material      string `json:"material"`
	BilledQty     string `json:"billedQty"`
	Dpp           string `json:"dpp"`
	CurrDpp       string `json:"currDpp"`
	Ppn           string `json:"ppn"`
	CurrPpn       string `json:"currPpn"`
	Total         string `json:"total"`
	CurrTotal     string `json:"currTotal"`
	Plant         string `json:"plant"`
	FpCreatedBy   string `json:"fpCreatedBy"`
	FpCreatedDate string `json:"fpCreatedDate"`
	PriceListType string `json:"priceListType"`
	TaxClass      string `json:"taxClass"`
	FpBranchCod   string `json:"fpBranchCod"`
}
