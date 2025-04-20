package handler

import (
	"Graduation-Project/common"
	"Graduation-Project/log"
	"Graduation-Project/service"
	"Graduation-Project/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	resp := LoginResp{
		UserId: userInfo.Id,
	}
	Response(c, http.StatusOK, common.SUCCESS, resp)

}

func GetScheduleInfo(c *gin.Context) {
	// 1. 参数获取与验证
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		log.Logger.Errorf("Missing user_id parameter")
		Response(c, common.ERRCODE_PARAMETER_INVALID, "缺少user_id参数", nil)
		return
	}

	// 转换user_id为uint
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		log.Logger.Errorf("Invalid user_id format: %v", err)
		Response(c, common.ERRCODE_PARAMETER_INVALID, "非法的用户ID格式", nil)
		return
	}

	// 2. 创建服务处理器
	serviceHandler, err := service.NewUserServiceHandler() // 假设存在对应服务处理器
	if err != nil {
		log.Logger.Errorf("创建服务处理器失败: %v", err)
		Response(c, common.ERRCODE_SERVER_ERROR, common.SERVER_ERROR, nil)
		return
	}

	// 3. 查询数据库
	schedules, _, err := serviceHandler.GetSchedulesByUserID(uint(userID))
	if err != nil {
		log.Logger.Errorf("查询日程失败: %v", err)
		Response(c, common.ERRCODE_SERVER_ERROR, common.SERVER_ERROR, nil)
		return
	}

	// 4. 构造响应
	resp := GetScheduleInfoResp{
		Schedules: schedules,
		Msg:       common.SUCCESS,
		ErrCode:   http.StatusOK,
	}

	// 5. 返回响应
	Response(c, http.StatusOK, common.SUCCESS, resp)
}
func UpdateScheduleInfo(c *gin.Context) {}

func TranslateSchedule(c *gin.Context) {}
