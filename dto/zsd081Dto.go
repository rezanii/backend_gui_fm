package dto

type HistoryZsd081 struct {
	ID         int    `json:"id"`
	UploadBy   string `json:"uploadBy"`
	FileName   string `json:"fileName"`
	DateInsert string `json:"dateInsert"`
	DateUpdate string `json:"dateUpdate"`
	StatusFile string `json:"statusFile"`
}

type FileProcessZsd081 struct {
	ID       int    `gorm:"primary_key:auto_increment" json:"id"`
	FileName string `json:"fileName"`
}

type DtoFpmDetailInputZsd081 struct {
	ID           int    `json:"id"`
	IDHistory    int    `json:"idHistory"`
	DocNumber    string `json:"docNumber"`
	EmailAddress string `json:"emailAddress"`
	CreateOn     string `json:"createOn"`
	Time         string `json:"time"`
	CreateBy     string `json:"createBy"`
	EmailToOrCc  string `json:"emailToOrCc"`
	Payer        string `json:"payer"`
	PayerName    string `json:"payerName"`
	ShipTo       string `json:"ShipTo"`
	ShipToName   string `json:"ShipToName"`
	Type         string `json:"Type"`
}
