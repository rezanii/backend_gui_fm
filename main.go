package main

import (
	"backend_gui/connectionSetup"
	"backend_gui/controller"
	"backend_gui/middleware"
	"backend_gui/repository"
	"backend_gui/utils"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"regexp"
)

var (
	db                       *gorm.DB                            = connectionSetup.SetupDbConnection()
	mapUploadRepo            repository.UploadRepository         = repository.InstanceUploadRepository(db)
	mapProcessZsd081Repo     repository.ProcessZsd081Repository  = repository.InstanceProcessZsd081Repository(db)
	mapProcessFileRepo       repository.FpmRepository            = repository.InstanceFpmRepository(db)
	mapProcessCompareRepo    repository.ProcessCompareRepository = repository.InstanceProcessCompareRepository(db)
	mapEmailRepo             repository.EmailRepository          = repository.InstanceEmailRepository(db)
	userFpmRepo              repository.UserFpmRepository        = repository.InstanceUserFpmRepository(db)
	testRepo                 repository.TestRepository           = repository.InstanceTestRepository(db)
	processZsd081Controller  controller.ProcessZsd081Controller  = controller.InstanceProcessZsd081Controller(mapProcessZsd081Repo)
	uploadFileController     controller.UploadController         = controller.InstanceUploadController(mapUploadRepo)
	processFileController    controller.FpmController            = controller.InstanceFpmController(mapProcessFileRepo)
	processCompareController controller.ProcessCompareController = controller.InstanceProcessCompareController(mapProcessCompareRepo)
	emailController          controller.EmailController          = controller.InstanceEmailController(mapEmailRepo)
	userFpmController        controller.UserFpmController        = controller.InstanceUserFpmController(userFpmRepo)
	testController           controller.TestController           = controller.InstanceTestController(testRepo)

)

func main() {
	defer connectionSetup.CloseDbConnection(db)
	log.Info("************* deployment "+os.Getenv("ENV_VERSION")+" *************")
	r := gin.Default()

	utils.RedisConnection = connectionSetup.RedisConnection()
	r.SetTrustedProxies([]string{os.Getenv("IP")})
	r.Use(CORSMiddlewareSession())
	r.Static("/backend-gui/file-fpm", os.Getenv("ROOT_UPLOAD_FPM_PATH"))
	r.Static("/backend-gui/source-data", os.Getenv("ROOT_UPLOAD_SOURCE_DATA_PATH"))
	r.Static("/backend-gui/temp-file", os.Getenv("ROOT_TEMP_UPLOAD"))
	r.Use(LogRequst())
	privateRoutes := r.Group("/backend-gui/private", middleware.Authorization(mapUploadRepo))
	{

		privateRoutes.GET("/session", uploadFileController.CheckSession)

		privateRoutes.POST("/finance_tax", uploadFileController.UploadFpm)
		privateRoutes.POST("/saveDataUpload", uploadFileController.SaveDataUpload)
		privateRoutes.POST("/deleteUploadFile", uploadFileController.DeleteTempFiles)
		privateRoutes.POST("/getHistoryFile", uploadFileController.GetHistoryUploadFpm)

		privateRoutes.POST("/source_data", uploadFileController.UploadSourceData)
		privateRoutes.POST("/saveSourceData", uploadFileController.SaveSourceData)


		adminRoutes := privateRoutes.Group("/admin", middleware.AuthorizationAdmin(userFpmRepo))
		{
			adminRoutes.POST("/addUser", userFpmController.AddUserFpm)
			adminRoutes.POST("/deleteUser", userFpmController.DeleteUserFpm)
			adminRoutes.POST("/listUserFpm", userFpmController.ListUserFpm)
			adminRoutes.GET("/listRoleUser", userFpmController.GetListRoleUser)

		}
		privateRoutes.POST("/collectionInv",processCompareController.GetCollectionInv)
		privateRoutes.POST("/downloadFpm", uploadFileController.DownloadFpm)
	}



	backProsesRoutes := r.Group("/backend-gui/process")
	{
		//backProsesRoutes.POST("/insertZsd081", processZsd081Controller.InsertDataZsd081)
		backProsesRoutes.POST("/insertZsd001nToCompare", processCompareController.InsertZsd001nToCompare)
		backProsesRoutes.POST("/checkInvoiceAlreadySend", processCompareController.CheckInvoiceAlreadySend)
		backProsesRoutes.POST("/compareInvoiceCancel", processCompareController.CompareInvoiceCancel)
		backProsesRoutes.POST("/compareZsd001nZv60", processCompareController.CompareZsd001nZv60)
		backProsesRoutes.POST("/compareFpCancel", processCompareController.CompareFpCancel)
		backProsesRoutes.POST("/checkFpAlreadySend", processCompareController.CheckFpAlreadySend)

		backProsesRoutes.POST("/moveFpm", processFileController.MoveFpm)
		backProsesRoutes.POST("/updateStatusFpm", processFileController.UpdateStatusFile)

		backProsesRoutes.POST("/emailCancelInv", emailController.EmailCancelInv)
		backProsesRoutes.POST("/emailInvNoFp", emailController.EmailInvNoFp)
		backProsesRoutes.POST("/emailInvNotSend", emailController.EmailInvNotSend)
		backProsesRoutes.POST("/emailNoFileFp", emailController.EmailNoFileFp)
		backProsesRoutes.POST("/emailFpCancel", emailController.EmailFpCancel)
		backProsesRoutes.POST("/emailFpNotSend", emailController.EmailFpNotSend)
		backProsesRoutes.POST("/deleteZip", uploadFileController.HouseKeepingZipFile)

	}
	publicRoutes:= r.Group("/backend-gui/test/")
	{
		publicRoutes.POST("/db", testController.SelectAllDb)
		publicRoutes.POST("/folder", testController.SelectAllFolder)
	}

	r.Run(":" + os.Getenv("PORT_GO"))
}

/*
fuction to set cors allow access api
by irma 30/12/2021
*/
func CORSMiddlewareSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("FLUTTER_URL"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()

	}
}

func LogRequst() gin.HandlerFunc {
	return func(context *gin.Context) {
		buffer, _ := ioutil.ReadAll(context.Request.Body)
		body := utils.ReadBody(ioutil.NopCloser(bytes.NewBuffer(buffer)))
		context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buffer))
		response := &utils.BodyLogWriter{Body: bytes.NewBufferString(""), ResponseWriter: context.Writer}
		context.Writer = response
		re := regexp.MustCompile(`\r?\n`)

		_, err := context.MultipartForm()
		if err != nil {
			log.Info(fmt.Sprintf("requestBody: %s",re.ReplaceAllString(body,"") ))
		}
		log.Info(fmt.Sprintf("requestHeader: %s",context.Request.Header))
		context.Next()
		log.Info(fmt.Sprintf("response: %s",re.ReplaceAllString(response.Body.String(),"")))

	}
}