package handler

import "Graduation-Project/model"

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResp struct {
	UserId uint   `json:"user_id"`
	Msg    string `json:"msg"`
}

// 注册请求结构体
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetScheduleInfoResp struct {
	Schedules []model.Schedule `json:"schedules"`
	Msg       string           `json:"msg"`
	ErrCode   int              `json:"errcode"`
}

type UpdateScheduleInfoReq struct {
	Schedules []model.Schedule `json:"schedules"`
}

type UpdateScheduleInfoResp struct {
	Msg     string `json:"msg"`
	ErrCode int    `json:"errcode"`
}
type UpdateScheduleReq struct {
	ScheduleID uint   `json:"schedule_id" binding:"required"`
	UserID     uint   `json:"user_id" binding:"required"` // 可选，根据实际需求
	Title      string `json:"title"`
	Location   string `json:"location"`
	Comment    string `json:"comment"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}
