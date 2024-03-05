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
	"strconv"
)

type UserFpmController interface {
	AddUserFpm(ctx *gin.Context)
	DeleteUserFpm(ctx *gin.Context)
	ListUserFpm(ctx *gin.Context)
	GetListRoleUser(ctx *gin.Context)

}

type userFpmController struct {
	mapUserFpmRepo repository.UserFpmRepository
}
/*
	Method get list role user
*/
func (u userFpmController) GetListRoleUser(ctx *gin.Context) {
   data := u.mapUserFpmRepo.GetListRoleUser()
	ctx.JSON(http.StatusOK, dto.SuccessResponse(data, "", utils.SUCCESS_CODE))
}
/*
	Method get list user
*/
func (u userFpmController) ListUserFpm(ctx *gin.Context) {
	username, err := utils.GetSession(ctx, "username")
	if err != nil {
		log.Error("delete user fpm, error : ", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}
	var reqData dto.ReqPaginationDto
	errScan := ctx.ShouldBindJSON(&reqData)
	if errScan != nil {
		log.Error("get history user, username :  ",fmt.Sprintf("%s", username), err.Error())
		ctx.JSON(http.StatusOK, dto.ErrorResponse(err.Error(), utils.ERR_VALIDATE_DATA))
		return
	}
	start := (reqData.Page * reqData.MaxDataDisplay) - reqData.MaxDataDisplay
	data := dto.ResPaginationDto{
		BaseUrl: "",
		TotalData: u.mapUserFpmRepo.GetTotalUserFpm(reqData.Search),
		Record:    u.mapUserFpmRepo.GetHistoryUserFpm(start, reqData.MaxDataDisplay, reqData.Search),
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse(data, "", utils.SUCCESS_CODE))
}
/*
	Method delete user
*/
func (u userFpmController) DeleteUserFpm(ctx *gin.Context) {
	username, err := utils.GetSession(ctx, "username")
	if err != nil {
		log.Error("delete user fpm, error : ", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}

	var reqDeleteUser dto.ReqDeleteUserFpmDto
	errScan := ctx.ShouldBindJSON(&reqDeleteUser)

	if errScan != nil {
		log.Error("delete user fpm, username : ", fmt.Sprintf("%s", username), err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse(err.Error(), utils.ERR_VALIDATE_DATA))
		return
	}

	data := u.mapUserFpmRepo.CheckUserById(reqDeleteUser.IdUser)

	if data == (models.UserFpmModels{}) {
		log.Error("delete user fpm, username : ",fmt.Sprintf("%s",username)+", error :", "data user not found")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("data user not found", utils.ERR_VALIDATE_DATA))
		return
	}

	data.StsActive = "0"
	errSave := u.mapUserFpmRepo.DeleteUserFpm(data)
	if errSave != nil {
		log.Error("delete user fpm, username : ",fmt.Sprintf("%s",username)+", error :", errSave.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("somenthing wrong", utils.ERR_BAD_REQUEST))
		return
	}
	log.Error("delete user fpm, username : ",fmt.Sprintf("%s",username)+", success delete , user : ", data)
	ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "Success delete user fpm", utils.SUCCESS_CODE))
}
/*
   Method add user
*/
func (u userFpmController) AddUserFpm(ctx *gin.Context) {
	username, err := utils.GetSession(ctx, "username")
	if err != nil {
		log.Error("add user fpm, error : ", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse("unauthorized", utils.ERR_AUTHORIZATION))
		return
	}
	var reqAddUser dto.ReqAddUserFpmDto
	errScan := ctx.ShouldBindJSON(&reqAddUser)
	if errScan != nil {
		log.Error("add user fpm, username : ",fmt.Sprintf("%s", username)+", error : ", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse(err.Error(), utils.ERR_VALIDATE_DATA))
		return
	}

	checkRoleUser := u.mapUserFpmRepo.CheckRoleUserExist(reqAddUser.RoleUser)

	if checkRoleUser.Id == 0 {
		log.Error("add user fpm, username : ",fmt.Sprintf("%s", username)+", error : ", "role user not found")
		ctx.JSON(http.StatusOK, dto.ErrorResponse("role user not found", utils.ERR_VALIDATE_DATA))
		return
	}

	dataUser := u.mapUserFpmRepo.CheckUserByUsername(reqAddUser.Username)

	if dataUser.ID > 0 {
		log.Error("add user fpm, username : ",fmt.Sprintf("%s", username)+", error : ", "username already exist at db")
		ctx.JSON(http.StatusOK, dto.ErrorResponse("username already exist at db", utils.ERR_VALIDATE_DATA))
		return
	}
	idRole := strconv.Itoa(checkRoleUser.Id)
	dataSave := models.UserFpmModels{
		Username:   reqAddUser.Username,
		IdRoleUser: idRole,
		StsActive:  "1",
	}

	errSave := u.mapUserFpmRepo.SaveUserFpm(dataSave)

	if errSave != nil {
		log.Error("add user fpm, username : ",fmt.Sprintf("%s", username)+", error : ", errSave.Error())
		ctx.JSON(http.StatusOK, dto.ErrorResponse("Failed save user fpm", utils.ERR_VALIDATE_DATA))
		return
	}
	log.Info("add user fpm, username : ",fmt.Sprintf("%s", username)+", success add user, with user : ", dataSave.Username)
	ctx.JSON(http.StatusOK, dto.SuccessResponse(nil, "Success save user fpm", utils.SUCCESS_CODE))

}

func InstanceUserFpmController(mapUserRepo repository.UserFpmRepository) UserFpmController {
	return &userFpmController{
		mapUserFpmRepo: mapUserRepo,
	}
}
