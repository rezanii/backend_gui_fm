package dto

type DtoUserRoleFpm struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	IdRoleUser string `json:"idRoleUser"`
	StsActive  string `json:"stsActive"`
}
