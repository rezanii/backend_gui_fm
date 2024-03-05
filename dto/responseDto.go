package dto

import "time"

type ResponseDto struct {
	Status  bool        `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    string      `json:"code"`
}

func SuccessResponse(data interface{}, msg string, code string) ResponseDto {
	return ResponseDto{
		Status:  true,
		Data:    data,
		Message: msg,
		Code:    code,
	}
}

func ErrorResponse(msgError string, code string) ResponseDto {
	return ResponseDto{
		Status:  false,
		Data:    nil,
		Message: msgError,
		Code:    code,
	}
}

type RespTempUploadDto struct {
	Total           int      `json:"total"`
	Urls            []string `json:"urls"`
	FileNamesUpload []string `json:"fileNamesUpload"`
}
type RespTempUploadDtoNew struct {
	Total int                 `json:"total"`
	Data  []ResSaveUploadFile `json:"dataUpload"`
}

type ResPaginationDto struct {
	BaseUrl   string      `json:"baseUrl"`
	TotalData int64       `json:"totalData"`
	Record    interface{} `json:"record"`
}

type RespFpmHistoryDto struct {
	ID               int    `json:"id"`
	UploadBy         string `json:"uploadBy"`
	JenisFakturPajak string `json:"jenisFakturPajak"`
	FileNameUpload   string `json:"fileNameUpload"`
	Url              string `json:"url"`
	NoFaktur         string `json:"noFaktur"`
	Status           string `json:"status"`
	CreateAt         string `json:"createAt"`
	UpdateAt         string `json:"updateAt"`
}

type RespHistoryZsd081Dto struct {
	ID          int    `json:"id"`
	UploadBy    string `json:"uploadBy"`
	FileName    string `json:"fileName"`
	StatusFile  string `json:"statusFile"`
	DateInsert  string `json:"dateInsert"`
	DateUpdate  string `json:"dateUpdate"`
	Description string `json:"description"`
}

type ResSaveUploadFile struct {
	FileName string
	Url      string
}

type ResNotifDto struct {
	ResponseCode int    `json:"ResponseCode"`
	ResponseDesc string `json:"ResponseDesc"`
}

type ResSaveEmailDto struct {
	Id             int       `json:"id"`
	Recipient      string    `json:"recipient"`
	Cc             string    `json:"cc"`
	Bcc            string    `json:"bcc"`
	Sender         string    `json:"sender"`
	Subject        string    `json:"subject"`
	PathAttachment string    `json:"pathAttachment"`
	LogStatus      int       `json:"logStatus"`
	Body           string    `json:"body"`
	SendDateTime   time.Time `json:"sendDateTime"`
}

type RespSaveFpmUserDto struct {
	Username   string `json:"username"`
	IdRoleUser int    `json:"idRoleUser"`
	StsActive  int    `json:"stsActive"`
}

type RespHisUserFpm struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	RoleUser  string `json:"roleUser"`
	StsActive string `json:"stsActive"`
}

type RespHisFp struct {
	NoInv           string `json:"noInv"`
	NoFp            string `json:"noFp"`
	BillingDate     string `json:"billingDate"`
	CompCode        string `json:"compCode"`
	Efaktur         string `json:"efaktur"`
	Customer        string `json:"customer"`
	CreateDateInv   string `json:"createDateInv"`
	CreateDateFp    string `json:"createDateFp"`
	SendDateInv     string `json:"sendDateInv"`
	SendDateFp      string `json:"sendDateFp"`
	StatusInv       string `json:"statusInv"`
	StatusFp        string `json:"statusFp"`
	EmailReceiptInv string `json:"emailReceiptInv"`
	EmailReceiptFp string `json:"emailReceiptFp"`
	NoReference     string `json:"noReference"`
	IdEfaktur       int    `json:"idEfaktur"`
	IdInvoice       int    `json:"idInvoice"`
	StsFaktur       string `json:"stsFaktur"`
	Dc              string `json:"dc"`
	Cmd             string `json:"cmd"`
}
type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}
type RespUploadSourceData struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type RespRoleUser struct {
	Id       int    `json:"id"`
	RoleUser string `json:"roleUser"`
}

type RespDownloadFpm struct {
	BaseUrl  string `json:"baseUrl"`
	FileName string `json:"fileName"`
}
type RespDataFpDto struct {
	ID               int    `json:"id"`
	UploadBy         string `json:"uploadBy"`
	JenisFakturPajak string `json:"jenisFakturPajak"`
	FileNameUpload   string `json:"fileNameUpload"`
	NoFaktur         string `json:"noFaktur"`
	Status           string `json:"status"`
	InvDc            string `json:"invDc"`
	InvNumber        string `json:"invNumber"`
	InvDate          string `json:"invDate"`
	FpCreatedDate    string `json:"fpCreatedDate"`
	CreateAt         string `json:"createAt"`
	UpdateAt         string `json:"updateAt"`
}
