package handler

import (
	"Graduation-Project/common"
	"Graduation-Project/log"
	"Graduation-Project/model"
	"Graduation-Project/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
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
		Response(c, http.StatusOK, common.USERNAME_OR_PASSWORD_INCORRECT, nil)
		return
	}

	userId, err := serviceHandler.Login(req.Username, req.Password)
	if err != nil {
		log.Logger.Error(err)
		Response(c, http.StatusOK, common.USERNAME_OR_PASSWORD_INCORRECT, nil)
		return
	}
	resp := LoginResp{
		UserId: userId,
	}
	Response(c, http.StatusOK, common.SUCCESS, resp)

}

func Register(c *gin.Context) {
	// 1. 参数获取与验证
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Logger.Errorf("注册参数绑定失败: %v", err)
		Response(c, common.ERRCODE_PARAMETER_INVALID, "请求参数格式错误", nil)
		return
	}

	// 空值校验
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		log.Logger.Errorf("用户名或密码为空")
		Response(c, common.ERRCODE_PARAMETER_INVALID, "用户名和密码不能为空", nil)
		return
	}

	// 2. 创建服务处理器
	serviceHandler, err := service.NewUserServiceHandler()
	if err != nil {
		log.Logger.Errorf("创建用户服务处理器失败: %v", err)
		Response(c, common.ERRCODE_SERVER_ERROR, common.SERVER_ERROR, nil)
		return
	}

	err = serviceHandler.CreateUser(req.Username, req.Password)
	if err != nil {
		log.Logger.Errorf("用户注册失败: %v", err)
		// 处理特定错误类型
		if strings.Contains(err.Error(), "已存在") {
			Response(c, common.ERRCODE_SERVER_ERROR, "用户名已存在", nil)
		} else {
			Response(c, common.ERRCODE_SERVER_ERROR, common.SERVER_ERROR, nil)
		}
		return
	}
	// 5. 返回响应
	Response(c, http.StatusCreated, common.SUCCESS, nil)
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
	}

	// 5. 返回响应
	Response(c, http.StatusOK, common.SUCCESS, resp)
}

// ------------------------ 更新日程 ------------------------
func UpdateScheduleInfo(c *gin.Context) {
	// 1. 参数绑定与校验
	var req UpdateScheduleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Logger.Errorf("请求参数解析失败: %v", err)
		Response(c, common.ERRCODE_PARAMETER_INVALID, "无效的请求格式", nil)
		return
	}

	// 2. 获取当前用户ID（假设从Token或Session中获取）
	currentUserID, exists := c.Get("user_id")
	if !exists {
		log.Logger.Errorf("用户未登录")
		Response(c, common.ERRCODE_UNAUTHORIZED, "请先登录", nil)
		return
	}

	// 3. 创建服务处理器
	serviceHandler, err := service.NewUserServiceHandler()
	if err != nil {
		log.Logger.Errorf("创建服务处理器失败: %v", err)
		Response(c, common.ERRCODE_SERVER_ERROR, common.SERVER_ERROR, nil)
		return
	}

	// 4. 调用Service层
	schedule := &model.Schedule{
		Id:        req.ScheduleID,
		UserID:    req.UserID, // 需要验证是否与currentUserID一致
		Title:     req.Title,
		Location:  req.Location,
		Comment:   req.Comment,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	err = serviceHandler.UpdateSchedule(currentUserID.(uint), schedule)
	if err != nil {
		log.Logger.Errorf("更新日程失败: %v", err)
		// 根据错误类型返回不同提示
		if strings.Contains(err.Error(), "无权") {
			Response(c, common.ERRCODE_UNAUTHORIZED, "无权修改该日程", nil)
		} else {
			Response(c, common.ERRCODE_SERVER_ERROR, common.SERVER_ERROR, nil)
		}
		return
	}

	// 5. 返回成功
	Response(c, http.StatusOK, common.SUCCESS, nil)
}

// ------------------------ 删除日程 ------------------------
func DeleteSchedule(c *gin.Context) {
	// 1. 参数获取与验证
	scheduleIDStr := c.Param("schedule_id")
	if scheduleIDStr == "" {
		log.Logger.Errorf("Missing schedule_id parameter")
		Response(c, common.ERRCODE_PARAMETER_INVALID, "缺少日程ID参数", nil)
		return
	}

	scheduleID, err := strconv.ParseUint(scheduleIDStr, 10, 64)
	if err != nil {
		log.Logger.Errorf("Invalid schedule_id format: %v", err)
		Response(c, common.ERRCODE_PARAMETER_INVALID, "非法的日程ID格式", nil)
		return
	}

	// 2. 获取当前用户ID
	currentUserID, exists := c.Get("user_id")
	if !exists {
		log.Logger.Errorf("用户未登录")
		Response(c, common.ERRCODE_UNAUTHORIZED, "请先登录", nil)
		return
	}

	// 3. 创建服务处理器
	serviceHandler, err := service.NewUserServiceHandler()
	if err != nil {
		log.Logger.Errorf("创建服务处理器失败: %v", err)
		Response(c, common.ERRCODE_SERVER_ERROR, common.SERVER_ERROR, nil)
		return
	}

	// 4. 调用Service层
	err = serviceHandler.DeleteSchedule(currentUserID.(uint), uint(scheduleID))
	if err != nil {
		log.Logger.Errorf("删除日程失败: %v", err)
		if strings.Contains(err.Error(), "无权") {
			Response(c, common.ERRCODE_UNAUTHORIZED, "无权删除该日程", nil)
		} else {
			Response(c, common.ERRCODE_SERVER_ERROR, common.SERVER_ERROR, nil)
		}
		return
	}

	// 5. 返回成功
	Response(c, http.StatusOK, common.SUCCESS, nil)
}

func TranslateSchedule(c *gin.Context) {

}
