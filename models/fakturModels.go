package models

type FakturModels struct {
	ID             int    `gorm:"primary_key:auto_increment" json:"id"`
	InvNumber      string `json:"invNumber"`
	InvDc          string `json:"invDc"`
	Faktur         string `json:"faktur"`
	StatusSendFp   string `json:"statusSendFp"`
	Status         int    `json:"status"`
	EmailReceiptFp string `json:"emailReceiptFp"`
	SendDateFp     string `json:"sendDateFp"`
	SendTimeFp     string `json:"sendTimeFp"`
	FakturCreated  string `json:"fakturCreated"`
}

func (FakturModels) TableName() string {
	return "sf_fpm_faktur"
}
