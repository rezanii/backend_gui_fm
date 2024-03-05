package models

type UserFpmModels struct {
	ID         int    `gorm:"primary_key:auto_increment" json:"id"`
	Username   string `json:"username"`
	IdRoleUser string    `json:"idRoleUser"`
	StsActive  string `json:"stsActive"`
}

func (UserFpmModels) TableName() string {
	return "sf_fpm_user"
}
