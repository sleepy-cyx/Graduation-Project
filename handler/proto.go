package handler

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResp struct {
	UserId uint   `json:"user_id"`
	Msg    string `json:"msg"`
}

type Schedule struct {
	Id        uint   `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint   `gorm:"not null" json:"user_id"`
	Title     string `gorm:"not null" json:"title"`
	Comment   string `gorm:"not null" json:"comment"`
	StartTime string `gorm:"not null" json:"start_time"`
	EndTime   string `gorm:"not null" json:"end_time"`
}

type GetScheduleInfoResp struct {
	Schedules []Schedule `json:"schedules"`
	Msg       string     `json:"msg"`
	ErrCode   int        `json:"errcode"`
}

type UpdateScheduleInfoReq struct {
	Schedules []Schedule `json:"schedules"`
}

type UpdateScheduleInfoResp struct {
	Msg     string `json:"msg"`
	ErrCode int    `json:"errcode"`
}
