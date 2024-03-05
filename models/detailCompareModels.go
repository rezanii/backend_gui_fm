package models

type DetailCompareModels struct {
	Id                     int    `gorm:"primary_key:auto_increment" json:"id"`
	BillingDocumentZsd001n string `json:"billingDocumentZsd001n"`
	BillingTypeZsd001n     string `json:"billingTypeZsd001n"`
	BillingDateZsd001n     string `json:"billingDateZsd001n"`
	DcZsd001n              string `json:"dcZsd001n"`
	PayerZsd001n           string `json:"payerZsd001n"`
	PayerNameZsd001n       string `json:"payerNameZsd001n"`
	SlorZsd001n            string `json:"slorZsd001n"`
	DistrictZsd001n        string `json:"districtZsd001n"`
	MaterialNumberZsd001n  string `json:"materialNumberZsd001n"`
	MaterialDescZsd001n    string `json:"materialDescZsd001n"`
	CreateOnZsd001n        string `json:"createOnZsd001n"`
	BillingDocCancel       string `json:"billingDocCancel"`
	StsCancelInv           string `json:"stsCancelInv"`
	StsEmailInvCancel      string `json:"stsEmailInvCancel"`
	StsSendInv             string `json:"stsSendInv"`
	BillingNumberZv60      string `json:"billingNumberZv60"`
	BillingDateZv60        string `json:"billingDateZv60"`
	FpNumberZv60           string `json:"fpNumberZv60"`
	FpCreatedDateZv60      string `json:"fpCreatedDateZv60"`
	PayerZv60              string `json:"payerZv60"`
	NameZv60               string `json:"nameZv60"`
	NpwpZv60               string `json:"npwpZv60"`
	MaterialZv60           string `json:"materialZv60"`
	StsCompare             string `json:"stsCompare"`
	IdHistoryFpm           int    `json:"idHistoryFpm"`
	StsEmailCompare        string `json:"stsEmailCompare"`
	EmailReceiptInv        string `json:"emailReceiptInv"`
	SendDateInv            string `json:"sendDateInv"`
	SendTimeInv            string `json:"sendTimeInv"`
	IdHisEmail             int    `json:"idHisEmail"`
}

func (DetailCompareModels) TableName() string {
	return "sf_fpm_compare"
}
