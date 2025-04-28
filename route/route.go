package route

import (
	"Graduation-Project/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.Use(cors.Default())
	r.POST("/login", handler.Login)
	// 获取日程列表
	r.GET("/get_schedules", handler.GetScheduleInfo)
	r.GET("/get_user_info", handler.GetUserInfo)
	r.POST("update_user_info", handler.UpdateUserInfo)
	r.POST("/register", handler.Register)
	// 更新日程(使用POST保持现有风格)
	r.POST("/update_schedule", handler.UpdateScheduleInfo)
	// 删除日程(更推荐RESTful风格，但根据现有示例保持POST)
	r.POST("/delete_schedule", handler.DeleteSchedule)
	// 原有其他接口
	r.POST("/translate_schedule", handler.TranslateSchedule)
	r.POST("/create_schedule", handler.CreateSchedule)
	r.POST("/parse_file", handler.ParseFile)
}
