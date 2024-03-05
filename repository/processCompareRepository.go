package repository

import (
	"backend_gui/dto"
	"backend_gui/models"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ProcessCompareRepository interface {
	CompareZsd001nZv60(dateInput string) []DetailCompare
	UpdateCompareZsd001nZv60(data models.DetailCompareModels) error
	CompareInvoiceCancel(date string) []DetailCompare
	UpdateInvoiceCancel(dataInvoiceCancel models.DetailCompareModels) error
	GetFpByCreated(dateInput string) []GetFpCancel
	UpdateFpAlreadySend(dataFp models.FakturModels) error
	GetDataZsd001nByCreateOn(date string) []Zsd001n
	CreateDataCompare([]models.DetailCompareModels) error
	GetHistoryFaktur(faktur string) FakturHistory
	UpdateHistoryFaktur(billNumber string) error
	SaveHistoryFaktur(data models.FakturModels) error
	GetHistoryFpCancel() []FakturHistory
	GetCompareInvZsd081(dateInput string) []DetailCompare
	UpdateInvAlreadySend(data models.DetailCompareModels) error
	//GetCompareFpHisZsd081(dateInput string) []CompareZsd081Faktur
	GetCompareFpHisZsd081() []CompareZsd081Faktur
	GetCompareByInvoice(invoice string) DetailCompare
	GetTotalCollectionInv(req dto.ReqCollectionInvDto, dc []string) int64
	GetRecordCollectionInv(start int, end int, req dto.ReqCollectionInvDto, dc []string) []dto.RespHisFp
	CheckUserAlreadyExist(username string) UserWithRole
}
type mapProcessCompareRepositoryCon struct {
	mapProcessCompareRepositoryCon *gorm.DB
}

func (m mapProcessCompareRepositoryCon) CheckUserAlreadyExist(username string) UserWithRole {
	var records UserWithRole
	m.mapProcessCompareRepositoryCon.Table("sf_fpm_user").Select("sf_fpm_user.username, sf_fpm_user.id,"+
		"sf_fpm_role.role_user").
		Joins("inner join sf_fpm_role on sf_fpm_role.id = sf_fpm_user.id_role_user").
		Where("sf_fpm_user.username=?", username).
		Where("sf_fpm_user.sts_active=?", "1").First(&records)
	return records
}

func (m mapProcessCompareRepositoryCon) GetRecordCollectionInv(start int, end int, req dto.ReqCollectionInvDto, dc []string) []dto.RespHisFp {
	var data []dto.RespHisFp
	query := m.mapProcessCompareRepositoryCon.Table("sf_fpm_compare")
	if req.TypeDc == "cluster" {
		query.Select(
			"IF (sf_fpm_compare.sts_cancel_inv = 'yes', sf_fpm_compare.billing_doc_cancel, sf_fpm_compare.billing_document_zsd001n )AS no_inv," +
				"IF(sf_fpm_faktur.faktur IS NULL,'', sf_fpm_faktur.faktur) AS no_fp," +
				"DATE_FORMAT(STR_TO_DATE(sf_fpm_compare.billing_date_zsd001n, '%Y%m%d'),'%Y-%m-%d') AS billing_date," +
				"sf_dump_zsd001n.slor AS comp_code," +
				"IF(sf_fpm_history_upload_pjk.file_name_upload IS NULL,'',sf_fpm_history_upload_pjk.file_name_upload) AS efaktur," +
				"sf_fpm_compare.payer_name_zsd001n AS customer," +
				"DATE_FORMAT(STR_TO_DATE(sf_fpm_compare.create_on_zsd001n, '%Y%m%d'),'%Y-%m-%d') AS create_date_inv," +
				"IF((sf_fpm_compare.fp_created_date_zv60 = '' OR sf_fpm_compare.fp_created_date_zv60 IS NULL),'', DATE_FORMAT(STR_TO_DATE(sf_fpm_compare.fp_created_date_zv60, '%d.%m.%Y'), '%Y-%m-%d')) AS create_date_fp," +
				"IF((sf_fpm_compare.send_date_inv = '' OR sf_fpm_compare.send_date_inv IS NULL),'', CONCAT(sf_fpm_compare.send_date_inv,' ',sf_fpm_compare.send_time_inv)) AS send_date_inv," +
				"IF((sf_fpm_faktur.send_date_fp = '' OR sf_fpm_faktur.send_date_fp IS NULL),'', CONCAT(sf_fpm_faktur.send_date_fp,' ',sf_fpm_faktur.send_time_fp)) AS send_date_fp," +
				"IF(sf_fpm_compare.sts_send_inv = 'yes', 'Berhasil Terkirim', '-') AS status_inv," +
				"IF(sf_fpm_faktur.status_send_fp = 'yes', 'Berhasil Terkirim', '-') AS status_fp," +
				"IF((sf_fpm_compare.email_receipt_inv = '' OR sf_fpm_compare.email_receipt_inv IS NULL),'', sf_fpm_compare.email_receipt_inv) AS email_receipt_inv,"+
				"IF((sf_fpm_faktur.email_receipt_fp = '' OR sf_fpm_faktur.email_receipt_fp IS NULL),'', sf_fpm_faktur.email_receipt_fp) AS email_receipt_fp,"+
				"IF(sf_fpm_compare.sts_cancel_inv = 'yes', sf_fpm_compare.billing_document_zsd001n, '') AS no_reference," +
				"sf_fpm_history_upload_pjk.id AS id_efaktur," +
				"sf_fpm_compare.id AS id_invoice," +
				"sf_fpm_history_upload_pjk.jenis_faktur_pajak AS cmd," +
				"sf_fpm_compare.dc_zsd001n AS dc," +
				"sf_fpm_history_upload_pjk.status AS sts_faktur")
	} else {
		query.Select(
			"IF (sf_fpm_compare.sts_cancel_inv = 'yes', sf_fpm_compare.billing_doc_cancel, sf_fpm_compare.billing_document_zsd001n )AS no_inv," +
				"IF(sf_fpm_faktur.faktur IS NULL,'', sf_fpm_faktur.faktur) AS no_fp," +
				"DATE_FORMAT(STR_TO_DATE(sf_fpm_compare.billing_date_zsd001n, '%Y%m%d'),'%Y-%m-%d') AS billing_date," +
				"sf_dump_zsd001n.slor AS comp_code," +
				"IF(sf_fpm_history_upload_pjk.file_name_upload IS NULL,'',sf_fpm_history_upload_pjk.file_name_upload) AS efaktur," +
				"sf_fpm_compare.payer_name_zsd001n AS customer," +
				"DATE_FORMAT(STR_TO_DATE(sf_fpm_compare.create_on_zsd001n, '%Y%m%d'),'%Y-%m-%d') AS create_date_inv," +
				"IF((sf_fpm_compare.fp_created_date_zv60 = '' OR sf_fpm_compare.fp_created_date_zv60 IS NULL),'', DATE_FORMAT(STR_TO_DATE(sf_fpm_compare.fp_created_date_zv60, '%d.%m.%Y'), '%Y-%m-%d')) AS create_date_fp," +
				"IF(sf_fpm_compare.sts_cancel_inv = 'yes', sf_fpm_compare.billing_document_zsd001n, '') AS no_reference," +
				"sf_fpm_history_upload_pjk.id AS id_efaktur," +
				"sf_fpm_compare.id AS id_invoice," +
				"sf_fpm_history_upload_pjk.jenis_faktur_pajak AS cmd," +
				"sf_fpm_compare.dc_zsd001n AS dc," +
				"sf_fpm_history_upload_pjk.status AS sts_faktur")
	}

	query.Joins("inner join sf_dump_zsd001n on sf_dump_zsd001n.billing_document = sf_fpm_compare.billing_document_zsd001n").
		Joins("left join sf_fpm_faktur on sf_fpm_compare.billing_document_zsd001n = sf_fpm_faktur.inv_number").
		Joins("left join sf_fpm_history_upload_pjk on sf_fpm_compare.fp_number_zv60 = sf_fpm_history_upload_pjk.no_faktur")

	if len(dc) > 0 {
		query.Where("sf_fpm_compare.dc_zsd001n IN (?)", dc)
	}

	query.Where("STR_TO_DATE(sf_fpm_compare.billing_date_zsd001n,'%Y%m%d') BETWEEN STR_TO_DATE('" + req.StartBillDate + "','%Y-%m-%d') AND STR_TO_DATE('" + req.EndBillDate + "','%Y-%m-%d')")

	if req.FilterCustomer != "" && req.FilterInvoice == "" {
		query.Where("sf_fpm_compare.payer_name_zsd001n LIKE ?", "%"+req.FilterCustomer+"%")
	} else if req.FilterCustomer == "" && req.FilterInvoice != "" {
		if string(req.FilterInvoice[0]) == "8" {
			query.Where("sf_fpm_compare.billing_doc_cancel LIKE ?", "%"+req.FilterInvoice+"%")
		} else {
			query.Where("sf_fpm_compare.billing_document_zsd001n LIKE ?", "%"+req.FilterInvoice+"%")
		}

	} else if req.FilterCustomer != "" && req.FilterInvoice != "" {
		if string(req.FilterInvoice[0]) == "8" {
			query.Where("(sf_fpm_compare.payer_name_zsd001n LIKE ? OR sf_fpm_compare.billing_doc_cancel LIKE ?)",
				"%"+req.FilterCustomer+"%", "%"+req.FilterInvoice+"%")
		} else {
			query.Where("(sf_fpm_compare.payer_name_zsd001n LIKE ? OR sf_fpm_compare.billing_document_zsd001n LIKE ?)",
				"%"+req.FilterCustomer+"%", "%"+req.FilterInvoice+"%")
		}

	}
	query.Group("sf_fpm_compare.billing_document_zsd001n, sf_fpm_faktur.faktur")

	err := query.Limit(end).Offset(start).Scan(&data)
	if err.Error != nil {
		log.Error("error get record collection: ", err.Error.Error())
	}
	return data
}

func (m mapProcessCompareRepositoryCon) GetTotalCollectionInv(req dto.ReqCollectionInvDto, dc []string) int64 {
	var totalData int64
	query := m.mapProcessCompareRepositoryCon.Table("sf_fpm_compare").Select("sf_fpm_compare.id").
		Joins("inner join sf_dump_zsd001n on sf_dump_zsd001n.billing_document = sf_fpm_compare.billing_document_zsd001n").
		Joins("left join sf_fpm_history_upload_pjk on sf_fpm_compare.fp_number_zv60 = sf_fpm_history_upload_pjk.no_faktur").
		Joins("left join sf_fpm_faktur on sf_fpm_compare.billing_document_zsd001n = sf_fpm_faktur.inv_number")
	if len(dc) > 0 {
		query.Where("sf_fpm_compare.dc_zsd001n IN (?)", dc)
	}
	query.Where("STR_TO_DATE(sf_fpm_compare.billing_date_zsd001n,'%Y%m%d') BETWEEN STR_TO_DATE('" + req.StartBillDate + "','%Y-%m-%d') AND STR_TO_DATE('" + req.EndBillDate + "','%Y-%m-%d')")
	if req.FilterCustomer != "" && req.FilterInvoice == "" {
		query.Where("sf_fpm_compare.payer_name_zsd001n LIKE ?", "%"+req.FilterCustomer+"%")
	} else if req.FilterCustomer == "" && req.FilterInvoice != "" {
		query.Where("sf_fpm_compare.billing_document_zsd001n LIKE ?", "%"+req.FilterInvoice+"%")
	} else if req.FilterCustomer != "" && req.FilterInvoice != "" {
		query.Where("(sf_fpm_compare.payer_name_zsd001n LIKE ? OR sf_fpm_compare.billing_document_zsd001n LIKE ?)",
			"%"+req.FilterCustomer+"%", "%"+req.FilterInvoice+"%")
	}
	query.Count(&totalData)
	return totalData
}

func (m mapProcessCompareRepositoryCon) GetCompareByInvoice(invoice string) DetailCompare {
	var records DetailCompare
	errSelect := m.mapProcessCompareRepositoryCon.Table("sf_fpm_compare").Select("sf_fpm_compare.*").
		Where("sf_fpm_compare.billing_document_zsd001n=?", invoice).First(&records)
	if errSelect.Error != nil{
		log.Error("error get data compare by invoice : ", errSelect.Error.Error())
	}
	return records
}

func (m mapProcessCompareRepositoryCon) GetCompareFpHisZsd081() []CompareZsd081Faktur {
	var records []CompareZsd081Faktur
	query := m.mapProcessCompareRepositoryCon.Table("sf_fpm_faktur").Select(
		"sf_fpm_faktur.*," +
			"GROUP_CONCAT(sf_dump_zsd081.`email_address`) AS email_receipt_fp," +
			"sf_dump_zsd081.create_on AS send_date_fp," +
			"sf_dump_zsd081.time AS send_time_fp").
		Joins(" inner join sf_dump_zsd081 ON sf_fpm_faktur.faktur = TRIM(LEADING '0' FROM `sf_dump_zsd081`.`doc_number`)")
	query.Where("sf_fpm_faktur.status_send_fp = ?", "no").
		Where("sf_fpm_faktur.status = ?", 1).
		Where("sf_dump_zsd081.type=?", "FPM").
		Group("sf_dump_zsd081.doc_number")
	err := query.Scan(&records)
	if err.Error != nil {
		log.Error(" error get fp in zsd081 : ", err.Error.Error())
	}
	return records
}

func (m mapProcessCompareRepositoryCon) UpdateInvAlreadySend(data models.DetailCompareModels) error {
	err := m.mapProcessCompareRepositoryCon.Save(&data)
	if err.Error != nil {
		log.Error("error update flag invoice already send to customer with data : ", data, err.Error.Error())
		return fmt.Errorf("%s", "error update flag invoice already send to customer")
	}
	return nil
}

func (m mapProcessCompareRepositoryCon) GetCompareInvZsd081(dateInput string) []DetailCompare {
	var records []DetailCompare
	query := m.mapProcessCompareRepositoryCon.Table("sf_fpm_compare").
		Select("sf_fpm_compare.Id," +
			" sf_fpm_compare.billing_document_zsd001n," +
			" sf_fpm_compare.billing_type_zsd001n," +
			" sf_fpm_compare.billing_date_zsd001n," +
			" sf_fpm_compare.dc_zsd001n," +
			" sf_fpm_compare.payer_zsd001n," +
			" sf_fpm_compare.payer_name_zsd001n," +
			" sf_fpm_compare.slor_zsd001n," +
			" sf_fpm_compare.district_zsd001n," +
			" sf_fpm_compare.material_number_zsd001n," +
			" sf_fpm_compare.material_desc_zsd001n," +
			" sf_fpm_compare.create_on_zsd001n," +
			" sf_fpm_compare.billing_doc_cancel," +
			" sf_fpm_compare.sts_cancel_inv," +
			" sf_fpm_compare.sts_email_inv_cancel," +
			" sf_fpm_compare.sts_send_inv," +
			" sf_fpm_compare.billing_number_zv60," +
			" sf_fpm_compare.billing_date_zv60," +
			" sf_fpm_compare.fp_number_zv60," +
			" sf_fpm_compare.fp_created_date_zv60," +
			" sf_fpm_compare.payer_zv60," +
			" sf_fpm_compare.name_zv60," +
			" sf_fpm_compare.npwp_zv60," +
			" sf_fpm_compare.material_zv60," +
			" sf_fpm_compare.sts_compare," +
			" sf_fpm_compare.id_history_fpm," +
			" sf_fpm_compare.sts_email_compare," +
			" GROUP_CONCAT(sf_dump_zsd081.`email_address`) AS email_receipt_inv," +
			" sf_dump_zsd081.create_on AS send_date_inv," +
			" sf_dump_zsd081.time AS send_time_inv," +
			" sf_fpm_compare.id_his_email").
		Joins("inner join sf_dump_zsd081 on sf_fpm_compare.billing_document_zsd001n = sf_dump_zsd081.doc_number")
	query.Where("sf_fpm_compare.sts_send_inv = ?", "no").
		Where("sf_dump_zsd081.type=?", "INV").
		Group("sf_dump_zsd081.doc_number")
	errSelect := query.Scan(&records)
	if errSelect.Error != nil{
		log.Error("error data compare vs zsd081 by invoice : ", errSelect.Error.Error())
	}
	return records
}

func (m mapProcessCompareRepositoryCon) GetHistoryFpCancel() []FakturHistory {
	var recordHis []FakturHistory

	errSelect:=m.mapProcessCompareRepositoryCon.Table("sf_fpm_faktur").Select("sf_fpm_faktur.ID,"+
		"sf_fpm_faktur.inv_number,"+
		"sf_fpm_faktur.inv_dc,"+
		"sf_fpm_faktur.faktur,"+
		"sf_fpm_faktur.faktur_created,"+
		"sf_fpm_faktur.status,"+
		"sf_fpm_faktur.create_at,"+
		"sf_fpm_faktur.update_at").
		Where("status=?", 0).
		Where("sf_fpm_faktur.update_at IS NOT NULL").
		Where("sf_fpm_faktur.inv_dc=?", "10").
		Where("DATE(sf_fpm_faktur.update_at) = CURRENT_DATE()").Scan(&recordHis)
	if errSelect.Error != nil{
		log.Error("error get data fp cancel : ",errSelect.Error.Error())
	}
	return recordHis
}

func (m mapProcessCompareRepositoryCon) SaveHistoryFaktur(data models.FakturModels) error {
	err := m.mapProcessCompareRepositoryCon.Create(&data)
	if err.Error != nil {
		log.Error("error save history faktur, with data: ", data, err.Error.Error())
		return fmt.Errorf("%s", "err save history faktur")
	}
	return nil
}

func (m mapProcessCompareRepositoryCon) UpdateHistoryFaktur(billNumber string) error {
	log.Info(" non active fp base on : ", billNumber)
	updateData := map[string]interface{}{
		"update_at": time.Now(),
		"status":    0,
	}
	err := m.mapProcessCompareRepositoryCon.Table("sf_fpm_faktur").
		Where("sf_fpm_faktur.inv_number=?", billNumber).
		Where("sf_fpm_faktur.update_at IS NULL").Updates(&updateData)
	if err.Error != nil {
		log.Error("error update history update", err.Error.Error())
		return err.Error
	}
	return nil
}

func (m mapProcessCompareRepositoryCon) GetHistoryFaktur(faktur string) FakturHistory {
	var records FakturHistory
	errSelect:=m.mapProcessCompareRepositoryCon.Model(&models.FakturModels{}).
		Where("faktur=?", faktur).First(&records)
	if errSelect.Error != nil{
		log.Error("error get his faktur : ", errSelect.Error.Error())
	}
	return records
}

func (m mapProcessCompareRepositoryCon) CreateDataCompare(compareModels []models.DetailCompareModels) error {
	db := m.mapProcessCompareRepositoryCon
	err := db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&compareModels)
	if err.Error != nil {
		log.Error("error copy bulk zsd001n to compare ", err.Error.Error())
		return fmt.Errorf("%s", "Failed to insert batch copy data")
	}
	return nil
}

func InstanceProcessCompareRepository(db *gorm.DB) ProcessCompareRepository {
	return &mapProcessCompareRepositoryCon{
		mapProcessCompareRepositoryCon: db,
	}
}
func (m mapProcessCompareRepositoryCon) GetDataZsd001nByCreateOn(date string) []Zsd001n {
	var records []Zsd001n
	query := m.mapProcessCompareRepositoryCon.Table("sf_dump_zsd001n").Select(
		"sf_dump_zsd001n.billing_document AS billing_document_zsd001n,"+
			"sf_dump_zsd001n.billing_type AS billing_type_zsd001n,"+
			"sf_dump_zsd001n.billing_date AS billing_date_zsd001n,"+
			"sf_dump_zsd001n.dc AS dc_zsd001n,"+
			"sf_dump_zsd001n.payer AS payer_zsd001n,"+
			"sf_dump_zsd001n.payer_name AS payer_name_zsd001n,"+
			"sf_dump_zsd001n.slor AS slor_zsd001n,"+
			"sf_dump_zsd001n.district_name AS district_zsd001n,"+
			"sf_dump_zsd001n.material_number AS material_number_zsd001n,"+
			"sf_dump_zsd001n.material_desc AS material_desc_zsd001n,"+
			"sf_dump_zsd001n.created_on AS create_on_zsd001n").
		Where("SUBSTR(sf_dump_zsd001n.billing_document,1,1)=?", "7")
	if date != "" {
		query.Where("sf_dump_zsd001n.created_on=?", date)
	} else {
		query.Where("STR_TO_DATE(sf_dump_zsd001n.created_on,'%Y%m%d') = CURRENT_DATE() - INTERVAL 1 DAY")
	}
	errSelect :=query.Group("sf_dump_zsd001n.billing_document, sf_dump_zsd001n.billing_type").Scan(&records)
	if errSelect.Error != nil{
		log.Error("error get data zsd001n for copy : ", errSelect.Error.Error())
	}
	return records
}

func (m mapProcessCompareRepositoryCon) UpdateFpAlreadySend(dataFp models.FakturModels) error {

	updateData := map[string]interface{}{
		"status_send_fp":   dataFp.StatusSendFp,
		"send_date_fp" : dataFp.SendDateFp,
		"send_time_fp" : dataFp.SendTimeFp,
		"email_receipt_fp": dataFp.EmailReceiptFp,
	}

	log.Info("data update status send fp history : ", dataFp)
	errUpdtFp := m.mapProcessCompareRepositoryCon.Table("sf_fpm_faktur").
		Where("sf_fpm_faktur.id=?", dataFp.ID).Updates(&updateData)
	if errUpdtFp.Error != nil {
		log.Error("error update history update with data : ", dataFp, errUpdtFp.Error.Error())
		return fmt.Errorf("error update history update")
	}
	return nil
}

func (m mapProcessCompareRepositoryCon) GetFpByCreated(dateInput string) []GetFpCancel {
	var records []GetFpCancel
	query := m.mapProcessCompareRepositoryCon.Table("sf_dump_zv60").Select(
		"sf_fpm_compare.billing_document_zsd001n AS inv_number, " +
			"sf_fpm_compare.dc_zsd001n AS inv_channel, " +
			"sf_fpm_compare.fp_number_zv60 AS fp_old, " +
			"sf_fpm_compare.fp_created_date_zv60 AS fp_created_old, " +
			"sf_dump_zv60.fp_number AS new_fp, " +
			"sf_dump_zv60.fp_created_date AS new_fp_created",
	).Joins("inner join sf_fpm_compare on sf_dump_zv60.billing_number = sf_fpm_compare.billing_number_zv60")

	if dateInput != "" {
		query.Where("sf_dump_zv60.fp_created_date =?", dateInput)
	} else {
		query.Where("STR_TO_DATE(sf_dump_zv60.fp_created_date,'%d.%m.%Y') = CURRENT_DATE() - INTERVAL 1 DAY")
	}
	query.Group("sf_dump_zv60.billing_number").Scan(&records)
	return records
}

func (m mapProcessCompareRepositoryCon) UpdateInvoiceCancel(data models.DetailCompareModels) error {
	err := m.mapProcessCompareRepositoryCon.Save(&data)
	if err.Error != nil {
		log.Error("Error update invoice cancel ", data, "with error : "+err.Error.Error())
		return fmt.Errorf("%s", "Failed to update invoice cancel")
	}
	return nil
}

func (m mapProcessCompareRepositoryCon) CompareInvoiceCancel(date string) []DetailCompare {
	var records []DetailCompare
	query := m.mapProcessCompareRepositoryCon.Table("sf_fpm_compare").Select(
		"sf_fpm_compare.*,"+
			"sf_dump_zsd001n.billing_document AS inv_cancel").
		Joins("inner join sf_dump_zsd001n on sf_fpm_compare.billing_document_zsd001n = sf_dump_zsd001n.bill_cancelled_ref").
		Where("sf_fpm_compare.sts_cancel_inv != ?", "yes").
		Where("(SUBSTR(sf_dump_zsd001n.billing_document,1,1)=? AND sf_dump_zsd001n.billing_type NOT IN(?))", "8", "ZG2").
		Where("sf_fpm_compare.billing_type_zsd001n NOT IN (?,?,?)", "ZINC", "ZINP", "ZINQ")

	if date != "" {
		query.Where("sf_dump_zsd001n.created_on=?", date)
	} else {
		query.Where("STR_TO_DATE(sf_dump_zsd001n.created_on,'%Y%m%d') = CURRENT_DATE() - INTERVAL 1 DAY")
	}
	errSelect :=query.Group("sf_fpm_compare.billing_document_zsd001n, sf_fpm_compare.billing_type_zsd001n").Scan(&records)
	if errSelect.Error != nil{
		log.Error("error get compare inv cancel : ", errSelect.Error.Error())
	}
	return records
}

func (m mapProcessCompareRepositoryCon) UpdateCompareZsd001nZv60(data models.DetailCompareModels) error {
	err := m.mapProcessCompareRepositoryCon.Save(&data)
	if err.Error != nil {
		log.Error("Error save or update data compare zsd001n vs zv60 ", data, "with error : "+err.Error.Error())
		return fmt.Errorf("%s", "Failed to save data compare zsd001n vs zv60")
	}
	return nil
}

func (m mapProcessCompareRepositoryCon) CompareZsd001nZv60(dateInput string) []DetailCompare {
	var records []DetailCompare
	query := m.mapProcessCompareRepositoryCon.Table("sf_dump_zv60").Select(
		"sf_fpm_compare.id,"+
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
			" sf_dump_zv60.fp_created_date AS fp_created_date_zv60,"+
			" sf_dump_zv60.billing_number AS billing_number_zv60,"+
			" sf_dump_zv60.billing_date AS billing_date_zv60,"+
			" sf_dump_zv60.fp_number AS fp_number_zv60,"+
			" sf_dump_zv60.payer AS payer_zv60,"+
			" sf_dump_zv60.payer AS name_zv60,"+
			" sf_dump_zv60.npwp AS npwp_zv60,"+
			" sf_dump_zv60.material AS material_zv60,"+
			" sf_fpm_compare.id_history_fpm,"+
			" sf_fpm_compare.sts_compare,"+
			" sf_fpm_compare.sts_email_compare,"+
			" sf_fpm_compare.email_receipt_inv,"+
			" sf_fpm_compare.send_date_inv,"+
			" sf_fpm_compare.send_time_inv,"+
			" sf_fpm_compare.id_his_email").
		Joins("inner join sf_fpm_compare on sf_dump_zv60.billing_number = sf_fpm_compare.billing_document_zsd001n").
		Where("sf_fpm_compare.sts_compare=?", "not_match").
		Where("sf_fpm_compare.sts_cancel_inv=?", "no").
		Where("sf_fpm_compare.billing_type_zsd001n NOT IN (?,?,?)", "ZINC", "ZINP", "ZINQ")
	if dateInput != "" {
		query.Where("sf_dump_zv60.fp_created_date=?", dateInput)
	} else {
		query.Where("STR_TO_DATE(sf_dump_zv60.fp_created_date,'%d.%m.%Y') = CURRENT_DATE() - INTERVAL 1 DAY")
	}

	errSelect := query.Group("sf_dump_zv60.billing_number , sf_dump_zv60.fp_number").Scan(&records)
	if errSelect.Error != nil{
		log.Error("error compare vs zv60 : ", errSelect.Error.Error())
	}
	return records

}

type CompareZsd001nZv60 struct {
	IdCompare              int    `json:"idCompare"`
	BillingDocumentZsd001n string `json:"billingDocumentZsd001n"`
	BillingTypeZsd001n     string `json:"billingTypeZsd001n"`
	BillingDateZsd001n     string `json:"billingDateZsd001n"`
	DcZsd001n              string `json:"dCZsd001n"`
	PayerNameZsd001n       string `json:"payerNameZsd001n"`
	BillingNumberZv60      string `json:"billingNumberZv60"`
	BillingDateZv60        string `json:"billingDateZv60"`
	FpNumberZv60           string `json:"fpNumberZv60"`
	StsCompare             string `json:"statusCompare"`
}

type DetailCompare struct {
	Id                     int    `json:"id"`
	BillingDocumentZsd001n string `json:"billingDocumentZsd001n"`
	BillingTypeZsd001n     string `json:"billingTypeZsd001n"`
	BillingDateZsd001n     string `json:"billingDateZsd001n"`
	DcZsd001n              string `json:"dcZsd001n"`
	SlorZsd001n            string `json:"slorZsd001n"`
	PayerZsd001n           string `json:"payerZsd001n"`
	PayerNameZsd001n       string `json:"payerNameZsd001n"`
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
	InvCancel              string `json:"invCancel"`
	EmailReceiptInv        string `json:"emailReceiptInv"`
	SendDateInv            string `json:"sendDateInv"`
	SendTimeInv            string `json:"sendTimeInv"`
	IdHisEmail             int    `json:"idHisEmail"`
}
type GetFpCancel struct {
	InvNumber    string `json:"invNumber"`
	InvChannel   string `json:"invChannel"`
	FpOld        string `json:"fpOld"`
	FpCreatedOld string `json:"fpCreatedOld"`
	NewFp        string `json:"newFp"`
	NewFpCreated string `json:"newFpCreated"`
}

type Zsd001n struct {
	BillingDocumentZsd001n string `json:"billingDocumentZsd001n"`
	BillingTypeZsd001n     string `json:"billingTypeZsd001n"`
	BillingDateZsd001n     string `json:"billingDateZsd001n"`
	DcZsd001n              string `json:"dcZsd001n"`
	SlorZsd001n            string `json:"slorZsd001n"`
	PayerZsd001n           string `json:"payerZsd001n"`
	PayerNameZsd001n       string `json:"payerNameZsd001n"`
	DistrictZsd001n        string `json:"districtZsd001n"`
	CreateOnZsd001n        string `json:"createOnZsd001n"`
	MaterialNumberZsd001n  string `json:"materialNumberZsd001n"`
	MaterialDescZsd001n    string `json:"materialDescZsd001n"`
}

type FakturHistory struct {
	ID            int       `json:"id"`
	InvNumber     string    `json:"invNumber"`
	InvDc         string    `json:"invDc"`
	Faktur        string    `json:"faktur"`
	FakturCreated string    `json:"fakturCreated"`
	StatusSendFp  string    `json:"statusSendFp"`
	Status        int       `json:"status"`
	UpdateAt      time.Time `json:"updateAt"`
	CreateAt      time.Time `json:"createAt"`
}

type CompareZsd081Faktur struct {
	ID             int       `json:"id"`
	InvNumber      string    `json:"invNumber"`
	InvDc          string    `json:"invDc"`
	Faktur         string    `json:"faktur"`
	FakturCreated  string    `json:"fakturCreated"`
	StatusSendFp   string    `json:"statusSendFp"`
	Status         int       `json:"status"`
	UpdateAt       time.Time `json:"updateAt"`
	CreateAt       time.Time `json:"createAt"`
	EmailReceiptFp string    `json:"emailReceiptFp"`
	SendDateFp    string    `json:"sendDateFp"`
	SendTimeFp    string    `json:"sendDateFp"`
}

type DataEmailCancel struct {
	BillingDocumentZsd001n string `json:"billingDocumentZsd001n"`
	BillingDateZsd001n     string `json:"billingDateZsd001n"`
	FpNumberZv60           string `json:"fpNumberZv60"`
}
