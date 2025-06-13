package main

import (
	"agritrace/cmd/server"
	"agritrace/cmd/web"
	"log"
)

func main() {
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("API server error: %v", err)
		}
	}()

	if err := web.Start(); err != nil {
		log.Fatalf("Web UI server error: %v", err)
	}
}
