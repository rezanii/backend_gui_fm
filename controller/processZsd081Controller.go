package controller

import (
	"backend_gui/dto"
	"backend_gui/models"
	"backend_gui/repository"
	"backend_gui/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ProcessZsd081Controller interface {
	InsertDataZsd081(ctx *gin.Context)
}

type processZsd081Controller struct {
	mapProcessFileRepo repository.ProcessZsd081Repository
}

func InstanceProcessZsd081Controller(mapProcessFileRepo repository.ProcessZsd081Repository) ProcessZsd081Controller {
	return &processZsd081Controller{
		mapProcessFileRepo: mapProcessFileRepo,
	}
}
/*
fuction to process insert data file upload zsd081 to database
by irma 30/12/2021
*/
func (p processZsd081Controller) InsertDataZsd081(ctx *gin.Context) {
	utils.CheckDayLog()
	files, err := ioutil.ReadDir(os.Getenv("ROOT_UPLOAD_ZSD081_PATH"))
	if err != nil {
		log.Error("Error background process insert file zsd081 :", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("err open folder source_data", utils.ERR_GLOBAL))
		return

	}
	isAnyErrorProcessFile := false
	isAnyErrorProcessRecord := false
	countProcessFile := 0
	successInsert := 0
	failedInsert := 0
	totalOfRecord := 0

	for _, f := range files {
		if !f.IsDir() {
			//match, _ := regexp.MatchString("export_([0-9]{2}[0-9]{2}[0-9]{4})_invoicesendemail.xlsx", strings.ToLower(f.Name()))
			//if match {
				dtaFileProcess := p.mapProcessFileRepo.GetFileZsd081Process(f.Name())
				if dtaFileProcess.ID > 0 {
					log.Error("Error background process insert file zsd081 :" + f.Name() + " on process")
					isAnyErrorProcessFile = true
					continue
				}

				errSaveFlagFile := p.mapProcessFileRepo.FlagFileZsd081Process(f.Name())
				if errSaveFlagFile != nil {
					log.Info("Error background process insert file zsd081 :"+f.Name(), errSaveFlagFile)
					isAnyErrorProcessFile = true
					continue
				}

				xlsx, errReadlsx := excelize.OpenFile(os.Getenv("ROOT_UPLOAD_ZSD081_PATH") + f.Name())
				if errReadlsx != nil {
					log.Error("Error background process insert file zsd081 :", f.Name(), errReadlsx.Error())
					isAnyErrorProcessFile = true
					continue
				}
				sheet1Name := "Sheet1"

				rows, errGetRowXlsx := xlsx.GetRows(sheet1Name)
				if errGetRowXlsx != nil {
					log.Error("Error background process insert file zsd081 :", f.Name(), errGetRowXlsx.Error())
					isAnyErrorProcessFile = true
					continue
				}

				resultCheckHeader := utils.CheckHeaderFile(rows[0])
				if resultCheckHeader {
					countProcessFile += 1
					dataHistoryUpload := p.mapProcessFileRepo.GetHistoryZsd081FileByName(f.Name())
					if dataHistoryUpload.ID == 0 {
						dataSave := models.FpmHistoryFileZsd081{
							UploadBy:   "SFTP",
							FileName:   f.Name(),
							StatusFile: "uploaded",
						}
						dataHistoryUpload.ID = p.mapProcessFileRepo.SaveDataHistoryZsd081(dataSave)
						if dataHistoryUpload.ID == 0 {
							log.Error("Error background process insert file zsd081 :", f.Name(), "error insert history zsd081 to db")
							continue
						} else {
							dataHistoryUpload.UploadBy = dataSave.UploadBy
							dataHistoryUpload.FileName = f.Name()
							dataHistoryUpload.StatusFile = dataSave.StatusFile

						}
					}
					//isAnySuccessInsert := false
					for rowIdx, _ := range rows {

						if rowIdx > 0 {
							totalOfRecord += 1
							docNumber, errDocNumber := xlsx.GetCellValue("Sheet1", fmt.Sprintf("A%d", rowIdx+1))
							if errDocNumber != nil {
								log.Error("Error background process insert file zsd081 :", f.Name(), errDocNumber.Error())
								isAnyErrorProcessRecord = true
								continue
							}

							emailAddress, errEmailAddres := xlsx.GetCellValue("Sheet1", fmt.Sprintf("B%d", rowIdx+1))
							if errEmailAddres != nil {
								log.Error("Error background process insert file zsd081 :", f.Name(), errEmailAddres.Error())
								isAnyErrorProcessRecord = true
								continue
							}

							createOn, errCreateOn := xlsx.GetCellValue("Sheet1", fmt.Sprintf("C%d", rowIdx+1))

							if errCreateOn != nil {
								log.Error("Error background process insert file zsd081 :", f.Name(), errCreateOn.Error())
								isAnyErrorProcessRecord = true
								continue
							}
							flt, errParsetoFloat := strconv.ParseFloat(createOn, 8)
							if errParsetoFloat == nil {
								date, errParsing := excelize.ExcelDateToTime(flt, false)
								if errParsing != nil {
									log.Error("Error background process insert file zsd081 :", f.Name(), errParsing.Error())
									isAnyErrorProcessRecord = true
									continue
								}
								createOn = date.Format("2006-01-02")
							} else {
								//fmt.Println(fmt.Sprintf("Parsing date error: %s, index %d ----",f.Name(),rowIdx),"error parsing date :", errParsetoFloat)
								//continue
								temp := strings.Split(createOn, "-")
								if len(temp) < 3 {
									log.Error("Error background process insert file zsd081 :", f.Name(), "error parsing date wrong format")
									isAnyErrorProcessRecord = true
									continue
								}

								createOn = "20" + temp[2] + "-" + temp[1] + "-" + temp[0]

							}

							cellNumberTime := fmt.Sprintf("D%d", rowIdx+1)
							style, errStyle := xlsx.NewStyle(`{"number_format":21}`)
							if errStyle != nil {
								log.Error("Error background process insert file zsd081 :", f.Name(), errStyle.Error())
								isAnyErrorProcessRecord = true
								continue
							}
							errSetStyle := xlsx.SetCellStyle("Sheet1", cellNumberTime, cellNumberTime, style)
							if errSetStyle != nil {
								log.Error("Error background process insert file zsd081 :", f.Name(), errSetStyle.Error())
								isAnyErrorProcessRecord = true
								continue
							}
							timeInput, errParsingTime := xlsx.GetCellValue("Sheet1", cellNumberTime)
							if errParsingTime != nil {
								log.Error("Error background process insert file zsd081 :", f.Name(), errParsingTime.Error())
								isAnyErrorProcessRecord = true
								continue
							}

							createBy, errCreateBy := xlsx.GetCellValue("Sheet1", fmt.Sprintf("E%d", rowIdx+1))
							if errCreateBy != nil {
								log.Error("Error background process insert file zsd081 :", f.Name(), errCreateBy.Error())
								isAnyErrorProcessRecord = true
								continue
							}

							emailToOrCc, errEmailToOrCc := xlsx.GetCellValue("Sheet1", fmt.Sprintf("F%d", rowIdx+1))
							if errEmailToOrCc != nil {
								log.Error("Error background process insert file zsd081 :", f.Name(), errEmailToOrCc.Error())
								isAnyErrorProcessRecord = true
								continue
							}

							payer, errPayer := xlsx.GetCellValue("Sheet1", fmt.Sprintf("G%d", rowIdx+1))
							if errPayer != nil {
								log.Error("Error background process insert file zsd081 :", f.Name(), errPayer.Error())
								isAnyErrorProcessRecord = true
								continue
							}

							payerName, errPayerName := xlsx.GetCellValue("Sheet1", fmt.Sprintf("H%d", rowIdx+1))
							if errPayerName != nil {
								log.Error("Error background process insert file zsd081 :", f.Name(), errPayerName.Error())
								isAnyErrorProcessRecord = true
								continue
							}

							shipTo, errShipTo := xlsx.GetCellValue("Sheet1", fmt.Sprintf("I%d", rowIdx+1))
							if errShipTo != nil {
								log.Error("Error background process insert file zsd081 :", f.Name(), errShipTo.Error())
								isAnyErrorProcessRecord = true
								continue
							}

							shipToName, errShipToName := xlsx.GetCellValue("Sheet1", fmt.Sprintf("J%d", rowIdx+1))
							if errShipToName != nil {
								log.Error("Error background process insert file zsd081 :", f.Name(), errShipToName.Error())
								isAnyErrorProcessRecord = true
								continue
							}

							dataSaveZsd081 := models.FpmDetailInputZsd081{
								IDHistory:    dataHistoryUpload.ID,
								DocNumber:    docNumber,
								EmailAddress: emailAddress,
								CreateOn:     createOn,
								Time:         timeInput,
								CreateBy:     createBy,
								EmailToOrCc:  emailToOrCc,
								Payer:        payer,
								PayerName:    payerName,
								ShipTo:       shipTo,
								ShipToName:   shipToName,
							}
							errSaveData := p.mapProcessFileRepo.InsertDataZsd081(dataSaveZsd081)
							if errSaveData != nil {
								failedInsert += 1
								log.Error("Error background process insert file zsd081 :", f.Name(), errSaveData.Error())
								isAnyErrorProcessRecord = true
							} else {
								successInsert += 1
							}

						}

					}
					log.Info("Info total","Total record : " + fmt.Sprintf("%d", totalOfRecord) + ", Success insert : " + fmt.Sprintf("%d", successInsert) + ", Failed insert :" + fmt.Sprintf("%d", failedInsert))
					errUpdateStatus := p.mapProcessFileRepo.UpdateStatusZsd081(models.FpmHistoryFileZsd081{
						ID:         dataHistoryUpload.ID,
						UploadBy:   dataHistoryUpload.UploadBy,
						FileName:   dataHistoryUpload.FileName,
						StatusFile: "finished_process",
						Description: "Total record : " + fmt.Sprintf("%d", totalOfRecord) + ", Success insert : " + fmt.Sprintf("%d", successInsert) + ", Failed insert :" + fmt.Sprintf("%d", failedInsert),
					})
					if errUpdateStatus != nil {
						log.Error("Error background process insert file zsd081 :", f.Name(), errUpdateStatus.Error())
						isAnyErrorProcessFile = true
					}
					errMoveFile := utils.MoveFile(os.Getenv("ROOT_UPLOAD_ZSD081_PATH")+f.Name(), os.Getenv("MOVE_SOURCE_PATH_SUCCESS")+f.Name(), false)

					if errMoveFile != nil {
						log.Error("Error background process insert file zsd081 :", f.Name(), errMoveFile.Error())
						isAnyErrorProcessFile = true
					}

				} else {
					errMoveFile := utils.MoveFile(os.Getenv("ROOT_UPLOAD_ZSD081_PATH")+f.Name(), os.Getenv("MOVE_SOURCE_PATH_UNSUCCESS")+f.Name(), false)
					if errMoveFile != nil {
						log.Error("Error background process insert file zsd081 :", f.Name()+" with error : ", errMoveFile.Error())
					}
					log.Error("Error background process insert file zsd081 :", f.Name()+" with error : ", "header tidak sesuai format")
					isAnyErrorProcessFile = true
				}

			//} else {
			//	errMoveFile := utils.MoveFile(os.Getenv("ROOT_UPLOAD_ZSD081_PATH")+f.Name(),  os.Getenv("MOVE_SOURCE_PATH_UNSUCCESS") + f.Name())
			//	if errMoveFile != nil {
			//		log.Error("Error insert file zsd081 :", f.Name(), errMoveFile.Error())
			//	}
			//	log.Error("Error insert file zsd081 :", f.Name(), "filename tidak sesuai format")
			//	isAnyErrorProcessFile = true
			//}
		}
		totalOfRecord = 0
		successInsert = 0
		failedInsert = 0
	}

	msg := ""
	if isAnyErrorProcessFile {
		msg = "Any wrong process file, "
	}

	if isAnyErrorProcessRecord {
		msg += " Any wrong process record"
	}

	if countProcessFile == 0 {
		msg = "no file process"
	}

	if msg == "" {
		log.Info("success insert zsd081")
		ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "success insert zsd081", utils.SUCCESS_CODE))
	} else {
		log.Info("Error background prosess insert :", msg)
		ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, msg, utils.ERR_GLOBAL))
	}

}
