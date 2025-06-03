package server

import (
	"fmt"
	"log"
	"net/http"

	"agritrace-api/internal/handler"
	"agritrace-api/internal/utils"
)

func Start() {
	http.HandleFunc("/submit", handler.HandleSubmit)
	http.HandleFunc("/trace", handler.HandleTrace)

	port := "8080"
	fmt.Println("Server is running at http://" + utils.GetPublicIP() + ":" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
