package controller

import (
	"backend_gui/dto"
	"backend_gui/models"
	"backend_gui/repository"
	"backend_gui/utils"
	b64 "encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type EmailController interface {
	EmailCancelInv(ctx *gin.Context)
	EmailInvNoFp(ctx *gin.Context)
	EmailInvNotSend(ctx *gin.Context)
	EmailNoFileFp(ctx *gin.Context)
	EmailFpCancel(ctx *gin.Context)
	EmailFpNotSend(ctx *gin.Context)
}

type emailController struct {
	mapEmailRepo repository.EmailRepository
}
/*
	Method Email FP not send to customer
*/
func (e emailController) EmailFpNotSend(ctx *gin.Context) {
	var reqDateInput dto.ReqInputDate
	_ = ctx.ShouldBindJSON(&reqDateInput)
	dateInput := ""
	dateFormat := ""
	if reqDateInput.DateInput != "" {
		dateTemp := strings.Split(reqDateInput.DateInput, ".")
		dateInput = dateTemp[2] + "-" + dateTemp[1] + "-" + dateTemp[0]
		dateFormat = reqDateInput.DateInput
	} else {
		dateTempFormat := strings.Split(utils.ConvertDateToString(time.Now(), utils.TYPE_DATE_RFC3339), "-")
		dateFormat = dateTempFormat[2] + "." + dateTempFormat[1] + "." + dateTempFormat[0]
	}
	dataFpNotSend := e.mapEmailRepo.GetDataFpNotSend(dateInput)

	log.Info("Size data Fp not send : ", len(dataFpNotSend))
	log.Info("process email Fp not send, receipt : ", os.Getenv("EMAIL_TO_NOT_SEND_FP"))
	log.Info("process Fp not send, CC : ", os.Getenv("EMAIL_CC_NOT_SEND_FP"))
	log.Info("process Fp not send, sender : ", os.Getenv("EMAIL_SENDER"))
	log.Info("process Fp not send, subject : ", os.Getenv("SUBJECT_EMAIL_FP_NOT_SEND") + " " + dateFormat)
	if len(dataFpNotSend) > 0 {
		bodyEmail := utils.TemplateEmail(dataFpNotSend, utils.TYPE_EMAIL_FP_NOT_SEND)
		if bodyEmail != "" {
			emailNotif := models.EmailNotificationModels{
				//Recipient:    os.Getenv("EMAIL_TO"),
				//Cc:           os.Getenv("EMAIL_CC"),
				Recipient:    os.Getenv("EMAIL_TO_NOT_SEND_FP"),
				Cc:           os.Getenv("EMAIL_CC_NOT_SEND_FP"),
				Sender:       os.Getenv("EMAIL_SENDER"),
				Subject:      os.Getenv("SUBJECT_EMAIL_FP_NOT_SEND") + " " + dateFormat,
				LogStatus:    0,
				Body:         b64.StdEncoding.EncodeToString([]byte(bodyEmail)),
				SendDateTime: time.Now(),
			}
			dataEmail := e.mapEmailRepo.SaveEmailNotif(emailNotif)
			if dataEmail.Id > 0 {
				respEmail, errSendEmail := utils.EmailNotificationService(dataEmail, int(dataEmail.Id), []string{})
				log.Info("Send email data fp not send with response: ", respEmail)
				if errSendEmail != nil && respEmail.ResponseCode != 200 {
					log.Error("Error send email data fp not send : ", errSendEmail.Error())
				} else {
					for i := 0; i < len(dataFpNotSend); i++ {
						intVar, errorParser := strconv.Atoi(fmt.Sprintf("%v", dataFpNotSend[i]["id_compare"]))

						if errorParser != nil {
							log.Info("update data compare, id his email : ", errorParser.Error())
						}
						errUpdate := e.mapEmailRepo.UpdateStatusEmail(dataEmail.Id, intVar)
						if errUpdate != nil {
							log.Info("update data compare, id his email : ", errUpdate.Error())
						}

					}

				}
			}
		} else {
			log.Error("Email data fp not send : ", "Error template email")
		}

	} else {
		log.Info("Email data fp not send : ", "data fp cancel not found")
	}
}
/*
	Method email send fp cancel
*/

func (e emailController) EmailFpCancel(ctx *gin.Context) {
	var reqDateInput dto.ReqInputDate
	_ = ctx.ShouldBindJSON(&reqDateInput)
	dateInput := ""
	dateFormat := ""
	if reqDateInput.DateInput != "" {
		dateTemp := strings.Split(reqDateInput.DateInput, ".")
		dateInput = dateTemp[2] + "-" + dateTemp[1] + "-" + dateTemp[0]
		dateFormat = reqDateInput.DateInput
	} else {
		dateTempFormat := strings.Split(utils.ConvertDateToString(time.Now(), utils.TYPE_DATE_RFC3339), "-")
		dateFormat = dateTempFormat[2] + "." + dateTempFormat[1] + "." + dateTempFormat[0]
	}

	dataFpCancel := e.mapEmailRepo.GetDataFpCancel(dateInput)
	log.Info("Size data Fp cancel : ", len(dataFpCancel))
	log.Info("process email fp cancel, receipt : ", os.Getenv("EMAIL_TO_CANCEL_FAKTUR"))
	log.Info("process fp cancel, CC : ", os.Getenv("EMAIL_CC_CANCEL_FAKTUR"))
	log.Info("process fp cancel, sender : ", os.Getenv("EMAIL_SENDER"))
	log.Info("process fp cancel, subject : ", os.Getenv("SUBJECT_EMAIL_FP_CANCEL") + " " + dateFormat)

	if len(dataFpCancel) > 0 {
		bodyEmail := utils.TemplateEmail(dataFpCancel, utils.TYPE_EMAIL_FP_CANCEL)
		if bodyEmail != "" {
			emailNotif := models.EmailNotificationModels{
				//Recipient:    os.Getenv("EMAIL_TO"),
				//Cc:           os.Getenv("EMAIL_CC"),
				Recipient:    os.Getenv("EMAIL_TO_CANCEL_FAKTUR"),
				Cc:           os.Getenv("EMAIL_CC_CANCEL_FAKTUR"),
				Sender:       os.Getenv("EMAIL_SENDER"),
				Subject:      os.Getenv("SUBJECT_EMAIL_FP_CANCEL") + " " + dateFormat,
				LogStatus:    0,
				Body:         b64.StdEncoding.EncodeToString([]byte(bodyEmail)),
				SendDateTime: time.Now(),
			}
			dataEmail := e.mapEmailRepo.SaveEmailNotif(emailNotif)
			if dataEmail.Id > 0 {
				respEmail, errSendEmail := utils.EmailNotificationService(dataEmail, int(dataEmail.Id), []string{})
				log.Info("Send email fp cancel with response: ", respEmail)
				if errSendEmail != nil && respEmail.ResponseCode != 200 {
					log.Error("Error send email fp cancel : ", errSendEmail.Error())
				} else {
					for i := 0; i < len(dataFpCancel); i++ {
						intVar, errorParser := strconv.Atoi(fmt.Sprintf("%v", dataFpCancel[i]["id_compare"]))

						if errorParser != nil {
							log.Info("update data compare, id his email : ", errorParser.Error())
						}
						errUpdate := e.mapEmailRepo.UpdateStatusEmail(dataEmail.Id, intVar)
						if errUpdate != nil {
							log.Info("update data compare, id his email : ", errUpdate.Error())
						}

					}
				}
			}
		} else {
			log.Error("Email fp cancel not send : ", "Error template email")
		}

	} else {
		log.Info("Email fp cancel not send : ", "data fp cancel not found")
	}

}
/*
	Method email invoice not have file e-faktur
*/

func (e emailController) EmailNoFileFp(ctx *gin.Context) {
	var reqDateInput dto.ReqInputDate
	_ = ctx.ShouldBindJSON(&reqDateInput)
	dateInput := ""
	dateFormat := ""
	if reqDateInput.DateInput != "" {
		dateTemp := strings.Split(reqDateInput.DateInput, ".")
		dateInput = dateTemp[2] + "-" + dateTemp[1] + "-" + dateTemp[0]
		dateFormat = reqDateInput.DateInput
	} else {
		dateTempFormat := strings.Split(utils.ConvertDateToString(time.Now(), utils.TYPE_DATE_RFC3339), "-")
		dateFormat = dateTempFormat[2] + "." + dateTempFormat[1] + "." + dateTempFormat[0]

	}
	dataFp := e.mapEmailRepo.GetDataNoFileFp(dateInput)
	log.Info("Size data Fp not have file faktur : ", len(dataFp))
	log.Info("process email invoice no file fp, receipt : ", os.Getenv("EMAIL_TO_NO_FILE_FP"))
	log.Info("process email invoice no file fp, CC : ", os.Getenv("EMAIL_CC_NO_FILE_FP"))
	log.Info("process email invoice no file fp, sender : ", os.Getenv("EMAIL_SENDER"))
	log.Info("process email invoice no file fp, subject : ", os.Getenv("SUBJECT_EMAIL_NO_FILE_FP") + " " + dateFormat)
	if len(dataFp) > 0 {
		bodyEmail := utils.TemplateEmail(dataFp, utils.TYPE_EMAIL_NO_FILE_FP)
		if bodyEmail != "" {
			emailNotif := models.EmailNotificationModels{
				//Recipient:    os.Getenv("EMAIL_TO"),
				//Cc:           os.Getenv("EMAIL_CC"),
				Recipient:    os.Getenv("EMAIL_TO_NO_FILE_FP"),
				Cc:           os.Getenv("EMAIL_CC_NO_FILE_FP"),
				Sender:       os.Getenv("EMAIL_SENDER"),
				Subject:      os.Getenv("SUBJECT_EMAIL_NO_FILE_FP") + " " + dateFormat,
				LogStatus:    0,
				Body:         b64.StdEncoding.EncodeToString([]byte(bodyEmail)),
				SendDateTime: time.Now(),
			}
			dataEmail := e.mapEmailRepo.SaveEmailNotif(emailNotif)
			if dataEmail.Id > 0 {
				respEmail, errSendEmail := utils.EmailNotificationService(dataEmail, int(dataEmail.Id), []string{})
				log.Info("Send email inv no fp with response: ", respEmail)
				if errSendEmail != nil && respEmail.ResponseCode != 200 {
					log.Error("Error send email inv no fp : ", errSendEmail.Error())
				} else {
					for i := 0; i < len(dataFp); i++ {
						intVar, errorParser := strconv.Atoi(fmt.Sprintf("%v", dataFp[i]["id_compare"]))

						if errorParser != nil {
							log.Info("update data compare, id his email : ", errorParser.Error())
						}
						errUpdate := e.mapEmailRepo.UpdateStatusEmail(dataEmail.Id, intVar)
						if errUpdate != nil {
							log.Info("update data compare, id his email : ", errUpdate.Error())
						}

					}
				}
			}
		} else {
			log.Error("Email Inv not Fp not send : ", "Error template email")
		}

	} else {
		log.Info("Email Inv not Fp not send : ", "data invoice no fp not found")
	}

}
/*
	Methhod process email invoice not match with zv60
*/
func (e emailController) EmailInvNoFp(ctx *gin.Context) {
	var reqDateInput dto.ReqInputDate
	_ = ctx.ShouldBindJSON(&reqDateInput)
	dateEmail := ""
	dateFormat := ""
	if reqDateInput.DateInput != "" {
		dateTemp := strings.Split(reqDateInput.DateInput, ".")
		//dateInput = dateTemp[2] + dateTemp[1] + dateTemp[0]
		//dateTemp := strings.Split(reqCopyData.DateInput, ".")
		dateInput := dateTemp[2] + "-" + dateTemp[1] + "-" + dateTemp[0]
		dateConvert, errDateParse := time.Parse(utils.TYPE_DATE_RFC3339, dateInput)
		if errDateParse != nil {
			log.Error("Error parsing date input process email Inv No Fp: ", errDateParse.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("Error Parsing date input", utils.ERR_GLOBAL))
			return
		}
		dateSubstracOne := dateConvert.AddDate(0, 0, -1)
		dateConvertSts := utils.ConvertDateToString(dateSubstracOne, utils.TYPE_DATE_RFC3339)
		dateSplit := strings.Split(dateConvertSts, "-")
		dateEmail = dateSplit[0]  + dateSplit[1]  + dateSplit[2]

		dateFormat = reqDateInput.DateInput
	} else {
		dateTempFormat := strings.Split(utils.ConvertDateToString(time.Now(), utils.TYPE_DATE_RFC3339), "-")
		dateFormat = dateTempFormat[2] + "." + dateTempFormat[1] + "." + dateTempFormat[0]
	}
	dataInvNoFp := e.mapEmailRepo.GetDataInvNoFp(dateEmail)
	log.Info("Size data invoice no fp : ", len(dataInvNoFp))
	log.Info("process email invoice no fp, receipt : ", os.Getenv("EMAIL_TO_NO_FAKTUR"))
	log.Info("process email invoice no fp, CC : ", os.Getenv("EMAIL_TO_NO_FAKTUR"))
	log.Info("process email invoice sender : ", os.Getenv("EMAIL_SENDER"))
	log.Info("process email, subject : ", os.Getenv("SUBJECT_EMAIL_INV_NO_FP") + " " + dateFormat)
	if len(dataInvNoFp) > 0 {
		bodyEmail := utils.TemplateEmail(dataInvNoFp, utils.TYPE_EMAIL_NO_FP)
		if bodyEmail != "" {
			emailNotif := models.EmailNotificationModels{
				//Recipient:    os.Getenv("EMAIL_TO"),
				//Cc:           os.Getenv("EMAIL_CC"),
				Recipient: os.Getenv("EMAIL_TO_NO_FAKTUR"),
				Cc: os.Getenv("EMAIL_CC_NO_FAKTUR"),
				Sender:       os.Getenv("EMAIL_SENDER"),
				Subject:      os.Getenv("SUBJECT_EMAIL_INV_NO_FP") + " " + dateFormat,
				LogStatus:    0,
				Body:         b64.StdEncoding.EncodeToString([]byte(bodyEmail)),
				SendDateTime: time.Now(),
			}
			dataEmail := e.mapEmailRepo.SaveEmailNotif(emailNotif)
			if dataEmail.Id > 0 {
				respEmail, errSendEmail := utils.EmailNotificationService(dataEmail, int(dataEmail.Id), []string{})
				log.Info("Send email inv no fp with response: ", respEmail)
				if errSendEmail != nil && respEmail.ResponseCode != 200 {
					log.Error("Error send email inv no fp : ", errSendEmail.Error())
				} else {
					for i := 0; i < len(dataInvNoFp); i++ {
						intVar, errorParser := strconv.Atoi(fmt.Sprintf("%v", dataInvNoFp[i]["id_compare"]))

						if errorParser != nil {
							log.Info("update data compare, id his email : ", errorParser.Error())
						}
						errUpdate := e.mapEmailRepo.UpdateStatusEmail(dataEmail.Id, intVar)
						if errUpdate != nil {
							log.Info("update data compare, id his email : ", errUpdate.Error())
						}

					}
				}
			}
		} else {
			log.Error("Email Inv not Fp not send : ", "Error template email")
		}

	} else {
		log.Info("Email Inv not Fp not send : ", "data invoice no fp not found")
	}
}
/*
	Method send email cancel base on date now
*/
func (e emailController) EmailCancelInv(ctx *gin.Context) {
	var reqDateInput dto.ReqInputDate
	_ = ctx.ShouldBindJSON(&reqDateInput)
	dateInput := ""
	dateFormat := ""
	if reqDateInput.DateInput != "" {
		dateTemp := strings.Split(reqDateInput.DateInput, ".")
		dateInput = dateTemp[2] + "-" + dateTemp[1] + "-" + dateTemp[0]
		dateFormat = reqDateInput.DateInput
	} else {
		dateTempFormat := strings.Split(utils.ConvertDateToString(time.Now(), utils.TYPE_DATE_RFC3339), "-")
		dateFormat = dateTempFormat[2] + "." + dateTempFormat[1] + "." + dateTempFormat[0]
	}
	dataInvCancel := e.mapEmailRepo.GetDataInvCancel(dateInput)

	log.Info("Size data invoice cancel : ", len(dataInvCancel))
	log.Info("process email invoice no fp, receipt : ", os.Getenv("EMAIL_TO_CANCEL_INV"))
	log.Info("process email invoice no fp, CC : ", os.Getenv("EMAIL_CC_CANCEL_INV"))
	log.Info("process email invoice sender : ", os.Getenv("EMAIL_SENDER"))
	log.Info("process email, subject : ", os.Getenv("SUBJECT_EMAIL_INV_CANCEL") + " " + dateFormat)
	if len(dataInvCancel) > 0 {
		bodyEmail := utils.TemplateEmail(dataInvCancel, utils.TYPE_EMAIL_INV_CANCEL)
		if bodyEmail != "" {
			emailNotif := models.EmailNotificationModels{
				//Recipient:    os.Getenv("EMAIL_TO"),
				//Cc:           os.Getenv("EMAIL_CC"),
				Recipient:    os.Getenv("EMAIL_TO_CANCEL_INV"),
				Cc:           os.Getenv("EMAIL_CC_CANCEL_INV"),
				Sender:       os.Getenv("EMAIL_SENDER"),
				Subject:      os.Getenv("SUBJECT_EMAIL_INV_CANCEL") + " " + dateFormat,
				LogStatus:    0,
				Body:         b64.StdEncoding.EncodeToString([]byte(bodyEmail)),
				SendDateTime: time.Now(),
			}

			dataEmail := e.mapEmailRepo.SaveEmailNotif(emailNotif)
			if dataEmail.Id > 0 {
				respEmail, errSendEmail := utils.EmailNotificationService(dataEmail, int(dataEmail.Id), []string{})
				log.Info("Send email cancel with response: ", respEmail)
				if errSendEmail != nil && respEmail.ResponseCode != 200 {
					log.Error("Error send email cancel : ", errSendEmail.Error())
				} else {
					for i := 0; i < len(dataInvCancel); i++ {
						intVar, errorParser := strconv.Atoi(fmt.Sprintf("%v", dataInvCancel[i]["id_inv_cancel"]))

						if errorParser != nil {
							log.Info("update data compare, id his email : ", errorParser.Error())
						}
						errUpdate := e.mapEmailRepo.UpdateStatusEmail(dataEmail.Id, intVar)
						if errUpdate != nil {
							log.Info("update data compare, id his email : ", errUpdate.Error())
						}

					}

				}

			}
		} else {
			log.Error("Email notif invoice cancel not send :", "Error template email")
		}

	} else {
		log.Info("Email notif invoice cancel not send : ", "invoice cancel not found")
	}

}
/*
	Method email invoice belom terkirim ke customer
*/
func (e emailController) EmailInvNotSend(ctx *gin.Context) {
	var reqDateInput dto.ReqInputDate
	_ = ctx.ShouldBindJSON(&reqDateInput)
	dateInput := ""
	dateFormat := ""
	if reqDateInput.DateInput != "" {
		dateTemp := strings.Split(reqDateInput.DateInput, ".")
		dateInput = dateTemp[2] + "-" + dateTemp[1] + "-" + dateTemp[0]
		dateFormat = reqDateInput.DateInput
	} else {
		dateTempFormat := strings.Split(utils.ConvertDateToString(time.Now(), utils.TYPE_DATE_RFC3339), "-")
		dateFormat = dateTempFormat[2] + "." + dateTempFormat[1] + "." + dateTempFormat[0]
	}
	dataInv := e.mapEmailRepo.GetDataInvNotSend(dateInput)
	log.Info("Size data invoice not send : ", len(dataInv))
	log.Info("process email invoice no send, receipt : ", os.Getenv("EMAIL_TO_NOT_SEND_INV"))
	log.Info("process email invoice no send, CC : ", os.Getenv("EMAIL_CC_NOT_SEND_INV"))
	log.Info("process email invoice no send, sender : ", os.Getenv("EMAIL_SENDER"))
	log.Info("process email, invoice no send,subject : ", os.Getenv("SUBJECT_EMAIL_INV_NOT_SEND") + " " + dateFormat)
	if len(dataInv) > 0 {
		bodyEmail := utils.TemplateEmail(dataInv, utils.TYPE_EMAIL_INV_NOT_SEND)
		if bodyEmail != "" {
			emailNotif := models.EmailNotificationModels{
				//Recipient:    os.Getenv("EMAIL_TO"),
				//Cc:           os.Getenv("EMAIL_CC"),
				Recipient: os.Getenv("EMAIL_TO_NOT_SEND_INV"),
				Cc: os.Getenv("EMAIL_CC_NOT_SEND_INV"),
				Sender:       os.Getenv("EMAIL_SENDER"),
				Subject:      os.Getenv("SUBJECT_EMAIL_INV_NOT_SEND") + " " + dateFormat,
				LogStatus:    0,
				Body:         b64.StdEncoding.EncodeToString([]byte(bodyEmail)),
				SendDateTime: time.Now(),
			}

			dataEmail := e.mapEmailRepo.SaveEmailNotif(emailNotif)
			if dataEmail.Id > 0 {
				respEmail, errSendEmail := utils.EmailNotificationService(dataEmail, int(dataEmail.Id), []string{})
				log.Info("Send email cancel with response: ", respEmail)
				if errSendEmail != nil && respEmail.ResponseCode != 200 {
					log.Error("Error send email cancel : ", errSendEmail.Error())
				} else {
					for i := 0; i < len(dataInv); i++ {
						intVar, errorParser := strconv.Atoi(fmt.Sprintf("%v", dataInv[i]["id_compare"]))

						if errorParser != nil {
							log.Info("update data compare, id his email : ", errorParser.Error())
						}
						errUpdate := e.mapEmailRepo.UpdateStatusEmail(dataEmail.Id, intVar)
						if errUpdate != nil {
							log.Info("update data compare, id his email : ", errUpdate.Error())
						}

					}
				}
			}
		} else {
			log.Error("Email notif invoice not in zsd081 not send : ", "Error template email")
		}

	} else {
		log.Info("Email notif invoice not in zsd081 not send : ", "invoice not found")
	}

}

func InstanceEmailController(db repository.EmailRepository) EmailController {
	return &emailController{
		mapEmailRepo: db,
	}

}
