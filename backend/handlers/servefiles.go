package handlers

import (
	"net/http"
	"os"
)

// ServeFilesHandler serves static files and handles SPA routing
func ServeFilesHandler(w http.ResponseWriter, r *http.Request) {
	// Try to read the requested file
	_, err := os.ReadFile("../frontend" + r.URL.Path)
	if err != nil {
		// If file doesn't exist, serve index.html for SPA routing
		file, err := os.ReadFile("../frontend/index.html")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// Return 404 status but serve index.html for client-side routing
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/html")
		w.Write(file)
		return
	}

	// Serve the requested file
	http.FileServer(http.Dir("../frontend")).ServeHTTP(w, r)
}
