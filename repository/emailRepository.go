package repository

import (
	"backend_gui/dto"
	"backend_gui/models"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmailRepository interface {
	GetDataInvCancel(dateInput string) []map[string]interface{}
	GetDataInvNoFp(dateInput string) []map[string]interface{}
	GetDataInvNotSend(dateInput string) []map[string]interface{}
	GetDataNoFileFp(dateInput string) []map[string]interface{}
	GetDataFpCancel(dateInput string) []map[string]interface{}
	GetDataFpNotSend(dateInput string) []map[string]interface{}
	SaveEmailNotif(data models.EmailNotificationModels) dto.ResSaveEmailDto
	UpdateStatusEmail(idHisEmail int, data int) error
}

type mapEmailRepositoryCon struct {
	emailRepositoryCon *gorm.DB
}

func (m mapEmailRepositoryCon) UpdateStatusEmail(idHisEmail int, idCompare int) error  {
	dtaUpdtHisFile := map[string]interface{}{
		"id_his_email": idHisEmail,
	}

	errUpdtFp := m.emailRepositoryCon.Table("sf_fpm_compare").
		Where("sf_fpm_compare.id=?", idCompare).Updates(&dtaUpdtHisFile)
    if errUpdtFp.Error != nil{
    	return fmt.Errorf("%s","error update id his email : ",  errUpdtFp.Error.Error())
	}
	return nil
}

func (m mapEmailRepositoryCon) GetDataFpNotSend(dateInput string) []map[string]interface{} {
	var results []map[string]interface{}
	query := m.emailRepositoryCon.Table("sf_fpm_compare").Select(
		"sf_fpm_compare.id AS id_compare,"+
		"sf_fpm_compare.billing_document_zsd001n,"+
			"sf_fpm_compare.payer_name_zsd001n,"+
			"sf_fpm_compare.dc_zsd001n,"+
			"sf_fpm_faktur.faktur,"+
			"sf_fpm_compare.billing_date_zsd001n,"+
			"sf_fpm_history_upload_pjk.file_name_upload").
		Joins("inner join sf_fpm_faktur on sf_fpm_compare.billing_document_zsd001n = sf_fpm_faktur.inv_number").
		Joins("inner join sf_fpm_history_upload_pjk on sf_fpm_faktur.faktur = sf_fpm_history_upload_pjk.no_faktur").
		Where("sf_fpm_faktur.status=?", 1).
		Where("sf_fpm_faktur.update_at IS NULL").
		Where("sf_fpm_compare.sts_compare=?", "match").
		Where("sf_fpm_compare.sts_cancel_inv=?", "no").
		Where("sf_fpm_compare.sts_send_inv=?", "yes").
		Where("sf_fpm_compare.dc_zsd001n=?", "10").
		Where("sf_fpm_faktur.status_send_fp=?", "no").
		Where("sf_fpm_history_upload_pjk.status=?", "failed_send_to_cust")
	if dateInput != "" {
		query.Where("DATE(sf_fpm_history_upload_pjk.update_at)=?", dateInput)
	} else {
		query.Where("DATE(sf_fpm_history_upload_pjk.update_at)= CURRENT_DATE()")
	}
	errorSelect := query.Scan(&results)
	if errorSelect.Error != nil{
		log.Error("error get data Fp Not send : ", errorSelect.Error.Error())
	}
	return results
}

func (m mapEmailRepositoryCon) GetDataFpCancel(dateInput string) []map[string]interface{} {
	var results []map[string]interface{}
	query := m.emailRepositoryCon.Table("sf_fpm_faktur").Select(
		"sf_fpm_compare.id AS id_compare,"+
		"sf_fpm_compare.billing_document_zsd001n,"+
			"sf_fpm_compare.payer_name_zsd001n,"+
			"sf_fpm_compare.dc_zsd001n,"+
			"sf_fpm_faktur.faktur,"+
			"sf_fpm_compare.billing_date_zsd001n").
		Joins("inner join sf_fpm_compare on sf_fpm_faktur.inv_number = sf_fpm_compare.billing_document_zsd001n").
		Where("sf_fpm_faktur.status=?", 0).
		Where("sf_fpm_faktur.status_send_fp=?", "yes").
		Where("sf_fpm_compare.sts_compare=?", "match").
		Where("sf_fpm_compare.sts_send_inv=?", "yes").
		Where("sf_fpm_compare.dc_zsd001n=?", "10")
	if dateInput != "" {
		query.Where("DATE(sf_fpm_faktur.update_at)=?", dateInput)
	} else {
		query.Where("DATE(sf_fpm_faktur.update_at)= CURRENT_DATE()")
	}
	errorSelect := query.Scan(&results)
	if errorSelect.Error != nil{
		log.Error("error get data fp cancel : ", errorSelect.Error.Error())
	}
	return results
}

func (m mapEmailRepositoryCon) GetDataNoFileFp(dateInput string) []map[string]interface{} {
	var results []map[string]interface{}
	query := m.emailRepositoryCon.Table("sf_fpm_compare").Select(
		"sf_fpm_compare.id AS id_compare,"+
		"sf_fpm_compare.billing_document_zsd001n,"+
			"sf_fpm_compare.payer_name_zsd001n,"+
			"sf_fpm_compare.dc_zsd001n,"+
			"sf_fpm_faktur.faktur,"+
			"sf_fpm_compare.billing_date_zsd001n").
		Joins("left join sf_fpm_faktur on sf_fpm_compare.billing_document_zsd001n = sf_fpm_faktur.inv_number").
		Joins("left join sf_fpm_history_upload_pjk on sf_fpm_faktur.faktur = sf_fpm_history_upload_pjk.no_faktur").
		Where("sf_fpm_faktur.status=?", 1).
		Where("sf_fpm_compare.sts_compare=?", "match").
		Where("sf_fpm_compare.sts_cancel_inv=?", "no").
		Where("sf_fpm_compare.dc_zsd001n IN (?,?,?,?,?)", "10", "20", "35", "40", "45").
		Where("sf_fpm_history_upload_pjk.id IS NULL")
	if dateInput != "" {
		query.Where("DATE(sf_fpm_faktur.create_at)=?", dateInput)
	} else {
		query.Where("DATE(sf_fpm_faktur.create_at)= CURRENT_DATE()")
	}
	errSelect := query.Scan(&results)
	if errSelect.Error != nil{
		log.Error("error get data no file fp : ", errSelect.Error.Error())
	}
	return results

}

func (m mapEmailRepositoryCon) SaveEmailNotif(data models.EmailNotificationModels) dto.ResSaveEmailDto {
	var dataNotif dto.ResSaveEmailDto
	errSave := m.emailRepositoryCon.Save(&data)
	if errSave.Error != nil {
		log.Error("Failed save data email cancel", errSave.Error.Error())
	}
	dataNotif.Id = data.Id
	dataNotif.Recipient = data.Recipient
	dataNotif.Cc = data.Cc
	dataNotif.Bcc = data.Bcc
	dataNotif.Sender = data.Sender
	dataNotif.Subject = data.Subject
	dataNotif.PathAttachment = data.PathAttachment
	dataNotif.LogStatus = data.LogStatus
	dataNotif.Body = data.Body
	dataNotif.SendDateTime = data.SendDateTime
	return dataNotif
}
func (m mapEmailRepositoryCon) GetDataInvNoFp(dateInput string) []map[string]interface{} {
	var results []map[string]interface{}
	query := m.emailRepositoryCon.Table("sf_fpm_compare").Select(
		"sf_fpm_compare.id AS id_compare,"+
		"sf_fpm_compare.billing_document_zsd001n,"+
			"sf_fpm_compare.billing_date_zsd001n,"+
			"sf_fpm_compare.payer_name_zsd001n,"+
			"sf_fpm_compare.dc_zsd001n").
		Where("sf_fpm_compare.sts_cancel_inv=?", "no").
		Where("sf_fpm_compare.sts_compare=?", "not_match")
	if dateInput != "" {
		query.Where("sf_fpm_compare.create_on_zsd001n = ?", dateInput)
	} else {
		query.Where("STR_TO_DATE(sf_fpm_compare.create_on_zsd001n,'%Y%m%d') = CURRENT_DATE() - INTERVAL 1 DAY")
	}
	errSelect := query.Scan(&results)
	if errSelect.Error != nil{
		log.Error("error get email dana inv no fp : ", errSelect.Error.Error())
	}
	return results

}
func (m mapEmailRepositoryCon) GetDataInvCancel(dateInput string) []map[string]interface{} {
	var results []map[string]interface{}
	query := m.emailRepositoryCon.Table("sf_fpm_compare").
		Select(
			"sf_fpm_compare.id as id_inv_cancel,"+
				"sf_fpm_compare.billing_document_zsd001n,"+
				"sf_fpm_compare.billing_doc_cancel,"+
				"sf_fpm_compare.fp_number_zv60,"+
				"sf_fpm_compare.billing_date_zsd001n,"+
				"sf_fpm_compare.payer_name_zsd001n,"+
				"sf_fpm_compare.dc_zsd001n").
		Where("sf_fpm_compare.sts_cancel_inv=?", "yes").
		Where("sf_fpm_compare.sts_send_inv=?", "yes").
		Where("sf_fpm_compare.dc_zsd001n=?", "10").
		Where("sf_fpm_compare.sts_email_inv_cancel=?", "ready_to_send")
	if dateInput != "" {
		query.Where("DATE(sf_fpm_compare.dtm_updated)=?", dateInput)
	} else {
		query.Where("DATE(sf_fpm_compare.dtm_updated)= CURRENT_DATE()")
	}
	errSelect :=query.Scan(&results)
	if errSelect.Error != nil{
		log.Error("error get inv cancel : ", errSelect.Error.Error())
	}
	return results
}
func (m mapEmailRepositoryCon) GetDataInvNotSend(dateInput string) []map[string]interface{} {
	var results []map[string]interface{}
	query := m.emailRepositoryCon.Table("sf_fpm_compare").Select(
		"sf_fpm_compare.id AS id_compare,"+
			"sf_fpm_compare.billing_document_zsd001n,"+
			"sf_fpm_compare.payer_name_zsd001n,"+
			"sf_fpm_compare.dc_zsd001n,"+
			"sf_fpm_faktur.faktur,"+
			"sf_fpm_compare.billing_date_zsd001n").
		Joins("inner join sf_fpm_faktur on sf_fpm_compare.billing_document_zsd001n = sf_fpm_faktur.inv_number").
		Joins("inner join sf_fpm_history_upload_pjk on sf_fpm_faktur.faktur = sf_fpm_history_upload_pjk.no_faktur").
		Where("sf_fpm_faktur.status=?", 1).
		Where("sf_fpm_faktur.update_at IS NULL").
		Where("sf_fpm_compare.sts_compare=?", "match").
		Where("sf_fpm_compare.sts_cancel_inv=?", "no").
		Where("sf_fpm_compare.sts_send_inv=?", "no").
		Where("sf_fpm_history_upload_pjk.status=?", "uploaded").
		Where("sf_fpm_compare.dc_zsd001n=?", "10")
	if dateInput != "" {
		query.Where("DATE(sf_fpm_faktur.create_at)=?", dateInput)
	} else {
		query.Where("DATE(sf_fpm_faktur.create_at)= CURRENT_DATE()")
	}
	errSelect :=query.Scan(&results)
	if errSelect.Error != nil{
		log.Error("error get data inv not send : ", errSelect.Error.Error())
	}
	return results
}

func InstanceEmailRepository(db *gorm.DB) EmailRepository {
	return &mapEmailRepositoryCon{
		emailRepositoryCon: db,
	}
}
