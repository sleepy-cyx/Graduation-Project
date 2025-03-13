package handler

import (
	"Graduation-Project/common"
	"Graduation-Project/log"
	"Graduation-Project/middleware"
	"Graduation-Project/service"
	"Graduation-Project/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(c *gin.Context, code int32, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"errcode": code,
		"msg":     msg,
		"data":    data,
	})
	return
}
func Login(c *gin.Context) {
	req := LoginReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Errorf("Bind err: %v", err)
		Response(c, common.ERRCODE_PARAMETER_INVALID, common.PARAMETER_INVALID, nil)
		return
	}
	serviceHandler, err := service.NewUserServiceHandler()
	if err != nil {
		log.Logger.Error(err)
	}

	//查询db
	userInfo, exists, err := serviceHandler.GetUserInfoByUsername(req.Username)
	if err != nil {
		log.Logger.Errorf("GetUserInfoByUsername err: %v", err)
	}
	if !exists {
		log.Logger.Errorf("userInfo not found")
	}
	userPwd := utils.Md5BySalt(req.Password, userInfo.Salt)
	if userPwd != userInfo.Password {
		log.Logger.Errorf("userPwd != userInfo.Password")
		Response(c, common.ERRCODE_USERNAME_OR_PASSWORD_INCORRECT, common.USERNAME_OR_PASSWORD_INCORRECT, nil)
		return
	}
	//生成token并返回
	token, err := middleware.GenerateToken(userInfo.Username, uint32(userInfo.Id))
	if err != nil {
		log.Logger.Errorf("GenerateToken err: %v", err)
		Response(c, common.ERRCODE_SERVER_ERROR, common.SERVER_ERROR, nil)
		return
	}
	resp := LoginResp{
		Token: token,
	}
	Response(c, http.StatusOK, common.SUCCESS, resp)

}

func Logout(c *gin.Context) {}

func GetScheduleInfo(c *gin.Context) {}

func UpdateScheduleInfo(c *gin.Context) {}
