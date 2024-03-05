package models

type FpmHistoryFileZsd081 struct {
	ID          int    `gorm:"primary_key:auto_increment" json:"id"`
	UploadBy    string `json:"uploadBy"`
	FileName    string `json:"fileName"`
	StatusFile  string `json:"statusFile"`
	Description string `json:"description"`
}

func (FpmHistoryFileZsd081) TableName() string {
	return "sf_fpm_history_upload_zsd081"
}

type FpmDetailInputZsd081 struct {
	ID           int    `gorm:"primary_key:auto_increment" json:"id"`
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

func (FpmDetailInputZsd081) TableName() string {
	return "sf_dump_zsd081"
}

type FpmListFileProcess struct {
	ID       int    `gorm:"primary_key:auto_increment" json:"id"`
	FileName string `json:"fileName"`
}

func (FpmListFileProcess) TableName() string {
	return "sf_fpm_list_file_process"
}
