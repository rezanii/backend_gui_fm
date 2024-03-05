package repository

import (
	"backend_gui/dto"
	"backend_gui/models"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserFpmRepository interface {
	CheckUserByUsername(username string) models.UserFpmModels
	CheckRoleUserExist(roleUser string) models.RoleUser
	CheckUserById(idUser int) models.UserFpmModels
	DeleteUserFpm(data models.UserFpmModels) error
	SaveUserFpm(data models.UserFpmModels) error
	GetTotalUserFpm(search string) int64
	GetHistoryUserFpm(start int, end int, search string) []dto.RespHisUserFpm
	IsRoleAdmin(username string) models.UserFpmModels
	GetListRoleUser() []dto.RespRoleUser
}

type mapUserFpmRepositoryCon struct {
	mapUserFpmRepositoryCon *gorm.DB
}

func (m mapUserFpmRepositoryCon) GetListRoleUser() []dto.RespRoleUser{
	var data []dto.RespRoleUser
	m.mapUserFpmRepositoryCon.Model(&models.RoleUser{}).Where("delete_at=?", 0).Scan(&data)
	return data

}

func (m mapUserFpmRepositoryCon) IsRoleAdmin(username string) models.UserFpmModels {
	var recordsUser models.UserFpmModels
	m.mapUserFpmRepositoryCon.Table("sf_fpm_user").Select(
		"sf_fpm_user.id,"+
			"sf_fpm_user.username,"+
			"sf_fpm_user.id_role_user,"+
			"sf_fpm_user.sts_active").
		Joins("inner join sf_fpm_role on sf_fpm_user.id_role_user = sf_fpm_role.id").
		Where("sf_fpm_user.username = ?", username).
		Where("sf_fpm_role.role_user=?", "ADMIN").
		Where("sf_fpm_user.sts_active", "1").First(&recordsUser)
	return recordsUser
}

func (m mapUserFpmRepositoryCon) GetHistoryUserFpm(start int, end int, search string) []dto.RespHisUserFpm {
	var recordHis []dto.RespHisUserFpm
	query := m.mapUserFpmRepositoryCon.Table("sf_fpm_user").Select(
		"sf_fpm_user.id," +
			"sf_fpm_user.username," +
			"sf_fpm_role.role_user," +
			"IF(sf_fpm_user.sts_active = '1', 'active', ' ') AS sts_active").
		Joins("inner join sf_fpm_role on sf_fpm_user.id_role_user = sf_fpm_role.id").
		Where("sf_fpm_user.sts_active = ?", "1")
	if search != "" {
		query.Where("(sf_fpm_user.username LIKE ? OR sf_fpm_user.sts_active LIKE ? OR sf_fpm_role.role_user LIKE ?)",
			"%"+search+"%", "%"+search+"%s", "%"+search+"%")
	}
	query.Order("sf_fpm_user.dtm_create DESC")
	err:= query.Limit(end).Offset(start).Scan(&recordHis)
	if err.Error != nil{
		log.Error("get history user, error :", err.Error.Error())
	}
	return recordHis
}

func (m mapUserFpmRepositoryCon) GetTotalUserFpm(search string) int64 {
	var totalUser int64
	query := m.mapUserFpmRepositoryCon.Table("sf_fpm_user").Select("sf_fpm_user.id").
		Joins("inner join sf_fpm_role on sf_fpm_user.id_role_user = sf_fpm_role.id").
		Where("sf_fpm_user.sts_active = ?", "1")
	if search != "" {
		query.Where("(sf_fpm_user.username LIKE ? OR sf_fpm_user.sts_active LIKE ? OR sf_fpm_role.role_user LIKE ?)",
			"%"+search+"%", "%"+search+"%s", "%"+search+"%")
	}
	query.Count(&totalUser)
	return totalUser
}

func (m mapUserFpmRepositoryCon) DeleteUserFpm(data models.UserFpmModels) error {
	errUpdate := m.mapUserFpmRepositoryCon.Save(&data)
	if errUpdate.Error != nil {
		log.Error("error non active user : ", errUpdate.Error.Error())
		return fmt.Errorf("%s", "error delete user fpm")
	}
	return nil
}

func (m mapUserFpmRepositoryCon) CheckUserById(idUser int) models.UserFpmModels {
	var data models.UserFpmModels
	m.mapUserFpmRepositoryCon.Where("id=?", idUser).First(&data)
	return data
}

func (m mapUserFpmRepositoryCon) CheckRoleUserExist(roleUser string) models.RoleUser {
	var data models.RoleUser
	m.mapUserFpmRepositoryCon.Where("role_user=?", roleUser).First(&data)
	return data
}

func (m mapUserFpmRepositoryCon) SaveUserFpm(data models.UserFpmModels) error {
	errSave := m.mapUserFpmRepositoryCon.Save(&data)
	if errSave.Error != nil {
		log.Error("error save user fpm : ", errSave.Error.Error())
		return fmt.Errorf("%s", "failed save username")
	}
	return nil
}

func (m mapUserFpmRepositoryCon) CheckUserByUsername(username string) models.UserFpmModels {
	var dataUser models.UserFpmModels
	m.mapUserFpmRepositoryCon.Where("username=?", username).Where("sts_active=?","1").First(&dataUser)
	return dataUser
}

func InstanceUserFpmRepository(db *gorm.DB) UserFpmRepository {
	return &mapUserFpmRepositoryCon{
		mapUserFpmRepositoryCon: db,
	}
}
