package middleware

import (
	"backend_gui/dto"
	"backend_gui/models"
	"backend_gui/repository"
	"backend_gui/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)
/*
	Global Method to interupt aksess endpoint group and check session
*/
func Authorization(mapUploadRepo repository.UploadRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.CheckDayLog()
		username, err := utils.GetSession(c, "username")
		if err != nil {
			log.Error("Error validate session :", err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
			return
		}
       fmt.Println(fmt.Sprintf("%s", username))
		userFpm := mapUploadRepo.CheckUserAlreadyExist(fmt.Sprintf("%s", username))
		if userFpm.Id == 0{
			log.Error("Error check session : "," user not exist at db")
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
			return
		}

	}
	
	

}
/*
	Method global authorization for access endpoint
*/
func AuthorizationAdmin(mapUserRepo repository.UserFpmRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.CheckDayLog()
		username, err := utils.GetSession(c, "username")
		if err != nil {
			log.Error("Error validate session :", err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
			return
		}

		userFpm := mapUserRepo.IsRoleAdmin(fmt.Sprintf("%s", username))
		if userFpm == (models.UserFpmModels{}){
			log.Error("Error check user access : "," not admin")
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
			return
		}

	}


}

