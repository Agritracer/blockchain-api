package server

import (
	"fmt"
	"log"
	"net/http"

	"agritrace-api/internal/config"
	"agritrace-api/internal/handler"
	"agritrace-api/internal/middleware"
	"agritrace-api/internal/utils"
)

func Start() {
	config.LoadConfig()

	http.Handle("/submit", middleware.Auth(http.HandlerFunc(handler.HandleSubmit)))
	http.Handle("/trace", middleware.Auth(http.HandlerFunc(handler.HandleTrace)))
	http.Handle("/query", middleware.Auth(http.HandlerFunc(handler.HandleTraceByID)))

	port := config.Cfg.Port
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server is running at http://" + utils.GetPublicIP() + ":" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
