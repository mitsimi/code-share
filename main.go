package main

import (
	"codeShare/frontend"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

func init() {
	// Add correct MIME types for modern web files
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".mjs", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
}

func main() {
	// Create a file server handler for the embedded dist directory
	fs := http.FileServer(http.FS(frontend.DistDirFS))

	// Custom handler to handle Vue.js routing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Get the file extension
		ext := filepath.Ext(path)
		mimeType := mime.TypeByExtension(ext)

		// If it's a static asset (has a mime type), serve it directly
		if mimeType != "" {
			// Set the correct content type
			w.Header().Set("Content-Type", mimeType)
			// Create a new request with the modified path
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = path
			fs.ServeHTTP(w, r2)
			return
		}

		// For all other routes, serve index.html
		indexFile, err := frontend.DistDirFS.Open("index.html")
		if err != nil {
			http.Error(w, "Error loading index.html", http.StatusInternalServerError)
			return
		}
		defer indexFile.Close()

		http.ServeContent(w, r, "index.html", time.Now(), indexFile.(io.ReadSeeker))
	})

	// Start the server on port 8080
	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

