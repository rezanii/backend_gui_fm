package repository

import (
	"backend_gui/models"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type TestRepository interface {
	SelectTableRole(tableName string) error
	SelectTableUser(tableName string) error
	SelectTableFaktur(tableName string) error
	SelectTableHisUploadPjk(tableName string) error
	SelectTableEmailFpm(tableName string) error
	SelectTableCompareFpm(tableName string) error
	SelectTableZsd081(tableName string) error
}

type mapTestRepositoryCon struct {
	testRepository *gorm.DB
}

func (m mapTestRepositoryCon) SelectTableZsd081(tableName string) error {
	var dataZsd081 models.FpmDetailInputZsd081
	errSelect := m.testRepository.Model(models.FpmDetailInputZsd081{}).First(&dataZsd081)
	if errSelect.Error != nil && errors.Is(errSelect.Error, gorm.ErrRecordNotFound) == false{
		log.Error("error test select db, with table "+tableName+" : ", errSelect.Error.Error())
		return fmt.Errorf("%s", "Somenthing wrong table "+tableName)
	}
	log.Info("data tes zsd081 : ", dataZsd081)
	return nil
}

func (m mapTestRepositoryCon) SelectTableCompareFpm(tableName string) error {
	var dataCompareFpm CompareFpm
	errSelect := m.testRepository.Table(tableName).Select(
		tableName+".id,"+
			tableName+".billing_document_zsd001n,"+
			tableName+".billing_type_zsd001n,"+
			tableName+".billing_date_zsd001n,"+
			tableName+".dc_zsd001n,"+
			tableName+".payer_zsd001n,"+
			tableName+".payer_name_zsd001n,"+
			tableName+".slor_zsd001n,"+
			tableName+".district_zsd001n,"+
			tableName+".material_number_zsd001n,"+
			tableName+".material_desc_zsd001n,"+
			tableName+".create_on_zsd001n,"+
			tableName+".billing_doc_cancel,"+
			tableName+".sts_cancel_inv,"+
			tableName+".sts_email_inv_cancel,"+
			tableName+".sts_send_inv,"+
			tableName+".billing_number_zv60,"+
			tableName+".fp_number_zv60,"+
			tableName+".billing_date_zv60,"+
			tableName+".fp_created_date_zv60,"+
			tableName+".payer_zv60,"+
			tableName+".name_zv60,"+
			tableName+".npwp_zv60,"+
			tableName+".material_zv60,"+
			tableName+".sts_compare,"+
			tableName+".sts_email_compare,"+
			tableName+".email_receipt_inv,"+
			tableName+".send_date_inv,"+
			tableName+".send_time_inv,"+
			tableName+".id_history_fpm,"+
			tableName+".id_his_email,"+
			tableName+".dtm_created,"+
			tableName+".dtm_updated").First(&dataCompareFpm)
	if errSelect.Error != nil && errors.Is(errSelect.Error, gorm.ErrRecordNotFound) == false{
		log.Error("error test select db, with table "+tableName+" : ", errSelect.Error.Error())
		return fmt.Errorf("%s", "Somenthing wrong table "+tableName)
	}
	log.Info("data tes his compare fpm : ", dataCompareFpm)
	return nil

}

func (m mapTestRepositoryCon) SelectTableEmailFpm(tableName string) error {
	var dataEmailFpm EmailFpm
	errSelect := m.testRepository.Table(tableName).Select(
		tableName+".id",
		tableName+".recipient,"+
			tableName+".cc,"+
			tableName+".bcc,"+
			tableName+".sender,"+
			tableName+".subject,"+
			tableName+".path_attachment,"+
			tableName+".log_status,"+
			tableName+".body,"+
			tableName+".send_date_time,"+
			tableName+".create_date").First(&dataEmailFpm)
	if errSelect.Error != nil && errors.Is(errSelect.Error, gorm.ErrRecordNotFound) == false {
		log.Error("error test select db, with table "+tableName+" : ", errSelect.Error.Error())
		return fmt.Errorf("%s", "Somenthing wrong table "+tableName)
	}
	log.Info("data tes his email fpm : ", dataEmailFpm)
	return nil
}

func (m mapTestRepositoryCon) SelectTableHisUploadPjk(tableName string) error {
	var dataHisUploadPjk HisUploadPjk
	errSelect := m.testRepository.Table(tableName).Select(tableName + ".id," +
		tableName + ".upload_by," +
		tableName + ".jenis_faktur_pajak," +
		tableName + ".file_name_upload," +
		tableName + ".no_faktur," +
		tableName + ".url," +
		tableName + ".status," +
		tableName + ".create_at," +
		tableName + ".update_at").First(&dataHisUploadPjk)
	if errSelect.Error != nil && errors.Is(errSelect.Error, gorm.ErrRecordNotFound) == false {
		log.Error("error test select db, with table "+tableName+" : ", errSelect.Error.Error())
		return fmt.Errorf("%s", "Somenthing wrong table "+tableName)
	}
	log.Info("data tes his upload pjk : ", dataHisUploadPjk)
	return nil
}

func (m mapTestRepositoryCon) SelectTableFaktur(tableName string) error {
	var dataFaktur faktur
	errSelect := m.testRepository.Table(tableName).Select(tableName + ".id," +
		tableName + ".inv_number," +
		tableName + ".inv_dc," +
		tableName + ".faktur," +
		tableName + ".faktur_created," +
		tableName + ".status_send_fp," +
		tableName + ".status," +
		tableName + ".email_receipt_fp," +
		tableName + ".send_date_fp," +
		tableName + ".send_time_fp," +
		tableName + ".create_at," +
		tableName + ".update_at").First(&dataFaktur)
	if errSelect.Error != nil && errors.Is(errSelect.Error, gorm.ErrRecordNotFound) == false{
		log.Error("error test select db, with table "+tableName+" : ", errSelect.Error.Error())
		return fmt.Errorf("%s", "Somenthing wrong table "+tableName)
	}
	log.Info("data tes faktur : ", dataFaktur)
	return nil
}

func (m mapTestRepositoryCon) SelectTableUser(tableName string) error {
	var dataUserFpm UserFpm
	errSelect := m.testRepository.Table(tableName).Select(tableName + ".id," +
		tableName + ".username," +
		tableName + ".id_role_user," +
		tableName + ".sts_active," +
		tableName + ".dtm_create," +
		tableName + ".dtm_updated").First(&dataUserFpm)
	log.Info("data tes user : ", dataUserFpm)
	if errSelect.Error != nil && errors.Is(errSelect.Error, gorm.ErrRecordNotFound) == false{
		log.Error("error test select db, with table "+tableName+" : ", errSelect.Error.Error())
		return fmt.Errorf("%s", "Somenthing wrong table "+tableName)
	}
	return nil
}

func (m mapTestRepositoryCon) SelectTableRole(tableName string) error {
	var dataRole RoleFpm
	errSelect := m.testRepository.Table(tableName).Select(tableName + ".id," +
		tableName + ".role_user," +
		tableName + ".delete_at," +
		tableName + ".dtm_created," +
		tableName + ".dtm_updated").First(&dataRole)
	log.Info("data tes role : ", dataRole)
	if errSelect.Error != nil && errors.Is(errSelect.Error, gorm.ErrRecordNotFound) == false{
		log.Error("error test select db, with table "+tableName+": ", errSelect.Error.Error())
		return fmt.Errorf("%s", "Somenthing wrong table "+tableName)
	}


	return nil
}

type RoleFpm struct {
	Id         int
	RoleUser   string
	DeleteAt   int
	DtmCreated time.Time
	DtmUpdated time.Time
}

type UserFpm struct {
	Id         int
	Username   string
	IdRoleUser string
	StsActive  string
	DtmCreate  time.Time
	DtmUpdated time.Time
}

type faktur struct {
	Id             int
	InvNumber      string
	InvDc          string
	Faktur         string
	FakturCreated  string
	StatusSendFp   string
	Status         int
	EmailReceiptFp string
	SendDatFp      string
	SendTimeFp     string
	CreateAt       time.Time
	UpdateAt       time.Time
}

type HisUploadPjk struct {
	Id               int
	UploadBy         string
	JenisFakturPajak string
	FileNameUpload   string
	NoFaktur         string
	Url              string
	Status           string
	CreateAt         time.Time
	UpdateAt         time.Time
}

type EmailFpm struct {
	Id             int
	Recipient      string
	Cc             string
	Bcc            string
	Sender         string
	Subject        string
	PathAttachment string
	LogStatus      int
	Body           string
	SendDateTime   time.Time
	CreateDate     time.Time
}

type CompareFpm struct {
	Id                     int
	BillingDocumentZsd001n string
	BillingTypeZsd001n     string
	BillingDateZsd001n     string
	DcZsd001n              string
	PayerZsd001n           string
	PayerNameZsd001n       string
	SlorZsd001n            string
	DistrictZsd001n        string
	MaterialNumberZsd001n  string
	MaterialDescZsd001n    string
	CreateOnZsd001n        string
	BillingDocCancel       string
	StsCancelInv           string
	StsEmailInvCancel      string
	StsSendInv             string
	BillingNumberZv60      string
	FpNumberZv60           string
	BillingDateZv60        string
	FpCreatedDateZv60      string
	PayerZv60              string
	NameZv60               string
	NpwpZv60               string
	MaterialZv60           string
	StsCompare             string
	StsEmailCompare        string
	EmailReceiptInv        string
	SendDateInv            string
	SendTimeInv            string
	IdHistoryFpm           int
	IdHisEmail             int
	DtmCreated             time.Time
	DtmUpdated             time.Time
}

func InstanceTestRepository(db *gorm.DB) TestRepository {
	return &mapTestRepositoryCon{
		testRepository: db,
	}
}
