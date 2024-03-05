package repository

import (
	"backend_gui/dto"
	"backend_gui/models"
	"fmt"
	"github.com/mashingan/smapping"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UploadRepository interface {
	SaveDataUploadFilePjk(history models.FpmHistoryFilePjk) error
	GetDataHistoryUploadFpm(start int, end int, search string, username string, roleUser string) []dto.RespDataFpDto
	GetTotalDataFpm(search string, username string, roleUser string) int64
	GetHistoryFpmByFileName(fileName string, username string) int
	CheckUserAlreadyExist(username string) UserWithRole
	GetDownloadFpmById(listId []int, username string) []FpmDownload
}

type mapUploadRepositoryCon struct {
	mapUploadRepositoryCon *gorm.DB
}

func (m mapUploadRepositoryCon) GetDownloadFpmById(listId []int, username string) []FpmDownload{
	var records []FpmDownload
	errSelect :=m.mapUploadRepositoryCon.Table("sf_fpm_compare").
		Select(
			"sf_fpm_history_upload_pjk.file_name_upload,"+
				"sf_fpm_history_upload_pjk.jenis_faktur_pajak,"+
				"sf_fpm_history_upload_pjk.status,"+
				"sf_fpm_faktur.inv_dc").
		Joins("inner join sf_fpm_faktur on sf_fpm_compare.billing_document_zsd001n = sf_fpm_faktur.inv_number").
		Joins("inner join sf_fpm_history_upload_pjk on sf_fpm_faktur.faktur = sf_fpm_history_upload_pjk.no_faktur").
		Where("sf_fpm_compare.id IN (?)", listId).Scan(&records)
	if errSelect.Error != nil{
		log.Error("Zip file, get data for download file, username "+username+", error : ", errSelect.Error.Error())
	}
	return records
}

func (m mapUploadRepositoryCon) CheckUserAlreadyExist(username string) UserWithRole {
	var records UserWithRole
	errSelect :=m.mapUploadRepositoryCon.Table("sf_fpm_user").Select("sf_fpm_user.username, sf_fpm_user.id,"+
		"sf_fpm_role.role_user").
		Joins("inner join sf_fpm_role on sf_fpm_role.id = sf_fpm_user.id_role_user").
		Where("sf_fpm_user.username=?", username).
		Where("sf_fpm_user.sts_active=?", "1").First(&records)
	if errSelect.Error != nil{
		log.Error("check user at db,username "+username+", error : ", errSelect.Error.Error())
	}
	return records
}

func (m mapUploadRepositoryCon) GetHistoryFpmByFileName(fileName string, username string) int {
	data := dto.RespFpmHistoryDto{}
	fpmHistory := models.FpmHistoryFilePjk{}
	m.mapUploadRepositoryCon.Where("file_name_upload=?", fileName).First(&fpmHistory)
	err := smapping.FillStructByTags(&data, smapping.MapTags(&fpmHistory, "json"), "json")
	if err != nil {
		log.Error("get history file fpm, username "+username+"Error mapping value : ", err.Error())
	}
	return data.ID
}

func (m mapUploadRepositoryCon) GetTotalDataFpm(search string, username string, roleUser string) int64 {
	var total int64
	query := m.mapUploadRepositoryCon.Model(&models.FpmHistoryFilePjk{})
	if roleUser!= "ADMIN"{
		query.Where("upload_by=?", username)
	}
	if search != "" {
		query.Where("(no_faktur LIKE ? OR jenis_faktur_pajak LIKE ? OR status LIKE ? )", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	query = query.Count(&total)
	return total
}

func (m mapUploadRepositoryCon) GetDataHistoryUploadFpm(start int, end int, search string, username string, roleUser string) ([]dto.RespDataFpDto ) {
	data := []dto.RespDataFpDto{}
	query := m.mapUploadRepositoryCon.Table("sf_fpm_history_upload_pjk").
		Select("sf_fpm_history_upload_pjk.id,sf_fpm_history_upload_pjk.upload_by,"+
			" sf_fpm_history_upload_pjk.jenis_faktur_pajak, "+
			"sf_fpm_history_upload_pjk.file_name_upload, "+
			"sf_fpm_history_upload_pjk.status, "+
		    "sf_fpm_faktur.inv_dc, "+
			"sf_fpm_faktur.inv_number,"+
			"DATE_FORMAT(STR_TO_DATE(sf_fpm_compare.billing_date_zsd001n,'%Y%m%d'),'%Y-%m-%d') AS inv_date,"+
			"sf_fpm_history_upload_pjk.no_faktur,"+
			"DATE_FORMAT(STR_TO_DATE(sf_fpm_faktur.faktur_created,'%d.%m.%Y'),'%Y-%m-%d') AS fp_created_date,"+
			"DATE_FORMAT(sf_fpm_history_upload_pjk.create_at, '%Y-%m-%d %H:%I:%S') AS create_at,"+
			"DATE_FORMAT(sf_fpm_history_upload_pjk.update_at, '%Y-%m-%d %H:%I:%S') AS update_at").
		Joins("left join sf_fpm_faktur on sf_fpm_history_upload_pjk.no_faktur = sf_fpm_faktur.faktur").
		Joins("left join sf_fpm_compare on sf_fpm_compare.billing_document_zsd001n = sf_fpm_faktur.inv_number")
	if roleUser != "ADMIN"{
		query.Where("upload_by = ?", username)
	}

	if search != "" {
		query.Where("(sf_fpm_history_upload_pjk.no_faktur LIKE ? OR sf_fpm_history_upload_pjk.jenis_faktur_pajak LIKE ? OR sf_fpm_history_upload_pjk.status LIKE ? )", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	err:= query.Limit(end).Offset(start).Scan(&data)
	if err.Error != nil{
		log.Error("get his upload faktur, username "+username+", error : ", err.Error.Error())
	}
	return data
}

func (m mapUploadRepositoryCon) SaveDataUploadFilePjk(historyModel models.FpmHistoryFilePjk) error {
	err := m.mapUploadRepositoryCon.Save(&historyModel)
	if err.Error != nil {
		log.Error("errorr save fpm data : ", historyModel, err.Error.Error())
		return fmt.Errorf("%s", "Failed to save data upload file")
	}
	return nil

}

func InstanceUploadRepository(db *gorm.DB) UploadRepository {
	return &mapUploadRepositoryCon{
		mapUploadRepositoryCon: db,
	}
}

type UserWithRole struct {
	Id       int
	Username string
	RoleUser string
}

type FpmDownload struct {
	FileNameUpload   string
	JenisFakturPajak string
	Status           string
	InvDc            string
}

