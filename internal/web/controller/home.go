package controller

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("internal/web/templates", "index.html"))
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"Messages": "Hello world from AgriTrace Web 👋",
	}

	tmpl.Execute(w, data)
}
