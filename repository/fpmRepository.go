package repository

import (
	"backend_gui/dto"
	"backend_gui/models"
	"fmt"
	"github.com/mashingan/smapping"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FpmRepository interface {
	GetHistoryFpmByName(fileName string, jenisFakturPjk string) dto.RespFpmHistoryDto
	CheckNoFpmZV60(noFaktur string) dto.DtoFpmZv60
	UpdateStsFpmHistory(data models.FpmHistoryFilePjk) error
	UpdateDataFpm(dtaFp DataFp) error
	GetHistoryFpmBySts(stsMvSap string, stsFailed string) []DataEfakturBySts
	GetDataFp() []DataFp
}

type mapFpmRepositoryCon struct {
	mapFpmRepositoryCon *gorm.DB
}

func (m mapFpmRepositoryCon) UpdateDataFpm(dtaFp DataFp) error {

	tx := m.mapFpmRepositoryCon.Begin()
	dtaUpdtCompare := models.DetailCompareModels{
		Id:                     dtaFp.IdCompare,
		BillingDocumentZsd001n: dtaFp.BillingDocumentZsd001n,
		BillingTypeZsd001n:     dtaFp.BillingTypeZsd001n,
		BillingDateZsd001n:     dtaFp.BillingDateZsd001n,
		DcZsd001n:              dtaFp.DcZsd001n,
		PayerNameZsd001n:       dtaFp.PayerNameZsd001n,
		CreateOnZsd001n:        dtaFp.CreateOnZsd001n,
		BillingDocCancel:       dtaFp.BillingDocCancel,
		StsCancelInv:           dtaFp.StsCancelInv,
		StsEmailInvCancel:      dtaFp.StsEmailInvCancel,
		StsSendInv:             dtaFp.StsSendInv,
		BillingNumberZv60:      dtaFp.BillingNumberZv60,
		BillingDateZv60:        dtaFp.BillingDateZv60,
		FpNumberZv60:           dtaFp.FpNumberZv60,
		FpCreatedDateZv60:      dtaFp.FpCreatedDateZv60,
		StsCompare:             dtaFp.StsCompare,
		IdHistoryFpm:           dtaFp.IdFileFaktur,
		StsEmailCompare:        dtaFp.StsEmailCompare,
		EmailReceiptInv:        dtaFp.EmailReceiptInv,
		SendDateInv:            dtaFp.SendDateInv,
		SendTimeInv:            dtaFp.SendTimeInv,
		IdHisEmail:         dtaFp.IdHisEmail,
	}
	log.Info("Update data FP for move file  : ", dtaUpdtCompare)
	errUpdtCompare := tx.Save(&dtaUpdtCompare)
	if errUpdtCompare.Error != nil {
		log.Error("Error update history fpm data ", dtaUpdtCompare, "with error : ", errUpdtCompare.Error.Error())
		tx.Rollback()
		return fmt.Errorf("error update status fpm")
	}
	dtaUpdtHisFile := map[string]interface{}{
		"status": dtaFp.StsFile,
		"url":    dtaFp.UrlFile,
	}
	log.Info("Update data history file FP for move file  : ", dtaUpdtHisFile)
	errUpdtFp := tx.Table("sf_fpm_history_upload_pjk").
		Where("sf_fpm_history_upload_pjk.id=?", dtaFp.IdFileFaktur).Updates(&dtaUpdtHisFile)
	if errUpdtFp.Error != nil {
		log.Error("rror update history fpm data : ", errUpdtFp.Error.Error())
		tx.Rollback()
		return fmt.Errorf("error update status fpm")
	}
	tx.Commit()
	return nil
}

func (m mapFpmRepositoryCon) GetDataFp() []DataFp {
	var records []DataFp
	errSelect :=m.mapFpmRepositoryCon.Table("sf_fpm_compare").Select(
		"sf_fpm_compare.id AS id_compare,"+
			" sf_fpm_compare.billing_document_zsd001n,"+
			" sf_fpm_compare.billing_type_zsd001n,"+
			" sf_fpm_compare.billing_date_zsd001n,"+
			" sf_fpm_compare.dc_zsd001n,"+
			" sf_fpm_compare.payer_zsd001n,"+
			" sf_fpm_compare.payer_name_zsd001n,"+
			" sf_fpm_compare.slor_zsd001n,"+
			" sf_fpm_compare.district_zsd001n,"+
			" sf_fpm_compare.material_number_zsd001n,"+
			" sf_fpm_compare.material_desc_zsd001n,"+
			" sf_fpm_compare.create_on_zsd001n,"+
			" sf_fpm_compare.billing_doc_cancel,"+
			" sf_fpm_compare.sts_cancel_inv,"+
			" sf_fpm_compare.sts_email_inv_cancel,"+
			" sf_fpm_compare.sts_send_inv,"+
			" sf_fpm_compare.billing_number_zv60,"+
			" sf_fpm_compare.billing_date_zv60,"+
			" sf_fpm_compare.fp_number_zv60,"+
			" sf_fpm_compare.fp_created_date_zv60,"+
			" sf_fpm_compare.payer_zv60,"+
			" sf_fpm_compare.name_zv60,"+
			" sf_fpm_compare.npwp_zv60,"+
			" sf_fpm_compare.material_zv60,"+
			" sf_fpm_compare.sts_compare,"+
			" sf_fpm_compare.id_history_fpm,"+
			" sf_fpm_compare.sts_email_compare,"+
			" sf_fpm_compare.email_receipt_inv,"+
			" sf_fpm_compare.send_date_inv,"+
			" sf_fpm_compare.send_time_inv,"+
			" sf_fpm_compare.id_his_email,"+
			" sf_fpm_faktur.id AS id_history_fp,"+
			" sf_fpm_faktur.faktur,"+
			" sf_fpm_history_upload_pjk.id AS id_file_faktur,"+
			" sf_fpm_history_upload_pjk.jenis_faktur_pajak,"+
			" sf_fpm_history_upload_pjk.status AS sts_file,"+
			" sf_fpm_history_upload_pjk.file_name_upload AS file_name_fp,"+
			" sf_fpm_history_upload_pjk.url as url_file").
		Joins("inner join sf_fpm_faktur on sf_fpm_compare.billing_document_zsd001n = sf_fpm_faktur.inv_number").
		Joins("inner join sf_fpm_history_upload_pjk on sf_fpm_faktur.faktur = sf_fpm_history_upload_pjk.no_faktur").
		Where("sf_fpm_compare.sts_compare = ?", "match").
		Where("sf_fpm_compare.sts_cancel_inv = ?", "no").
		Where("sf_fpm_faktur.status = ?", 1).
		Where("sf_fpm_faktur.update_at IS NULL").
		Where("(sf_fpm_compare.dc_zsd001n != '10' OR (sf_fpm_compare.dc_zsd001n = '10' AND sf_fpm_compare.sts_send_inv = 'yes'))").
		Where("sf_fpm_history_upload_pjk.status", "uploaded").Scan(&records)
	if errSelect.Error != nil{
		log.Error("Error get data fp : ", errSelect.Error.Error())
	}
	return records
}

func (m mapFpmRepositoryCon) GetHistoryFpmBySts(stsMvSap string, stsFailed string) []DataEfakturBySts {
	var dtoFpm []DataEfakturBySts
	errSelect := m.mapFpmRepositoryCon.Table("sf_fpm_compare").Select(
		"sf_fpm_history_upload_pjk.id,"+
			"sf_fpm_history_upload_pjk.jenis_faktur_pajak,"+
			"sf_fpm_faktur.faktur AS no_faktur,"+
			"sf_fpm_history_upload_pjk.status,"+
			"sf_fpm_history_upload_pjk.file_name_upload,"+
			"sf_fpm_history_upload_pjk.update_at,"+
			"sf_fpm_history_upload_pjk.upload_by,"+
			"sf_fpm_history_upload_pjk.url,"+
			"sf_fpm_faktur.status_send_fp").
		Joins("inner join sf_fpm_faktur on sf_fpm_compare.billing_document_zsd001n = sf_fpm_faktur.inv_number").
		Joins("inner join sf_fpm_history_upload_pjk on sf_fpm_faktur.faktur = sf_fpm_history_upload_pjk.no_faktur").
		Where("sf_fpm_compare.sts_compare=?", "match").
		Where("sf_fpm_compare.sts_cancel_inv=?", "no").
		Where("sf_fpm_faktur.status=?", 1).
		Where("sf_fpm_faktur.update_at IS NULL").
		Where("sf_fpm_compare.sts_send_inv=?", "yes").
		Where("sf_fpm_history_upload_pjk.status IN ('"+stsMvSap+"'"+","+"'"+stsFailed+"'"+")").
		Where("sf_fpm_compare.dc_zsd001n=?", "10").Scan(&dtoFpm)
	if errSelect.Error != nil{
		log.Error("Error get his fpm by history : ", errSelect.Error.Error())
	}
	return dtoFpm
}

func (m mapFpmRepositoryCon) UpdateStsFpmHistory(data models.FpmHistoryFilePjk) error {
	errUpdate := m.mapFpmRepositoryCon.Save(&data)
	if errUpdate.Error != nil {
		log.Error("Error update history fpm data ", data, "with error : "+errUpdate.Error.Error())
		return fmt.Errorf("error update status fpm")
	}
	return nil

}

func (m mapFpmRepositoryCon) CheckNoFpmZV60(noFaktur string) dto.DtoFpmZv60 {
	var data dto.DtoFpmZv60
	dataDetailZv60 := models.FpmDetailZv60{}
	m.mapFpmRepositoryCon.Where("fp_number=?", noFaktur).First(&dataDetailZv60)
	err := smapping.FillStructByTags(&data, smapping.MapTags(&dataDetailZv60, "json"), "json")
	if err != nil {
		log.Error("Error mapping value ", err.Error())
	}
	return data

}

func (m mapFpmRepositoryCon) GetHistoryFpmByName(fileName string, jenisFakturPjk string) dto.RespFpmHistoryDto {
	var dataHistoryFpm dto.RespFpmHistoryDto
	historyFpm := models.FpmHistoryFilePjk{}
	errSelect :=m.mapFpmRepositoryCon.Where("file_name_upload=? and jenis_faktur_pajak=?", fileName, jenisFakturPjk).First(&historyFpm)
	if errSelect.Error != nil{
		log.Error("Error get his file fpm : ", errSelect.Error.Error())
	}
	dataHistoryFpm.ID = historyFpm.ID
	dataHistoryFpm.UploadBy = historyFpm.UploadBy
	dataHistoryFpm.JenisFakturPajak = historyFpm.JenisFakturPajak
	dataHistoryFpm.FileNameUpload = historyFpm.FileNameUpload
	dataHistoryFpm.NoFaktur = historyFpm.NoFaktur
	dataHistoryFpm.Url = historyFpm.Url
	dataHistoryFpm.Status = historyFpm.Status
	return dataHistoryFpm
}

func InstanceFpmRepository(db *gorm.DB) FpmRepository {
	return &mapFpmRepositoryCon{
		mapFpmRepositoryCon: db,
	}
}

type DataFp struct {
	IdCompare              int    `json:"IdCompare"`
	BillingDocumentZsd001n string `json:"billingDocumentZsd001n"`
	BillingTypeZsd001n     string `json:"billingTypeZsd001n"`
	BillingDateZsd001n     string `json:"billingDateZsd001n"`
	DcZsd001n              string `json:"dcZsd001n"`
	PayerNameZsd001n       string `json:"payerNameZsd001n"`
	CreateOnZsd001n        string `json:"createOnZsd001n"`
	BillingDocCancel       string `json:"billingDocCancel"`
	StsCancelInv           string `json:"stsCancelInv"`
	StsEmailInvCancel      string `json:"stsEmailInvCancel"`
	StsSendInv             string `json:"stsSendInv"`
	BillingNumberZv60      string `json:"billingNumberZv60"`
	BillingDateZv60        string `json:"billingDateZv60"`
	FpNumberZv60           string `json:"fpNumberZv60"`
	FpCreatedDateZv60      string `json:"fpCreatedDateZv60"`
	StsCompare             string `json:"stsCompare"`
	IdHistoryFpm           int    `json:"idHistoryFpm"`
	StsEmailCompare        string `json:"stsEmailCompare"`
	EmailReceiptInv        string `json:"emailReceiptInv"`
	SendDateInv            string `json:"sendDateInv"`
	SendTimeInv            string `json:"sendTimeInv"`
	IdHistoryFp            int    `json:"idHistoryFp"`
	Faktur                 string `json:"faktur"`
	StatusSendFp           string `json:"statusSendFp"`
	Status                 int    `json:"status"`
	IdFileFaktur           int    `json:"idFileFaktur"`
	JenisFakturPajak       string `json:"jenisFakturPajak"`
	FileNameFp             string `json:"fileNameFp"`
	StsFile                string `json:"stsFile"`
	UrlFile                string `json:"urlFile"`
	IdHisEmail             int `json:"idHisEmail"`
}

type DataEfakturBySts struct {
	Id               int    `json:"Id"`
	JenisFakturPajak string `json:"jenisFakturPajak"`
	NoFaktur         string `json:"noFaktur"`
	Status           string `json:"status"`
	FileNameUpload   string `json:"fileNameUpload"`
	UpdateAt         string `json:"updateAt"`
	UploadBy         string `json:"uploadBy"`
	Url              string `json:"Url"`
	StatusSendFp     string `json:"statusSendFp"`
}
