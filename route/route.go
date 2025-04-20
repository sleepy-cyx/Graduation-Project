package route

import (
	"Graduation-Project/handler"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.POST("/login", handler.Login)
	// 获取日程列表
	r.GET("/get_schedules", handler.GetScheduleInfo)
	r.POST("/register", handler.Register)
	// 更新日程(使用POST保持现有风格)
	r.POST("/update_schedule", handler.UpdateScheduleInfo)
	// 删除日程(更推荐RESTful风格，但根据现有示例保持POST)
	r.POST("/delete_schedule", handler.DeleteSchedule)
	// 原有其他接口
	r.POST("/translate_schedule", handler.TranslateSchedule)
}
