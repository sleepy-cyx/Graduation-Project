package route

import (
	"Graduation-Project/handler"
	"Graduation-Project/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.POST("/login", handler.Login)
	r.Use(middleware.JWTAuthMiddleware())
	r.POST("/logout", handler.Logout)
	r.GET("/get_schedule_info", handler.GetScheduleInfo)
	r.POST("/update_schedule_info", handler.UpdateScheduleInfo)
}
