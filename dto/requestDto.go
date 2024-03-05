package dto

type ReqSaveDataUploadDto struct {
	JenisFakturPajak string   `json:"jenisFakturPajak" binding:"required" validate:"required"`
	FileNameUpload   []string `json:"fileNameUpload" binding:"required" validate:"required"`
}

type ReqDeleteFileDto struct {
	TypeFile string `json:"typeFile" binding:"required" validate:"required"`
	FileName string `json:"fileName" binding:"required" validate:"required"`
}

type ReqPaginationDto struct {
	Page           int    `json:"page" binding:"required" validate:"required"`
	MaxDataDisplay int    `json:"maxDataDisplay" binding:"required" validate:"required"`
	Search         string `json:"search,omiempty"`
}

type ReqSaveSourceDataDto struct {
	Files    []string `json:"files" binding:"required" validate:"required"`
	TypeFile string   `json:"typeFile" binding:"required" validate:"required"`
}

type ReqInputDate struct {
	DateInput string `json:"dateInput"`
}

type ReqDownloadFpmDto struct {
	IdFiles []int `json:"idFiles" binding:"required" validate:"required"`
}

type ReqNotifDto struct {
	RequestId      int      `json:"RequestId"`
	Partner        string   `json:"PartnerCode"`
	Timestamp      string   `json:"Timestamp"`
	Signature      string   `json:"Signature"`
	Recipient      string   `json:"Recipient"`
	Cc             string   `json:"Cc"`
	Bcc            string   `json:"Bcc"`
	Sender         string   `json:"Sender"`
	Body           string   `json:"Body"`
	Subject        string   `json:"Subject"`
	PathAttachment []string `json:"PathAttachment"`
	LogStatus      string   `json:"LogStatus"`
}

type ReqAddUserFpmDto struct {
	Username string `json:"username" binding:"required" validate:"required"`
	RoleUser string `json:"roleUser" binding:"required" validate:"required"`
}

type ReqDeleteUserFpmDto struct {
	IdUser int `json:"idUser" binding:"required" validate:"required"`
}

type ReqCollectionInvDto struct {
	Page           int    `json:"page" binding:"required" validate:"required"`
	MaxDataDisplay int    `json:"maxDataDisplay" binding:"required" validate:"required"`
	FilterCustomer string `json:"filterCustomer,omiempty"`
	FilterInvoice  string `json:"filterInvoice,omiempty"`
	StartBillDate  string `json:"startBillDate" binding:"required" validate:"required"`
	EndBillDate    string `json:"endBillDate" binding:"required" validate:"required"`
	TypeDc         string `json:"typeDc" binding:"required" validate:"required"`
}
