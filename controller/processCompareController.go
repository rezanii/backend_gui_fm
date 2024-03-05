package controller

import (
	"backend_gui/dto"
	"backend_gui/models"
	"backend_gui/repository"
	"backend_gui/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type ProcessCompareController interface {
	CompareZsd001nZv60(ctx *gin.Context)
	CompareInvoiceCancel(ctx *gin.Context)
	CompareFpCancel(ctx *gin.Context)
	InsertZsd001nToCompare(ctx *gin.Context)
	CheckInvoiceAlreadySend(ctx *gin.Context)
	CheckFpAlreadySend(ctx *gin.Context)
	GetCollectionInv(ctx *gin.Context)
}

type processCompareController struct {
	mapProcessCompareRepo repository.ProcessCompareRepository
}

/*
	Method get data collection inv
*/
func (p processCompareController) GetCollectionInv(ctx *gin.Context) {
	username, errSession := utils.GetSession(ctx, "username")
	if errSession != nil {
		log.Error("collection invoice , error :", errSession.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}
	userFpm := p.mapProcessCompareRepo.CheckUserAlreadyExist(fmt.Sprintf("%s", username))
	if userFpm.RoleUser == "COLLECTION_CLUSTER" || userFpm.RoleUser == "COLLECTION_NON_CLUSTER" || userFpm.RoleUser == "COLLECTION_ENTERPRISE" || userFpm.RoleUser == "ADMIN" || userFpm.RoleUser == "FINANCE_COLLECTION"{
		var reqCollectionInv dto.ReqCollectionInvDto
		err := ctx.ShouldBindJSON(&reqCollectionInv)
		if err != nil{
			log.Error("collection invoice, username : "+ fmt.Sprintf("%s", username)+", error", err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse(err.Error(), utils.ERR_VALIDATE_DATA))
			return
		}
		var dc []string
		if reqCollectionInv.TypeDc == "enterprise"{
			dc = append(dc,"20")
			dc = append(dc,"35")
		}else if reqCollectionInv.TypeDc == "non_cluster"{
			dc = append(dc, "30")
			dc = append(dc,"40")
			dc = append(dc,"45")
		}else if reqCollectionInv.TypeDc == "cluster"{
			dc = append(dc,"10")
			dc = append(dc, "50")
			dc = append(dc, "51")
		}

		start := (reqCollectionInv.Page * reqCollectionInv.MaxDataDisplay) - reqCollectionInv.MaxDataDisplay
		dataEnterprise := dto.ResPaginationDto{
			BaseUrl: os.Getenv("BASE_URL_FILE_FAKTUR_PAJAK")+"file-fpm",
			TotalData: p.mapProcessCompareRepo.GetTotalCollectionInv(reqCollectionInv,dc),
			Record:    p.mapProcessCompareRepo.GetRecordCollectionInv(start, reqCollectionInv.MaxDataDisplay,reqCollectionInv,dc),
		}
		ctx.JSON(http.StatusOK, dto.SuccessResponse(dataEnterprise, "", utils.SUCCESS_CODE))
	}else{
		log.Error("collection invoice, username",fmt.Sprintf("%s",username)+", error :", " user have not access this service")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}


}
/*
	Method process cek fp send to customer with StatusSendFp = yes
*/
func (p processCompareController) CheckFpAlreadySend(ctx *gin.Context) {
	dataFp := p.mapProcessCompareRepo.GetCompareFpHisZsd081()
	log.Info("Size data flag FP in zsd081 e-faktur : ", len(dataFp))
	msg := ""
	countSuccess := 0
	countFailed := 0
	anySuccess := false
	for n := 0; n < len(dataFp); n++ {
		dataFp[n].StatusSendFp = "yes"
		dataUptHisFp := models.FakturModels{
			ID:             dataFp[n].ID,
			InvNumber:      dataFp[n].InvNumber,
			InvDc:          dataFp[n].InvDc,
			Faktur:         dataFp[n].Faktur,
			StatusSendFp:   dataFp[n].StatusSendFp,
			Status:         dataFp[n].Status,
			EmailReceiptFp: dataFp[n].EmailReceiptFp,
			SendDateFp:     dataFp[n].SendDateFp,
			SendTimeFp:     dataFp[n].SendTimeFp,
			FakturCreated:  dataFp[n].FakturCreated,
		}
		errUpdate := p.mapProcessCompareRepo.UpdateFpAlreadySend(dataUptHisFp)
		if errUpdate != nil {
			msg += dataUptHisFp.Faktur + " failed flag"
			countFailed++
		} else {
			anySuccess = true
			countSuccess++
			msg += dataUptHisFp.Faktur + " success flag"
		}

	}
	sizeConvert := strconv.Itoa(len(dataFp))
	if len(dataFp) > 0 {
		log.Info("Process flag faktur : ", msg)
		log.Info("Process flag faktur success size : ", countSuccess)
		log.Info("Process flag faktur failed size : ", countFailed)
		if anySuccess {
			ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "Success flag faktur send to cust", utils.SUCCESS_CODE))
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("Failed flag faktur send to cust", utils.ERR_GLOBAL))
		}

	} else {
		ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "No process flag faktur, size data : "+sizeConvert, utils.SUCCESS_CODE))
	}

}
/*
	Method for check invoice send to customer flag sts_send_inv = yes
*/
func (p processCompareController) CheckInvoiceAlreadySend(ctx *gin.Context) {
	var reqDateInput dto.ReqInputDate
	_ = ctx.ShouldBindJSON(&reqDateInput)
	dateGetData := ""

	if reqDateInput.DateInput != "" {
		dateTemp := strings.Split(reqDateInput.DateInput, ".")
		dateInput := dateTemp[2] + "-" + dateTemp[1] + "-" + dateTemp[0]
		dateConvert, errDateParse := time.Parse(utils.TYPE_DATE_RFC3339, dateInput)
		if errDateParse != nil {
			log.Error("Error parsing date input process flaq invoice success send to customer: ", errDateParse.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("Error Parsing date input", utils.ERR_GLOBAL))
			return
		}
		dateSubstracOne := dateConvert.AddDate(0, 0, -1)
		dateConvertSts := utils.ConvertDateToString(dateSubstracOne, utils.TYPE_DATE_RFC3339)
		dateSplit := strings.Split(dateConvertSts, "-")
		dateGetData = dateSplit[2] + "." + dateSplit[1] + "." + dateSplit[0]
	}

	dataInv := p.mapProcessCompareRepo.GetCompareInvZsd081(dateGetData)

	log.Info("Process flag invoice success send to customer, size data : ", len(dataInv))
	AnyErrorUpdate := false
	countSuccess := 0
	countFailed := 0
	msg := ""
	for i := 0; i < len(dataInv); i++ {
		if dataInv[i].Id > 0 {
			dataInv[i].StsSendInv = "yes"
			dtaUpdate := models.DetailCompareModels{
				Id:                     dataInv[i].Id,
				BillingDocumentZsd001n: dataInv[i].BillingDocumentZsd001n,
				BillingDateZsd001n:     dataInv[i].BillingDateZsd001n,
				BillingTypeZsd001n:     dataInv[i].BillingTypeZsd001n,
				DcZsd001n:              dataInv[i].DcZsd001n,
				SlorZsd001n:            dataInv[i].SlorZsd001n,
				PayerZsd001n:           dataInv[i].PayerZsd001n,
				PayerNameZsd001n:       dataInv[i].PayerNameZsd001n,
				MaterialNumberZsd001n:  dataInv[i].MaterialNumberZsd001n,
				MaterialDescZsd001n:    dataInv[i].MaterialDescZsd001n,
				CreateOnZsd001n:        dataInv[i].CreateOnZsd001n,
				BillingDocCancel:       dataInv[i].BillingDocCancel,
				StsCancelInv:           dataInv[i].StsCancelInv,
				StsEmailInvCancel:      dataInv[i].StsEmailInvCancel,
				StsSendInv:             dataInv[i].StsSendInv,
				BillingNumberZv60:      dataInv[i].BillingNumberZv60,
				BillingDateZv60:        dataInv[i].BillingDateZv60,
				FpNumberZv60:           dataInv[i].FpNumberZv60,
				FpCreatedDateZv60:      dataInv[i].FpCreatedDateZv60,
				PayerZv60:              dataInv[i].PayerZv60,
				NameZv60:               dataInv[i].NameZv60,
				NpwpZv60:               dataInv[i].NpwpZv60,
				MaterialZv60:           dataInv[i].MaterialZv60,
				StsCompare:             dataInv[i].StsCompare,
				IdHistoryFpm:           dataInv[i].IdHistoryFpm,
				StsEmailCompare:        dataInv[i].StsEmailCompare,
				EmailReceiptInv:        dataInv[i].EmailReceiptInv,
				SendDateInv:            dataInv[i].SendDateInv,
				SendTimeInv:            dataInv[i].SendTimeInv,
				IdHisEmail:             dataInv[i].IdHisEmail,
			}
			errUpdateData := p.mapProcessCompareRepo.UpdateInvAlreadySend(dtaUpdate)
			if errUpdateData != nil {
				countFailed++
				msg += dtaUpdate.BillingDocumentZsd001n + " failed flag, "
				AnyErrorUpdate = true
			} else {
				msg += dtaUpdate.BillingDocumentZsd001n + " Success flag, "
				countSuccess++
			}
		}
	}
	sizeConvert := strconv.Itoa(len(dataInv))
	if len(dataInv) > 0 {
		if AnyErrorUpdate {
			log.Error("Any Error update flag invoice success send to customer")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("Any Error update flag invoice success send to customer", utils.ERR_GLOBAL))
			return
		}

		ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "Success update flag invoice success send to customer, size : "+sizeConvert, utils.SUCCESS_CODE))
	} else {
		ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "No process update flag invoice success send to customer, size : "+sizeConvert, utils.SUCCESS_CODE))
	}

}
/*
	Method for copy data zsd001n to sf_fpm_compare
*/
func (p processCompareController) InsertZsd001nToCompare(ctx *gin.Context) {
	var reqCopyData dto.ReqInputDate
	_ = ctx.ShouldBindJSON(&reqCopyData)
	dateZsd001n := ""
	anyFailedInsert := false
	//if reqCopyData.DateInput != "" {
	//	dateTemp := strings.Split(reqCopyData.DateInput, ".")
	//	dateInput = dateTemp[2] + dateTemp[1] + dateTemp[0]
	//}
	if reqCopyData.DateInput != "" {
		dateTemp := strings.Split(reqCopyData.DateInput, ".")
		dateInput := dateTemp[2] + "-" + dateTemp[1] + "-" + dateTemp[0]
		dateConvert, errDateParse := time.Parse(utils.TYPE_DATE_RFC3339, dateInput)
		if errDateParse != nil {
			log.Error("Error parsing date input process flaq invoice success send to customer: ", errDateParse.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("Error Parsing date input", utils.ERR_GLOBAL))
			return
		}
		dateSubstracOne := dateConvert.AddDate(0, 0, -1)
		dateConvertSts := utils.ConvertDateToString(dateSubstracOne, utils.TYPE_DATE_RFC3339)
		dateSplit := strings.Split(dateConvertSts, "-")
		dateZsd001n = dateSplit[0]  + dateSplit[1]  + dateSplit[2]
	}

	data := p.mapProcessCompareRepo.GetDataZsd001nByCreateOn(dateZsd001n)

	log.Info("Process batch copy file size : ", len(data))
	if len(data) > 0 {
		dtaMapCompare := []models.DetailCompareModels{}
		for i := 0; i < len(data); i++ {
			stsEmailCompare := "ready_to_send"
			stsCompare := "not_match"
			stsCancelInv := "no"
			stsSendInv := "no"

			recordsMap := models.DetailCompareModels{
				BillingDocumentZsd001n: data[i].BillingDocumentZsd001n,
				BillingTypeZsd001n:     data[i].BillingTypeZsd001n,
				BillingDateZsd001n:     data[i].BillingDateZsd001n,
				DcZsd001n:              data[i].DcZsd001n,
				PayerZsd001n:           data[i].PayerZsd001n,
				PayerNameZsd001n:       data[i].PayerNameZsd001n,
				SlorZsd001n:            data[i].SlorZsd001n,
				DistrictZsd001n:        data[i].DistrictZsd001n,
				MaterialNumberZsd001n:  data[i].MaterialNumberZsd001n,
				MaterialDescZsd001n:    data[i].MaterialDescZsd001n,
				CreateOnZsd001n:        data[i].CreateOnZsd001n,
				StsCancelInv:           stsCancelInv,
				StsSendInv:             stsSendInv,
				StsCompare:             stsCompare,
				StsEmailCompare:        stsEmailCompare,
			}
			dtaMapCompare = append(dtaMapCompare, recordsMap)
		}
		limit := 500
		recordProcess := 0
		if len(dtaMapCompare) > limit {
			dataInsert := []models.DetailCompareModels{}
			for m := 0; m < len(dtaMapCompare); m++ {
				dataInsert = append(dataInsert, dtaMapCompare[m])
				if (len(dtaMapCompare) - recordProcess) >= limit {
					if len(dataInsert) == limit {
						log.Info("Process insert batch size : ", limit)
						errInsertBatchLimit := p.mapProcessCompareRepo.CreateDataCompare(dataInsert)
						if errInsertBatchLimit != nil {
							anyFailedInsert = true
						}
						dataInsert = []models.DetailCompareModels{}
						recordProcess += limit
					}
				} else {
					if m == (len(dtaMapCompare) - 1) {
						log.Info("Process insert batch size : ", len(dataInsert))
						errInsertSisaData := p.mapProcessCompareRepo.CreateDataCompare(dataInsert)
						if errInsertSisaData != nil {
							anyFailedInsert = true
						}
					}
				}

			}

		} else {
			log.Info("Process insert bacth size : ", len(dtaMapCompare))
			errInsertDataLessLimit := p.mapProcessCompareRepo.CreateDataCompare(dtaMapCompare)
			if errInsertDataLessLimit != nil {
				anyFailedInsert = true
			}
		}
	}
	if anyFailedInsert {
		log.Error("Any Error insert bacth Zsd001n to compare")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("Any Error insert bacth Zsd001n to compare", utils.ERR_GLOBAL))
		return
	}
	sizeConvert := strconv.Itoa(len(data))
	ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "Success insert batch size : "+sizeConvert, utils.SUCCESS_CODE))

}
/*
	Method proses Fp cancel status = 0 and update_at fill
    Fp not cancel will be insert new row with status = 1
*/
func (p processCompareController) CompareFpCancel(ctx *gin.Context) {
	utils.CheckDayLog()
	var reqFpCancel dto.ReqInputDate
	_ = ctx.ShouldBindJSON(&reqFpCancel)
	dateGetData := ""
	if reqFpCancel.DateInput != "" {
		dateTemp := strings.Split(reqFpCancel.DateInput, ".")
		dateInput := dateTemp[2] + "-" + dateTemp[1] + "-" + dateTemp[0]
		dateConvert, errDateParse := time.Parse(utils.TYPE_DATE_RFC3339, dateInput)
		if errDateParse != nil {
			log.Error("Error Parsing date input flag FP cancel : ", errDateParse.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("Error Parsing date input", utils.ERR_GLOBAL))
			return
		}
		dateSubstracOne := dateConvert.AddDate(0, 0, -1)
		dateConvertSts := utils.ConvertDateToString(dateSubstracOne, utils.TYPE_DATE_RFC3339)
		dateSplit := strings.Split(dateConvertSts, "-")
		dateGetData = dateSplit[2] + "." + dateSplit[1] + "." + dateSplit[0]
	}

	dataFp := p.mapProcessCompareRepo.GetFpByCreated(dateGetData)
	log.Info("size of prosess data fp cancel : ", len(dataFp))
	AnyFailedUpdateOrSave := false
	msgUpdate := ""
	msgInsert := ""
	for i := 0; i < len(dataFp); i++ {
		checkNewFpExist := p.mapProcessCompareRepo.GetHistoryFaktur(dataFp[i].NewFp)
		if checkNewFpExist.ID > 0 {
			continue
		}
		checkOldFakturExist := p.mapProcessCompareRepo.GetHistoryFaktur(dataFp[i].FpOld)

		if checkOldFakturExist.ID > 0 {
			//ketika fp cancel fp comparenya kolom fp di update jika dia ada yang fp lagi ke cancel ke 2 dst fp yang old
			//yang di ambil bukan yang pertama melainkan yang terupdate
			if checkOldFakturExist.Faktur == dataFp[i].NewFp {
				continue
			}
			errUpdateHis := p.mapProcessCompareRepo.UpdateHistoryFaktur(checkOldFakturExist.InvNumber)
			if errUpdateHis != nil {
				AnyFailedUpdateOrSave = true
				msgUpdate += checkOldFakturExist.InvNumber + " Failed non active faktur, "
				continue
			} else {
				msgUpdate += checkOldFakturExist.InvNumber + " Success non active faktur, "
			}
		}
		dataUpdate := models.FakturModels{
			InvNumber:     dataFp[i].InvNumber,
			InvDc:         dataFp[i].InvChannel,
			Faktur:        dataFp[i].NewFp,
			FakturCreated: dataFp[i].NewFpCreated,
			StatusSendFp:  "no",
			Status:        1,
		}
		errSave := p.mapProcessCompareRepo.SaveHistoryFaktur(dataUpdate)
		if errSave != nil {
			AnyFailedUpdateOrSave = true
			msgInsert += dataUpdate.Faktur + " failed insert, "
		} else {
			msgInsert += dataUpdate.Faktur + " success insert, "
		}
	}
	if len(dataFp) > 0 {
		log.Info("process update old faktur : ", msgUpdate)
		log.Info("process insert new faktur : ", msgInsert)
		if AnyFailedUpdateOrSave {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("Any Error insert or update data history faktur", utils.ERR_GLOBAL))
			return
		}
		ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "Success insert or update data history faktur", utils.SUCCESS_CODE))
	} else {
		sizeConvert := strconv.Itoa(len(dataFp))
		ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "No process insert or update data history faktur, size : "+sizeConvert, utils.SUCCESS_CODE))
	}

}
/*
 Method for cek invoice cancel
*/
func (p processCompareController) CompareInvoiceCancel(ctx *gin.Context) {
	utils.CheckDayLog()
	var reqCancelInvoice dto.ReqInputDate
	_ = ctx.ShouldBindJSON(&reqCancelInvoice)
	dateGetData := ""
	//if reqCancelInvoice.DateInput != "" {
	//	dateTemp := strings.Split(reqCancelInvoice.DateInput, ".")
	//	dateGetData = dateTemp[2] + dateTemp[1] + dateTemp[0]
	//
	//}
	if reqCancelInvoice.DateInput != "" {
		dateTemp := strings.Split(reqCancelInvoice.DateInput, ".")
		dateInput := dateTemp[2] + "-" + dateTemp[1] + "-" + dateTemp[0]
		dateConvert, errDateParse := time.Parse(utils.TYPE_DATE_RFC3339, dateInput)
		if errDateParse != nil {
			log.Error("Error parsing date input process flaq invoice success send to customer: ", errDateParse.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("Error Parsing date input", utils.ERR_GLOBAL))
			return
		}
		dateSubstracOne := dateConvert.AddDate(0, 0, -1)
		dateConvertSts := utils.ConvertDateToString(dateSubstracOne, utils.TYPE_DATE_RFC3339)
		dateSplit := strings.Split(dateConvertSts, "-")
		dateGetData = dateSplit[0]  + dateSplit[1]  + dateSplit[2]
	}

	compareInvoiceCancel := p.mapProcessCompareRepo.CompareInvoiceCancel(dateGetData)
	log.Info("size of invoice cancel : ", len(compareInvoiceCancel))

	if len(compareInvoiceCancel) > 0 {
		message := ""
		anySuccessUpdate := false
		countSuccess := 0
		countFailed := 0
		for i := 0; i < len(compareInvoiceCancel); i++ {
			if compareInvoiceCancel[i].StsSendInv == "yes" {
				if compareInvoiceCancel[i].DcZsd001n == "10" {
					compareInvoiceCancel[i].StsEmailInvCancel = "ready_to_send"
					compareInvoiceCancel[i].StsCancelInv = "yes"
				} else {
					compareInvoiceCancel[i].StsEmailInvCancel = "not_process"
					compareInvoiceCancel[i].StsCancelInv = "yes"
				}
			} else {
				compareInvoiceCancel[i].StsCancelInv = "yes"
				compareInvoiceCancel[i].StsEmailInvCancel = "not_process"
			}
			data := models.DetailCompareModels{
				Id:                     compareInvoiceCancel[i].Id,
				BillingDocumentZsd001n: compareInvoiceCancel[i].BillingDocumentZsd001n,
				BillingTypeZsd001n:     compareInvoiceCancel[i].BillingTypeZsd001n,
				BillingDateZsd001n:     compareInvoiceCancel[i].BillingDateZsd001n,
				DcZsd001n:              compareInvoiceCancel[i].DcZsd001n,
				SlorZsd001n:            compareInvoiceCancel[i].SlorZsd001n,
				CreateOnZsd001n:        compareInvoiceCancel[i].CreateOnZsd001n,
				PayerZsd001n:           compareInvoiceCancel[i].PayerZsd001n,
				PayerNameZsd001n:       compareInvoiceCancel[i].PayerNameZsd001n,
				MaterialNumberZsd001n:  compareInvoiceCancel[i].MaterialNumberZsd001n,
				MaterialDescZsd001n:    compareInvoiceCancel[i].MaterialDescZsd001n,
				BillingDocCancel:       compareInvoiceCancel[i].InvCancel,
				StsCancelInv:           compareInvoiceCancel[i].StsCancelInv,
				StsEmailInvCancel:      compareInvoiceCancel[i].StsEmailInvCancel,
				StsSendInv:             compareInvoiceCancel[i].StsSendInv,
				BillingNumberZv60:      compareInvoiceCancel[i].BillingNumberZv60,
				BillingDateZv60:        compareInvoiceCancel[i].BillingDateZv60,
				FpNumberZv60:           compareInvoiceCancel[i].FpNumberZv60,
				FpCreatedDateZv60:      compareInvoiceCancel[i].FpCreatedDateZv60,
				PayerZv60:              compareInvoiceCancel[i].PayerZv60,
				NameZv60:               compareInvoiceCancel[i].NameZv60,
				NpwpZv60:               compareInvoiceCancel[i].NpwpZv60,
				MaterialZv60:           compareInvoiceCancel[i].MaterialZv60,
				StsCompare:             compareInvoiceCancel[i].StsCompare,
				IdHistoryFpm:           compareInvoiceCancel[i].IdHistoryFpm,
				StsEmailCompare:        compareInvoiceCancel[i].StsEmailCompare,
				EmailReceiptInv:        compareInvoiceCancel[i].EmailReceiptInv,
				SendDateInv:            compareInvoiceCancel[i].SendDateInv,
				SendTimeInv:            compareInvoiceCancel[i].SendTimeInv,
				IdHisEmail:             compareInvoiceCancel[i].IdHisEmail,
			}
			result := p.mapProcessCompareRepo.UpdateInvoiceCancel(data)
			if result != nil {
				countFailed++
				message += compareInvoiceCancel[i].BillingDocumentZsd001n + " Failed update, "
			} else {
				anySuccessUpdate = true
				countSuccess++
				message += compareInvoiceCancel[i].BillingDocumentZsd001n + " Success update, "
			}
		}
		log.Info("Size Success flag invoice cancel : ", countSuccess)
		log.Info("Size Failed flag invoice cancel : ", countFailed)
		if anySuccessUpdate {
			log.Info("update invoice cancel : ", message)
			sizeConvert := strconv.Itoa(countSuccess)
			ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "Success flaq invoice cancel : "+sizeConvert, utils.SUCCESS_CODE))
		} else {
			log.Info("update invoice cancel : ", "all invoice cancel failed update")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("all invoice cancel failed update", utils.ERR_GLOBAL))
			return
		}
	} else {
		log.Info("update invoice cancel : data inv cancel not found")
		ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "Invoice cancel not found", utils.SUCCESS_CODE))
	}

}
/*
 Method for compare data at table sf_fpm_compare vs zv60
*/
func (p processCompareController) CompareZsd001nZv60(ctx *gin.Context) {
	utils.CheckDayLog()
	var reqCompareData dto.ReqInputDate
	_ = ctx.ShouldBindJSON(&reqCompareData)
	dateGetData := ""
	if reqCompareData.DateInput != "" {
		dateTemp := strings.Split(reqCompareData.DateInput, ".")
		dateInput := dateTemp[2] + "-" + dateTemp[1] + "-" + dateTemp[0]
		dateConvert, errDateParse := time.Parse(utils.TYPE_DATE_RFC3339, dateInput)
		if errDateParse != nil {
			log.Error("Error Parsing date input compare zv60 vs tbl compare : ", errDateParse.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("Error Parsing date input", utils.ERR_GLOBAL))
			return
		}
		dateSubstracOne := dateConvert.AddDate(0, 0, -1)
		dateConvertSts := utils.ConvertDateToString(dateSubstracOne, utils.TYPE_DATE_RFC3339)
		dateSplit := strings.Split(dateConvertSts, "-")
		dateGetData = dateSplit[2] + "." + dateSplit[1] + "." + dateSplit[0]
	}

	compareZsd001nZv60 := p.mapProcessCompareRepo.CompareZsd001nZv60(dateGetData)
	log.Info("process compare zv60 vs tbl compare size match: ", len(compareZsd001nZv60))

	anySuccess := false
	message := ""
	countSuccess := 0
	countFailed := 0
	for i := 0; i < len(compareZsd001nZv60); i++ {
		data := models.DetailCompareModels{
			Id:                     compareZsd001nZv60[i].Id,
			BillingDocumentZsd001n: compareZsd001nZv60[i].BillingDocumentZsd001n,
			BillingTypeZsd001n:     compareZsd001nZv60[i].BillingTypeZsd001n,
			BillingDateZsd001n:     compareZsd001nZv60[i].BillingDateZsd001n,
			DcZsd001n:              compareZsd001nZv60[i].DcZsd001n,
			PayerZsd001n:           compareZsd001nZv60[i].DcZsd001n,
			PayerNameZsd001n:       compareZsd001nZv60[i].PayerNameZsd001n,
			SlorZsd001n:            compareZsd001nZv60[i].SlorZsd001n,
			DistrictZsd001n:        compareZsd001nZv60[i].DistrictZsd001n,
			MaterialNumberZsd001n:  compareZsd001nZv60[i].MaterialNumberZsd001n,
			MaterialDescZsd001n:    compareZsd001nZv60[i].MaterialDescZsd001n,
			CreateOnZsd001n:        compareZsd001nZv60[i].CreateOnZsd001n,
			BillingDocCancel:       compareZsd001nZv60[i].BillingDocCancel,
			StsCancelInv:           compareZsd001nZv60[i].StsCancelInv,
			StsEmailInvCancel:      compareZsd001nZv60[i].StsEmailInvCancel,
			StsSendInv:             compareZsd001nZv60[i].StsSendInv,
			BillingNumberZv60:      compareZsd001nZv60[i].BillingNumberZv60,
			BillingDateZv60:        compareZsd001nZv60[i].BillingDateZv60,
			FpNumberZv60:           compareZsd001nZv60[i].FpNumberZv60,
			FpCreatedDateZv60:      compareZsd001nZv60[i].FpCreatedDateZv60,
			PayerZv60:              compareZsd001nZv60[i].PayerZv60,
			NameZv60:               compareZsd001nZv60[i].NameZv60,
			NpwpZv60:               compareZsd001nZv60[i].NpwpZv60,
			MaterialZv60:           compareZsd001nZv60[i].MaterialZv60,
			StsCompare:             "match",
			IdHistoryFpm:           compareZsd001nZv60[i].IdHistoryFpm,
			StsEmailCompare:        "not_process",
			EmailReceiptInv:        compareZsd001nZv60[i].EmailReceiptInv,
			SendDateInv:            compareZsd001nZv60[i].SendDateInv,
			SendTimeInv:            compareZsd001nZv60[i].SendTimeInv,
			IdHisEmail:             compareZsd001nZv60[i].IdHisEmail,
		}

		result := p.mapProcessCompareRepo.UpdateCompareZsd001nZv60(data)

		if result == nil {
			anySuccess = true
			countSuccess++
			message += compareZsd001nZv60[i].BillingDocumentZsd001n + " Success, "
		} else {
			countFailed++
			message += compareZsd001nZv60[i].BillingDocumentZsd001n + " Failed, "
		}
	}
	log.Info("Process compare zv60 vs tbl compare success update match : ", countSuccess)
	log.Info("Process compare zv60 vs tbl compare failed update match : ", countFailed)
	if !anySuccess {
		log.Info("update data compare : ", "all invoice failed update")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("all invoice failed update", utils.ERR_GLOBAL))
		return
	} else {
		log.Info("update data compare : ", message)
		sizeConvert := strconv.Itoa(countSuccess)
		ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "Success update match, size data: "+sizeConvert, utils.SUCCESS_CODE))
	}

}

func InstanceProcessCompareController(mapProcessCompareRepo repository.ProcessCompareRepository) ProcessCompareController {
	return &processCompareController{
		mapProcessCompareRepo: mapProcessCompareRepo,
	}
}
