package handlers

import (
	"net/http"
	"os"
)

func ServeFilesHandler(w http.ResponseWriter, r *http.Request) {
	_, err := os.ReadFile("../frontend" + r.URL.Path)
	if err != nil && r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "../frontend/index.html")
		return
	}
	http.FileServer(http.Dir("../frontend")).ServeHTTP(w, r)
}
