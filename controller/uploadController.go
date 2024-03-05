package controller

import (
	"archive/zip"
	"backend_gui/dto"
	"backend_gui/models"
	"backend_gui/repository"
	"backend_gui/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type UploadController interface {
	DeleteTempFiles(ctx *gin.Context)
	SaveDataUpload(ctx *gin.Context)
	GetHistoryUploadFpm(ctx *gin.Context)
	SaveSourceData(ctx *gin.Context)
	UploadSourceData(ctx *gin.Context)
	UploadFpm(c *gin.Context)
	CheckSession(ctx *gin.Context)
	DownloadFpm(ctx *gin.Context)
	HouseKeepingZipFile(ctx *gin.Context)
}

type uploadController struct {
	mapUploadRepo repository.UploadRepository
}
/*
	Method process Housekeeping file zip
*/
func (u uploadController) HouseKeepingZipFile(ctx *gin.Context) {
	path := os.Getenv("ROOT_UPLOAD_SOURCE_DATA_PATH")
	files, errReadDir := ioutil.ReadDir(path)
	if errReadDir != nil {
		log.Error("house keeping file zip fpm, error :", dto.ErrorResponse(errReadDir.Error(), utils.ERR_VALIDATE_DATA))
		return
	}
	countFileProcess := 0
	for _, f := range files {
		if !f.IsDir(){
			var extFile = strings.Split(f.Name(),".")[1]
			if extFile == "zip"{
				sufix := []rune(f.Name())
				if string(sufix[:3]) == "FPM"{
					countFileProcess++
					log.Info("house keeping file zip fpm, process delete file zip : ", f.Name())
					errDelete := os.Remove(path+f.Name())
					if errDelete != nil {
						log.Error("house keeping file zip fpm, error :", errDelete.Error())
					}else{
						log.Info("house keeping file zip fpm, success delete file zip : ", f.Name())
					}
				}
			}
		}

	}
	if countFileProcess == 0{
		log.Info("house keeping file zip fpm : ", "file empty")
	}

}
/*
	Method download file fpm
*/

func (u uploadController) DownloadFpm(ctx *gin.Context) {
	username, errSession := utils.GetSession(ctx, "username")
	if errSession != nil {
		log.Error("upload source data, error :", dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}
	var reqDownloadFpm dto.ReqDownloadFpmDto
	errScan := ctx.ShouldBindJSON(&reqDownloadFpm)
	if errScan != nil {
		log.Error("add user fpm, username : "+fmt.Sprintf("%s", username)+", error : ", errScan.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse(errScan.Error(), utils.ERR_VALIDATE_DATA))
		return
	}

	if len(reqDownloadFpm.IdFiles) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse(errScan.Error(), utils.ERR_VALIDATE_DATA))
		return
	}
	currentTime := time.Now()
	dateNow := currentTime.Format("2006-01-02 15:04:05.000000")
	dateSplit :=strings.Split(dateNow," ")
	dateZip:= strings.Replace(dateSplit[0],"-","", -1)
	timeSplit:= strings.Replace(strings.Split(dateSplit[1],":")[2],".","",-1)

	AnyError := false
	data := u.mapUploadRepo.GetDownloadFpmById(reqDownloadFpm.IdFiles, fmt.Sprintf("%s",username))
	log.Info("Zip file size data, username :  "+fmt.Sprintf("%s", username), len(data))
	fileNameZip := "FPM"+dateZip+timeSplit+".zip"
	if len(data) > 0 {
		newZipFile, errCreate := os.Create(os.Getenv("ROOT_UPLOAD_SOURCE_DATA_PATH")+fileNameZip)
		if errCreate != nil {
			log.Error("zip file fpm, username : "+fmt.Sprintf("%s", username)+", error : ", errCreate)
		}
		defer newZipFile.Close()

		zipWriter := zip.NewWriter(newZipFile)
		defer zipWriter.Close()

		for i := 0; i < len(data); i++ {
			pathFile:=""
           if data[i].Status != "uploaded"{
           	if data[i].InvDc != "10" && !utils.CheckChannel(data[i].InvDc){
				pathFile = os.Getenv("ROOT_UPLOAD_FPM_PATH")+data[i].JenisFakturPajak+"/"+os.Getenv("FOLDER_OTHER_DC")+"/"+data[i].FileNameUpload
			}else{
				pathFile = os.Getenv("ROOT_UPLOAD_FPM_PATH")+data[i].JenisFakturPajak+"/"+data[i].InvDc+"/"+data[i].FileNameUpload
			}
			}else{
				pathFile = os.Getenv("ROOT_UPLOAD_FPM_PATH")+data[i].JenisFakturPajak+"/"+data[i].FileNameUpload
			}
			err:=utils.AddFileToZip(zipWriter, pathFile, data[i].FileNameUpload)
			if err!=nil{
				AnyError = true
				log.Error("zip file fpm , username :"+fmt.Sprintf("%s", username)+", error : ", err.Error())
				break
			}
		}
	}else{
		log.Error("zip file fpm, username :"+fmt.Sprintf("%s", username)+", error : ", "file not found")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("file not found", utils.ERR_GLOBAL))
		return
	}

	if AnyError{
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("Somenthing wrong", utils.ERR_GLOBAL))
		return
	}
	dataRes:= dto.RespDownloadFpm{
		BaseUrl: os.Getenv("BASE_URL_FILE_SOURCE_DATA")+"source-data",
		FileName: fileNameZip,

	}
	log.Info("zip file fpm, success create : ", dataRes)
	ctx.JSON(http.StatusOK, dto.SuccessResponse(dataRes,"success create zip fpm", utils.SUCCESS_CODE))


}

/*
fuction to check session
by irma 30/12/2021
*/
func (u uploadController) CheckSession(ctx *gin.Context) {
	username, err := utils.GetSession(ctx, "username")
	if err != nil {
		log.Error("Error check session :", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}
	userFpm := u.mapUploadRepo.CheckUserAlreadyExist(fmt.Sprintf("%s", username))
	if userFpm.Id == 0 {
		log.Error("Error check session : ", " user not exist at db")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}

	log.Info("success connection with username ", fmt.Sprintf("%s", username))
	ctx.JSON(http.StatusOK, dto.SuccessResponse(&map[string]string{
		"username": fmt.Sprintf("%s", username),
		"roleUser": userFpm.RoleUser,
	}, "", ""))
	return
}

/*
fuction to upload file fpm to temporary files and respon total, link, filenames
by irma 30/12/2021
*/
func (u uploadController) UploadFpm(c *gin.Context) {
	log.Info("Process upload fpm ...")
	folder := os.Getenv("ROOT_TEMP_UPLOAD")
	pathUrlDownload := os.Getenv("BASE_URL_FILE_FAKTUR_PAJAK") + "temp-file/"
	username, errSession := utils.GetSession(c, "username")
	if errSession != nil {
		log.Error("upload file fpm, error :", dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}

	userFpm := u.mapUploadRepo.CheckUserAlreadyExist(fmt.Sprintf("%s", username))

	if userFpm.RoleUser != "TAX_TEAM" {
		log.Error("upload file fpm, username", fmt.Sprintf("%s", username)+", error :", "user have not access this service")
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Error("upload file fpm, username", fmt.Sprintf("%s", username)+", error :", fmt.Sprintf("err: %s", err.Error()))
		c.JSON(http.StatusOK, dto.ErrorResponse(fmt.Sprintf("err: %s", err.Error()), utils.ERR_VALIDATE_DATA))
		return
	}
	files := form.File["files"]
	//var succesFile []dto.ResSaveUploadFile
	succesFile := []string{}
	fileNames := []string{}
	if len(files) > utils.LIMIT_OF_UPLOAD_FPM {
		log.Error("upload file fpm, username : ", fmt.Sprintf("%s", username)+", error :", "size upload must be <= 100")
		c.JSON(http.StatusOK, dto.ErrorResponse("size exceeds limit of upload faktur", utils.ERR_VALIDATE_DATA))
		return
	}
	for _, file := range files {
		path := folder + file.Filename
		errUploadFpm := c.SaveUploadedFile(file, path)
		if errUploadFpm == nil {
			succesFile = append(succesFile, pathUrlDownload+file.Filename)
			fileNames = append(fileNames, file.Filename)
			//items := dto.ResSaveUploadFile{
			//	FileName: file.Filename,
			//	Url:      pathUrlDownload + file.Filename,
			//}
			//succesFile = append(succesFile, items)

			log.Info("success upload fpm, username "+fmt.Sprintf("%s", username)+" ,filename "+file.Filename+" and url :", pathUrlDownload+file.Filename)
		} else {
			log.Error("upload file fpm, username  " + fmt.Sprintf("%s", username) + " ,filename " + file.Filename + ", error :" + fmt.Sprintf("err: %s", errUploadFpm.Error()))
		}
	}
	if len(succesFile) > 0 {
		//c.JSON(http.StatusOK, dto.SuccessResponse(dto.RespTempUploadDtoNew{
		//	Total: len(succesFile),
		//	Data:  succesFile,
		//}, "", ""))
		c.JSON(http.StatusOK, dto.SuccessResponse(dto.RespTempUploadDto{
			Total:           len(succesFile),
			Urls:            succesFile,
			FileNamesUpload: fileNames,
		}, "", ""))
		return
	} else {
		log.Error("upload file fpm, username "+fmt.Sprintf("%s", username)+", error :", "Failed save all file to server")
		c.JSON(http.StatusOK, dto.ErrorResponse("Failed", "86"))
	}
	log.Info("End process upload fpm ...")
}

/*
fuction to upload zsd081 to temporary files and return respon total data,
link, filename
by irma 30
*/
func (u uploadController) UploadSourceData(ctx *gin.Context) {
	log.Info("Process upload upload source data ...")
	folder := os.Getenv("ROOT_TEMP_UPLOAD")
	pathUrlDownload := os.Getenv("BASE_URL_FILE_SOURCE_DATA") + "temp-file/"
	username, errSession := utils.GetSession(ctx, "username")
	if errSession != nil {
		log.Error("upload source data, error :", dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}
	userFpm := u.mapUploadRepo.CheckUserAlreadyExist(fmt.Sprintf("%s", username))
	if userFpm.RoleUser == "ADMIN" || userFpm.RoleUser == "FOS" {
		form, err := ctx.MultipartForm()
		if err != nil {
			log.Error("Error upload file source data username, "+fmt.Sprintf("%s", username)+":", err.Error())
			ctx.JSON(http.StatusOK, dto.ErrorResponse(fmt.Sprintf("err: %s", err.Error()), utils.ERR_VALIDATE_DATA))
			return
		}
		files := form.File["files"]
		//var succesFile []dto.ResSaveUploadFile
		succesFile := []string{}
		fileNames := []string{}

		for _, fileNotDuplicate := range files {
			path := folder + fileNotDuplicate.Filename
			err := ctx.SaveUploadedFile(fileNotDuplicate, path)
			if err == nil {
				//items := dto.ResSaveUploadFile{
				//	FileName: fileNotDuplicate.Filename,
				//	Url:      pathUrlDownload + fileNotDuplicate.Filename,
				//}
				//succesFile = append(succesFile, items)
				succesFile = append(succesFile, pathUrlDownload+fileNotDuplicate.Filename)
				fileNames = append(fileNames, fileNotDuplicate.Filename)
				log.Info("success upload zsd081, username "+fmt.Sprintf("%s", username)+" ,filename :"+fileNotDuplicate.Filename+" and url :", pathUrlDownload+fileNotDuplicate.Filename)
			} else {
				log.Error("Error upload file zsd081, username "+fmt.Sprintf("%s", username)+" ,filename :"+fileNotDuplicate.Filename+" with error :", err.Error())
			}
		}

		if len(succesFile) > 0 {
			//ctx.JSON(http.StatusOK, dto.SuccessResponse(dto.RespTempUploadDtoNew{
			//	Total: len(succesFile),
			//	Data:  succesFile,
			//}, "", ""))
			ctx.JSON(http.StatusOK, dto.SuccessResponse(dto.RespTempUploadDto{
				Total:           len(succesFile),
				Urls:            succesFile,
				FileNamesUpload: fileNames,
			}, "", ""))
			return
		} else {
			log.Error("Error upload file zsd081, username "+fmt.Sprintf("%s", username)+":", "Failed save all file to server")
			ctx.JSON(http.StatusOK, dto.ErrorResponse("Failed", "86"))
		}
	} else {
		log.Error("upload source data, username : ", fmt.Sprintf("%s", username)+", error : ", " user have not access this service")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}
	log.Info("End process upload upload source data ...")
}

/*
fuction to validate, move file zsd081 to folder source_data and save history
upload zsd081 to db
by irma 30/12/2021
*/
func (u uploadController) SaveSourceData(ctx *gin.Context) {
	log.Info("save upload source data process ...")
	username, errSession := utils.GetSession(ctx, "username")
	if errSession != nil {
		log.Error("save upload source data , username, error :", errSession.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}
	userFpm := u.mapUploadRepo.CheckUserAlreadyExist(fmt.Sprintf("%s", username))
	if userFpm.RoleUser == "ADMIN" || userFpm.RoleUser == "FOS" {
		getJwt, errJwt := utils.GenerateToken(fmt.Sprintf("%s", username))
		keyJwt := "Bearer " + getJwt

		if errJwt != nil {
			log.Error("save upload source data , username, error:", errJwt.Error())
			ctx.JSON(http.StatusOK, dto.ErrorResponse("somenthing wrong", utils.ERR_VALIDATE_DATA))
			return
		}

		var reqData dto.ReqSaveSourceDataDto
		errScanJsonValue := ctx.ShouldBindJSON(&reqData)

		if errScanJsonValue != nil {
			log.Error("save upload source data , username "+fmt.Sprintf("%s", username)+", error:", errScanJsonValue.Error())
			ctx.JSON(http.StatusOK, dto.ErrorResponse(errScanJsonValue.Error(), utils.ERR_VALIDATE_DATA))

		}

		if len(reqData.Files) == 0 {
			log.Error("save upload source data , username "+fmt.Sprintf("%s", username)+", error :", "files can't empty")
			ctx.JSON(http.StatusOK, dto.ErrorResponse("files can't empty", utils.ERR_VALIDATE_DATA))
		}

		isAnySuccess := false
		errorMessage := ""

		pathFile := os.Getenv("ROOT_TEMP_UPLOAD")
		url, typeFileUpload := utils.GetUrlUploadSourceData(reqData.TypeFile)
		fmt.Println(reqData.TypeFile)
		fmt.Println(typeFileUpload)
		for i := 0; i < len(reqData.Files); i++ {
			res, errUpload := utils.SendUpload(url, keyJwt, reqData.Files[i], pathFile, typeFileUpload)
			log.Info("Finish process curl upload file source data : ", res)
			if errUpload == nil && res.Meta.Code == 200 {
				isAnySuccess = true
				errorMessage += "SUCCESS|File " + reqData.Files[i] + " Success to save\n"

			} else {
				errorMessage += "ERROR|File " + reqData.Files[i] + "  failed to save\n"
				log.Error("save upload source data, username : "+fmt.Sprintf("%s", username)+"error, :", res)
			}

		}

		log.Info("Info save upload source data ,username "+fmt.Sprintf("%s", username)+":", errorMessage)
		if isAnySuccess {
			log.Info("save upload source data process delete file...")
			for x := 0; x < len(reqData.Files); x++ {
				errRemove := os.Remove(pathFile + reqData.Files[x])
				if errRemove != nil {
					log.Error("save upload source data , username : "+reqData.Files[x]+", error :", errRemove)
				}
			}
		}
		if !isAnySuccess {
			log.Error("Error save file source data ,username "+fmt.Sprintf("%s", username), "Failed all save file source data")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse(errorMessage, utils.ERR_GLOBAL))
			return
		}

		log.Info("", dto.SuccessResponse(nil, errorMessage, utils.SUCCESS_CODE))
		ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, errorMessage, utils.SUCCESS_CODE))
	} else {
		log.Error("save upload source data, username : ", fmt.Sprintf("%s", username)+", error : ", " user have not access this service")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}

	log.Info("End process save upload source data ....")
}

/*
fuction get data history upload file faktur pajak with konsep pagination
by irma 30/12/2021
*/
func (u uploadController) GetHistoryUploadFpm(ctx *gin.Context) {
	log.Info("Process get history fpm ....")
	username, errSession := utils.GetSession(ctx, "username")
	if errSession != nil {
		log.Info("get history fpm, error :", errSession.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}

	userFpm := u.mapUploadRepo.CheckUserAlreadyExist(fmt.Sprintf("%s", username))

	if userFpm.RoleUser == "FOS" || userFpm.RoleUser == "ADMIN" || userFpm.RoleUser == "TAX_TEAM" {
		var reqData dto.ReqPaginationDto
		err := ctx.ShouldBindJSON(&reqData)
		if err != nil {
			log.Info("get history fpm,username "+fmt.Sprintf("%s", username)+", Error :", err.Error())
			ctx.JSON(http.StatusOK, dto.ErrorResponse(err.Error(), utils.ERR_VALIDATE_DATA))
			return
		}
		start := (reqData.Page * reqData.MaxDataDisplay) - reqData.MaxDataDisplay

		data := dto.ResPaginationDto{
			BaseUrl:  os.Getenv("BASE_URL_FILE_FAKTUR_PAJAK")+"file-fpm",
			TotalData: u.mapUploadRepo.GetTotalDataFpm(reqData.Search, fmt.Sprintf("%s", username), userFpm.RoleUser),
			Record:    u.mapUploadRepo.GetDataHistoryUploadFpm(start, reqData.MaxDataDisplay, reqData.Search, fmt.Sprintf("%s", username), userFpm.RoleUser),
		}
		ctx.JSON(http.StatusOK, dto.SuccessResponse(data, "", utils.SUCCESS_CODE))

	} else {
		log.Error("get history fpm, username: ", fmt.Sprintf("%s", username), " user have not access at this service")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}

	log.Info("End process get history fpm ....")

}

func InstanceUploadController(mapUploadRepo repository.UploadRepository) UploadController {
	return &uploadController{
		mapUploadRepo: mapUploadRepo,
	}

}

/*
fuction to move file faktur pajak to DSJ/ST/SF
and save history file faktur pajak to db
by irma 30/12/2021
*/
func (u uploadController) SaveDataUpload(ctx *gin.Context) {
	log.Info("Process save upload fpm ...")
	var reqData dto.ReqSaveDataUploadDto
	username, errSession := utils.GetSession(ctx, "username")
	if errSession != nil {
		log.Error("save fpm, error:", dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}
	userFpm := u.mapUploadRepo.CheckUserAlreadyExist(fmt.Sprintf("%s", username))

	if userFpm.RoleUser != "TAX_TEAM" {
		log.Error("save fpm, username", fmt.Sprintf("%s", username)+", error :", " user have not access this service")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}

	errScan := ctx.ShouldBindJSON(&reqData)
	if errScan != nil {
		log.Error("save fpm, username "+fmt.Sprintf("%s", username)+", error :", dto.ErrorResponse(errScan.Error(), utils.ERR_VALIDATE_DATA))
		ctx.JSON(http.StatusOK, dto.ErrorResponse(errScan.Error(), utils.ERR_VALIDATE_DATA))
		return
	}

	var folder = ""
	if reqData.JenisFakturPajak == "1" {
		folder = "SF"
	} else if reqData.JenisFakturPajak == "2" {
		folder = "ST"
	} else {
		folder = "DSJ"
	}
	errorMessage := ""
	isAnySuccess := false

	for x := 0; x < len(reqData.FileNameUpload); x++ {
		takeNoFp := strings.Split(reqData.FileNameUpload[x], "-")[1]
		noFp := takeNoFp[8:len(takeNoFp)]
		checkFileExistAtDb := u.mapUploadRepo.GetHistoryFpmByFileName((reqData.FileNameUpload[x]), fmt.Sprintf("%s", username))
		if checkFileExistAtDb > 0 {
			log.Error("save data upload fpm, username "+fmt.Sprintf("%s", username)+", filename "+(reqData.FileNameUpload[x])+", error : ", " data already exist at db")
			continue
		}
		err := utils.MoveFile(os.Getenv("ROOT_TEMP_UPLOAD")+reqData.FileNameUpload[x], os.Getenv("ROOT_UPLOAD_FPM_PATH")+folder+"/"+reqData.FileNameUpload[x], false)

		if err != nil {
			log.Info("save data upload fpm, username "+fmt.Sprintf("%s", username)+", filename "+reqData.FileNameUpload[x]+", error :", err.Error())
			errorMessage += "File " + reqData.FileNameUpload[x] + " Failed to save\n"
		} else {

			savaData := models.FpmHistoryFilePjk{
				UploadBy:         fmt.Sprintf("%s", username),
				JenisFakturPajak: folder,
				FileNameUpload:   reqData.FileNameUpload[x],
				Url:              os.Getenv("BASE_URL_FILE_FAKTUR_PAJAK") + "file-fpm/" + folder + "/" + reqData.FileNameUpload[x],
				NoFaktur:         noFp,
				Status:           "uploaded",
			}
			log.Info("Save data upload fpm, username "+fmt.Sprintf("%s", username)+", with data :", savaData)
			errSave := u.mapUploadRepo.SaveDataUploadFilePjk(savaData)
			if errSave == nil {
				errorMessage += "SUCCESS|File " + reqData.FileNameUpload[x] + " Success to save\n"
				isAnySuccess = true
			} else {
				errorMessage += "ERROR|File " + reqData.FileNameUpload[x] + " Failed to save\n"
				log.Info("Error Save data upload fpm, username "+fmt.Sprintf("%s", username)+" ,filename "+reqData.FileNameUpload[x]+",  error :", errSave.Error())
			}
		}

	}

	log.Info("Info save file fpm, username "+fmt.Sprintf("%s", username)+":", errorMessage)
	if !isAnySuccess {
		log.Error("save fpm, username "+fmt.Sprintf("%s", username)+":", "Failed save all file fpm")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("Failed to save", utils.ERR_GLOBAL))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, errorMessage, utils.SUCCESS_CODE))
	log.Info("End process save upload fpm ...")
}

/*
fuction to delete file upload
by irma 30/12/2021
*/
func (u uploadController) DeleteTempFiles(ctx *gin.Context) {
	log.Info("Process delete file upload ...")
	var reqData dto.ReqDeleteFileDto
	username, errSession := utils.GetSession(ctx, "username")
	if errSession != nil {
		log.Error("delete file, error ", dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}
	userFpm := u.mapUploadRepo.CheckUserAlreadyExist(fmt.Sprintf("%s", username))

	if userFpm.RoleUser == "TAX_TEAM" || userFpm.RoleUser == "FOS" || userFpm.RoleUser == "ADMIN" {
		errScan := ctx.ShouldBindJSON(&reqData)
		if errScan != nil {
			log.Error("delete file, username "+fmt.Sprintf("%s", username)+", error:", errScan.Error())
			ctx.JSON(http.StatusOK, dto.ErrorResponse(errScan.Error(), utils.ERR_VALIDATE_DATA))
			return
		}
		path := os.Getenv("ROOT_TEMP_UPLOAD")
		if _, errFindFile := os.Stat(path + reqData.FileName); os.IsNotExist(errFindFile) {
			log.Error("delete file, username "+fmt.Sprintf("%s", username)+", filename "+reqData.FileName+", error:", errFindFile.Error())
			ctx.JSON(http.StatusOK, dto.ErrorResponse("File not found", utils.ERR_FILE_NOT_FOUND))
			return
		}

		e := os.Remove((path + reqData.FileName))
		if e != nil {
			log.Error("delete file, username "+fmt.Sprintf("%s", username)+", filename "+reqData.FileName+", error:", e.Error())
			ctx.JSON(http.StatusOK, dto.ErrorResponse("Failed to delete file", utils.ERR_DELETE_FILE))
			return
		}
		log.Info("delete file, username "+fmt.Sprintf("%s", username)+" success, with file :", path+reqData.FileName)
		ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "Success delete file", utils.SUCCESS_CODE))
	} else {
		log.Error("delete file fpm, username", fmt.Sprintf("%s", username)+", error :", " user have not access this service")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}
	log.Info("End process delete file upload ...")
}
