package server

import (
	"fmt"

	"agritrace/internal/config"
	"agritrace/internal/handler"
	"agritrace/internal/middleware"
	"agritrace/internal/utils"

	"github.com/gin-gonic/gin"
)

func Start() error {
	config.LoadConfig()

	router := gin.Default()
	protected := router.Group("/api/", middleware.Auth())

	{
		protected.GET("/list", handler.HandleList)
		protected.POST("/submit", handler.HandleSubmit)
		protected.GET("/trace", handler.HandleTrace)
		protected.GET("/query", handler.HandleTraceByID)
		// protected.GET("/compare", hanlder.HandlerCompare)
	}

	port := config.Cfg.ApiPort
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server is running at http://" + utils.GetPublicIP() + ":" + port)
	return router.Run(":" + port)
}
