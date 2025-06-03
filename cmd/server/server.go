package server

import (
	"fmt"
	"log"
	"net/http"

	"agritrace-api/internal/config"
	"agritrace-api/internal/handler"
	"agritrace-api/internal/utils"
)

func Start() {
	config.LoadConfig()

	http.HandleFunc("/submit", handler.HandleSubmit)
	http.HandleFunc("/trace", handler.HandleTrace)
	// http.HandleFunc("/query", handler.HandleQuery)
	http.HandleFunc("/query", handler.HandleTraceByID)

	port := config.Cfg.Port
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server is running at http://" + utils.GetPublicIP() + ":" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
