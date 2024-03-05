package controller

import (
	"backend_gui/dto"
	"backend_gui/models"
	"backend_gui/repository"
	"backend_gui/utils"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

type FpmController interface {
	MoveFpm(ctx *gin.Context)
	UpdateStatusFile(ctx *gin.Context)
}
type mapFpmController struct {
	mapFpmRepository repository.FpmRepository
}

/*
update status file faktur pajak after SAP processed
by irma 30/12/2021
*/
func (m mapFpmController) UpdateStatusFile(ctx *gin.Context) {
	data := m.mapFpmRepository.GetHistoryFpmBySts("finish_move_sap", "failed_send_to_cust")
	log.Info("Size of data for update status e-faktur : ", len(data))
	isAnyErrUpdateSts := false
	message := ""
	var stsEfaktur string = ""
	for i := 0; i < len(data); i++ {
		if data[i].Status == "failed_send_to_cust" && data[i].StatusSendFp == "no" {
			log.Info(data[i].NoFaktur , " not update, value same as the previous day")
			continue
		}

		if data[i].StatusSendFp == "no" {
			stsEfaktur = "failed_send_to_cust"
		} else {
			stsEfaktur = "success_send_to_cust"
		}
		dataUpdt := models.FpmHistoryFilePjk{
			ID:               data[i].Id,
			UploadBy:         data[i].UploadBy,
			JenisFakturPajak: data[i].JenisFakturPajak,
			FileNameUpload:   data[i].FileNameUpload,
			Url:              data[i].Url,
			NoFaktur:         data[i].NoFaktur,
			Status:           stsEfaktur,
		}
		errUpdt := m.mapFpmRepository.UpdateStsFpmHistory(dataUpdt)
		if errUpdt == nil {
			message += "Succes update send customer, faktur " + data[i].NoFaktur
			isAnyErrUpdateSts = true
		} else {
			message += "Failed update send customer, faktur " + data[i].NoFaktur + ", "
		}
	}
	if len(data) > 0 {
		if isAnyErrUpdateSts == false {
			log.Error("Process update data faktur : ","failed process update All data faktur")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("Failed process update data faktur", utils.ERR_GLOBAL))
		} else {
			log.Info("Process update data faktur  : ", message)
			ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "Finish process update data faktur", utils.SUCCESS_CODE))
		}
	} else {
		sizeConvert := strconv.Itoa(len(data))
		ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "No process update, size data : "+sizeConvert, utils.SUCCESS_CODE))
	}

}

/*
	fuction to move file faktur pajak base on dc and
    if dc = 10 and invoice status alreadysend to customer
    move SAP
*/
func (m mapFpmController) MoveFpm(ctx *gin.Context) {
	utils.CheckDayLog()
	dataFpm := m.mapFpmRepository.GetDataFp()
	log.Info("Process move file faktur , size of data : ", len(dataFpm))

	for i := 0; i < len(dataFpm); i++ {
		if (dataFpm[i].DcZsd001n != "10") {
			folderTo:=os.Getenv("ROOT_UPLOAD_FPM_PATH")+dataFpm[i].JenisFakturPajak+"/"+dataFpm[i].FileNameFp
			folderDest:= os.Getenv("ROOT_UPLOAD_FPM_PATH")+dataFpm[i].JenisFakturPajak+"/"
			if !utils.CheckChannel(dataFpm[i].DcZsd001n){
				folderDest+=os.Getenv("MOVE_OTHERS_CHANEL")+dataFpm[i].FileNameFp
				dataFpm[i].StsFile = "finish_move_" + os.Getenv("FOLDER_OTHER_DC")
				dataFpm[i].UrlFile = os.Getenv("BASE_URL_FILE_FAKTUR_PAJAK") + "file-fpm/" + dataFpm[i].JenisFakturPajak + "/" + os.Getenv("FOLDER_OTHER_DC") + "/" + dataFpm[i].FileNameFp

			}else{
				folderDest+= dataFpm[i].DcZsd001n+"/"+dataFpm[i].FileNameFp
				dataFpm[i].StsFile = "finish_move_" + dataFpm[i].DcZsd001n
				dataFpm[i].UrlFile = os.Getenv("BASE_URL_FILE_FAKTUR_PAJAK") + "file-fpm/" + dataFpm[i].JenisFakturPajak + "/" + dataFpm[i].DcZsd001n + "/" + dataFpm[i].FileNameFp
			}

			errMove := utils.MoveFile(folderTo,folderDest , false)
			if errMove != nil {
				log.Error(dataFpm[i].Faktur+"failed file move with error : ", errMove.Error())
				continue
			} else {
				log.Info(dataFpm[i].Faktur + " success move ")
			}

			errUpdateHis := m.mapFpmRepository.UpdateDataFpm(dataFpm[i])
			if errUpdateHis != nil {
				log.Error(dataFpm[i].FpNumberZv60 + " failed move ")
			}
		}

		if dataFpm[i].DcZsd001n == "10" {
			if _, errSearch := os.Stat(os.Getenv("ROOT_UPLOAD_FPM_PATH") + dataFpm[i].JenisFakturPajak + "/" + dataFpm[i].FileNameFp); errors.Is(errSearch, os.ErrNotExist) {
				log.Error(dataFpm[i].Faktur + " file upload not exist")
				continue
			}
			resultCopy := utils.CopyFile((os.Getenv("ROOT_UPLOAD_FPM_PATH") + dataFpm[i].JenisFakturPajak + "/" + dataFpm[i].FileNameFp), (os.Getenv("ROOT_UPLOAD_FPM_PATH") + os.Getenv("PATH_SAP") + dataFpm[i].FileNameFp), true)
			if resultCopy {
				log.Info("Success copy file fp : ", dataFpm[i].FileNameFp)
				dataFpm[i].StsFile = "finish_move_sap"
				dataFpm[i].UrlFile = os.Getenv("BASE_URL_FILE_FAKTUR_PAJAK") + "file-fpm/" + dataFpm[i].JenisFakturPajak + "/" + dataFpm[i].DcZsd001n + "/" + dataFpm[i].FileNameFp
				resultUpdate := m.mapFpmRepository.UpdateDataFpm(dataFpm[i])
				if resultUpdate == nil {
					errMoveChannel10 := utils.MoveFile(os.Getenv("ROOT_UPLOAD_FPM_PATH")+dataFpm[i].JenisFakturPajak+"/"+dataFpm[i].FileNameFp, os.Getenv("ROOT_UPLOAD_FPM_PATH")+dataFpm[i].JenisFakturPajak+"/"+dataFpm[i].DcZsd001n+"/"+dataFpm[i].FileNameFp, false)
					if errMoveChannel10 != nil {
						log.Error(dataFpm[i].Faktur + " failed move to folder channel 10")
					} else {
						log.Info("Success move file fp : ", dataFpm[i].FileNameFp)
					}
				} else {
					errRemove := os.Remove(os.Getenv("ROOT_UPLOAD_FPM_PATH") + os.Getenv("PATH_SAP") + dataFpm[i].FileNameFp)
					if errRemove != nil {
						log.Error("err remove file : ", dataFpm[i].FileNameFp, errRemove)
					}
				}
			} else {
				log.Error(dataFpm[i].Faktur + " failed file move to sap")
			}

		}
	}
	sizeConvert := strconv.Itoa(len(dataFpm))
	if len(dataFpm) > 0 {
		ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "Finish process move file faktur", utils.SUCCESS_CODE))
	} else {
		ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "No file move, size data : "+sizeConvert, utils.SUCCESS_CODE))
	}

}

func InstanceFpmController(repo repository.FpmRepository) FpmController {
	return &mapFpmController{
		mapFpmRepository: repo,
	}

}
