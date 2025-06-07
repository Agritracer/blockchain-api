package web

import (
	"log"
	"net/http"

	"agritrace-api/internal/web/controller"
)

func Start() error {
	http.HandleFunc("/", controller.HomeHandler)

	log.Println("🌐 Web server running at http://localhost:8080")
	return http.ListenAndServe(":8081", nil)
}
