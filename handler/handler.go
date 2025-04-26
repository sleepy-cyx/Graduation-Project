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
	"time"
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
		Response(c, common.ERRCODE_SUCCESS, common.USERNAME_OR_PASSWORD_INCORRECT, nil)
		return
	}

	userId, err := serviceHandler.Login(req.Username, req.Password)
	if err != nil {
		log.Logger.Error(err)
		Response(c, common.ERRCODE_SUCCESS, common.USERNAME_OR_PASSWORD_INCORRECT, nil)
		return
	}
	resp := LoginResp{
		UserId: userId,
	}
	Response(c, common.ERRCODE_SUCCESS, common.SUCCESS, resp)

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
	println(userIDStr)
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
	Response(c, common.ERRCODE_SUCCESS, common.SUCCESS, resp)
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
	Response(c, common.ERRCODE_SUCCESS, common.SUCCESS, nil)
}

// ------------------------ 删除日程 ------------------------
func DeleteSchedule(c *gin.Context) {
	// 2. 绑定并验证请求体
	var reqBody DeleteScheduleRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Logger.Errorf("请求参数错误: %v", err)
		Response(c, common.ERRCODE_PARAMETER_INVALID, "请求参数格式错误", nil)
		return
	}
	userId, err := strconv.Atoi(reqBody.UserID)
	if err != nil {
		log.Logger.Errorf("非法参数 user_id:%d schedule_id:%d", reqBody.UserID, reqBody.ScheduleID)
		Response(c, common.ERRCODE_PARAMETER_INVALID, "id非法", nil)
		return
	}
	// 3. 参数校验
	if userId <= 0 || reqBody.ScheduleID <= 0 {
		log.Logger.Errorf("非法参数 user_id:%d schedule_id:%d", reqBody.UserID, reqBody.ScheduleID)
		Response(c, common.ERRCODE_PARAMETER_INVALID, "ID必须为正整数", nil)
		return
	}

	// 4. 创建服务处理器
	serviceHandler, err := service.NewUserServiceHandler()
	if err != nil {
		log.Logger.Errorf("服务初始化失败: %v", err)
		Response(c, common.ERRCODE_SERVER_ERROR, common.SERVER_ERROR, nil)
		return
	}

	// 5. 调用Service层（需确保有权限校验）
	err = serviceHandler.DeleteSchedule(uint(userId), uint(reqBody.ScheduleID))
	if err != nil {
		log.Logger.Errorf("删除失败: %v", err)
		if strings.Contains(err.Error(), "无权") {
			Response(c, common.ERRCODE_UNAUTHORIZED, "无权操作该日程", nil)
		} else {
			Response(c, common.ERRCODE_SERVER_ERROR, common.SERVER_ERROR, nil)
		}
		return
	}

	// 6. 返回成功
	Response(c, common.ERRCODE_SUCCESS, common.SUCCESS, nil)
}

func TranslateSchedule(c *gin.Context) {

	// 2. 绑定请求参数
	var req TextToScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Response(c, common.ERRCODE_PARAMETER_INVALID, "请求参数错误", nil)
		return
	}

	// 3. 创建服务实例
	service, err := service.NewScheduleService()
	if err != nil {
		log.Logger.Errorf("创建服务失败: %v", err)
		Response(c, common.ERRCODE_SERVER_ERROR, common.SERVER_ERROR, nil)
		return
	}

	// 4. 调用服务处理
	schedule, err := service.TextToSchedule(req.Text)
	if err != nil {
		Response(c, common.ERRCODE_SERVER_ERROR, common.SERVER_ERROR, nil)
		return
	}

	// 5. 返回成功响应
	Response(c, common.ERRCODE_SUCCESS, common.SUCCESS, schedule)
}

// ------------------------ HTTP 处理层 ------------------------
func CreateSchedule(c *gin.Context) {

	// 绑定并校验参数
	var req CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数校验失败: " + err.Error(),
		})
		return
	}
	log.Logger.Infof("req:%v", req)
	// 3. 时间格式校验（ 示例：2024-05-01T14:00:00+08:00）
	if _, err := time.Parse(time.DateTime, req.StartTime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "开始时间格式无效，需使用 ISO8601 格式",
		})
		return
	}
	if _, err := time.Parse(time.DateTime, req.EndTime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "结束时间格式无效，需使用 ISO8601 格式",
		})
		return
	}

	// 4. 构造领域模型对象
	schedule := &model.Schedule{
		UserID:    req.UserID,
		Title:     strings.TrimSpace(req.Title),
		Location:  strings.TrimSpace(req.Location),
		Comment:   strings.TrimSpace(req.Comment),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// 5. 调用 Service 层
	serviceHandler, err := service.NewUserServiceHandler()
	if err != nil {
		log.Logger.Errorf("创建服务处理器失败: %v", err)
		Response(c, common.ERRCODE_SERVER_ERROR, common.SERVER_ERROR, nil)
		return
	}
	if err := serviceHandler.CreateSchedule(schedule); err != nil {
		// 6. 处理已知业务错误
		switch {
		case strings.Contains(err.Error(), "用户ID"):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case strings.Contains(err.Error(), "日程冲突"):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
		return
	}

	Response(c, common.ERRCODE_SUCCESS, common.SUCCESS, nil)
}
func ParseFile(c *gin.Context) {
	// 获取上传文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		Response(c, common.ERRCODE_PARAMETER_INVALID, "文件上传失败", nil)
		return
	}
	defer file.Close()

	// 验证文件类型
	allowedTypes := map[string]bool{
		"application/pdf": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       true,
	}
	if !allowedTypes[header.Header.Get("Content-Type")] {
		Response(c, common.ERRCODE_PARAMETER_INVALID, "不支持的文件类型", nil)
		return
	}

	// 调用服务层解析
	serviceHandler, err := service.NewScheduleService()
	if err != nil {
		// ...错误处理...
	}

	schedule, err := serviceHandler.FileToSchedule(file, header.Filename)
	if err != nil {
		// ...错误处理...
	}

	Response(c, common.ERRCODE_SUCCESS, common.SUCCESS, schedule)
}
