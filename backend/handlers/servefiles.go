package handlers

import (
	"net/http"
	"os"
)

func ServeFilesHandler(w http.ResponseWriter, r *http.Request) {
	_, err := os.ReadFile("../frontend" + r.URL.Path)
	if err != nil {
		file, err := os.ReadFile("../frontend/index.html")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// serve file with header
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/html")
		w.Write(file)
		return
	}

	http.FileServer(http.Dir("../frontend")).ServeHTTP(w, r)
}
