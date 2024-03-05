package controller

import (
	"backend_gui/dto"
	"backend_gui/repository"
	"backend_gui/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

type TestController interface {
  SelectAllDb(ctx *gin.Context)
  SelectAllFolder(ctx *gin.Context)

}

type testController struct {
	testRepo repository.TestRepository
}

func (t testController) SelectAllFolder(ctx *gin.Context) {
	files, err := ioutil.ReadDir(os.Getenv("ROOT_UPLOAD_FPM_PATH"))
	if err != nil{
		log.Error("error check folder : ", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("Somenthing wrong", utils.ERR_BAD_REQUEST))
		return
	}
	var folder []string
	anyError := false
	for _, file := range files {
		if file.IsDir(){
			if file.Name() == "SF" || file.Name() == "DSJ" ||  file.Name() == "ST" || file.Name() == "SAP"{
				folder = append(folder,file.Name())
				subFolder, errRead := ioutil.ReadDir(os.Getenv("ROOT_UPLOAD_FPM_PATH")+"/"+file.Name())
				if errRead != nil{
					anyError = true
					log.Error("error check folder : ", errRead.Error())
					break
				}
				for _, checkSub :=range subFolder{
                  if checkSub.IsDir(){
                  	if checkSub.Name() == "10" || checkSub.Name() == "20" || checkSub.Name() == "35" ||
                  		checkSub.Name() == "40" ||  checkSub.Name() == "45" || checkSub.Name() == "45" || checkSub.Name() == "others" || checkSub.Name() == "Archive_toSAP" || checkSub.Name() == "Email_Archive"{
						folder = append(folder,checkSub.Name())
					  }
				  }
				}
			}

		}
	}
	if anyError{
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("Somenthing wrong", utils.ERR_BAD_REQUEST))
		return
	}
	_, errReadFolderZip := ioutil.ReadDir(os.Getenv("ROOT_UPLOAD_SOURCE_DATA_PATH"))
	if errReadFolderZip != nil{
		log.Error("error check folder : ", errReadFolderZip.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("Somenthing wrong", utils.ERR_BAD_REQUEST))
		return
	}
	folder = append(folder, "source_data")
    log.Info("folder at server : ", folder)
	if len(folder) != 25{
		log.Error("error check folder : ", "count folder != 24")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("folder not completed", utils.ERR_BAD_REQUEST))
		return
	}
	ctx.JSON(http.StatusOK, dto.SuccessResponse(nil,"success test all folder test at "+os.Getenv("ENV_VERSION"), utils.SUCCESS_CODE))
}

func (t testController) SelectAllDb(ctx *gin.Context) {
   errSelectRole := t.testRepo.SelectTableRole("sf_fpm_role")
   if errSelectRole != nil{
	   ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse(errSelectRole.Error(), utils.ERR_BAD_REQUEST))
	   return
   }
   errSelectUserFpm := t.testRepo.SelectTableUser("sf_fpm_user")
   if errSelectUserFpm != nil{
	   ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse(errSelectRole.Error(), utils.ERR_BAD_REQUEST))
	   return
   }
   errSelectFaktur := t.testRepo.SelectTableFaktur("sf_fpm_faktur")
   if errSelectFaktur != nil{
	   ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse(errSelectFaktur.Error(), utils.ERR_BAD_REQUEST))
	   return
   }
   errSelectHisPjk := t.testRepo.SelectTableHisUploadPjk("sf_fpm_history_upload_pjk")
   if errSelectHisPjk !=nil{
	   ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse(errSelectHisPjk.Error(), utils.ERR_BAD_REQUEST))
	   return
   }
   errSelectEmailFpm := t.testRepo.SelectTableEmailFpm("sf_fpm_email_notification")
	if errSelectEmailFpm !=nil{
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse(errSelectEmailFpm.Error(), utils.ERR_BAD_REQUEST))
		return
	}

	errSelecCompareFpm := t.testRepo.SelectTableCompareFpm("sf_fpm_compare")
	if errSelecCompareFpm !=nil{
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse(errSelectEmailFpm.Error(), utils.ERR_BAD_REQUEST))
		return
	}
	errSelecZsd081 := t.testRepo.SelectTableZsd081("sf_dump_zsd081")
	if errSelecZsd081 !=nil{
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse(errSelecZsd081.Error(), utils.ERR_BAD_REQUEST))
		return
	}
	ctx.JSON(http.StatusOK, dto.SuccessResponse(nil,"success test all table test at "+os.Getenv("ENV_VERSION"), utils.SUCCESS_CODE))
}

func InstanceTestController(db repository.TestRepository) TestController {
	return &testController{
		testRepo: db,
	}

}
