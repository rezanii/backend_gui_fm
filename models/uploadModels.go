package models

type FpmHistoryFilePjk struct {
	ID               int    `gorm:"primary_key:auto_increment" json:"id"`
	UploadBy         string `json:"uploadBy"`
	JenisFakturPajak string `json:"jenisFakturPajak"`
	FileNameUpload   string `json:"fileNameUpload"`
	Url              string `json:"url"`
	NoFaktur         string `json:"noFaktur"`
	Status           string `json:"status"`
}

func (FpmHistoryFilePjk) TableName() string {
	return "sf_fpm_history_upload_pjk"
}
