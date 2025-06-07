package web

import (
	"log"
	"net/http"

	"agritrace/internal/config"
	"agritrace/internal/utils"
	"agritrace/internal/web/controller"
	"agritrace/internal/web/middleware"
)

func Start() error {
	config.LoadConfig()

	mux := http.NewServeMux()

	mux.Handle("/login", middleware.RedirectIfLoggedIn(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controller.LoginFormHandler(w, r)
		} else if r.Method == http.MethodPost {
			controller.LoginHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/", middleware.JWTAuthRedirect(http.HandlerFunc(controller.HomeHandler)))

	mux.Handle("/logout", middleware.JWTAuthRedirect(http.HandlerFunc(controller.LogoutHandler)))

	port := config.Cfg.WebPort
	if port == "" {
		port = "8081"
	}

	log.Println(utils.GetPublicIP() + ":" + port)
	return http.ListenAndServe(":"+port, mux)
}
