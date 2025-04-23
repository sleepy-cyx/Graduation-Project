package main

import (
	"Graduation-Project/config"
	"Graduation-Project/log"
	"Graduation-Project/route"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	config.InitConfig()
	router.Use(gin.Recovery())
	//pprof.Register(router)
	route.InitRouter(router)
	err := router.Run(":8000")
	if err != nil {
		log.Logger.Errorf("httpsever start failed : %v", err)
	}
}
