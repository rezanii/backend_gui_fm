package models

import "time"

type EmailNotificationModels struct {
	Id             int       `gorm:"primary_key:auto_increment"`
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

func (EmailNotificationModels) TableName() string {
	return "sf_fpm_email_notification"
}
