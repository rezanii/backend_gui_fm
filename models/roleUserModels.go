package models

type RoleUser struct {
	Id       int    `gorm:"primary_key:auto_increment" json:"id"`
	RoleUser string `json:"roleUser"`
	DeleteAt int    `json:"deleteAt"`
}

func (RoleUser) TableName() string {
	return "sf_fpm_role"
}
