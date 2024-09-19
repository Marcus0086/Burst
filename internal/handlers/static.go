package handlers

import (
	"fmt"
	"html/template"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func serveStaticContent(writer http.ResponseWriter, body string) {

	content := template.Must(template.New("static_content").Parse(body))
	content.Execute(writer, nil)

}


func serveStaticFile(writer http.ResponseWriter, request *http.Request, root, uri string) {
	// Clean the URI to prevent directory traversal attacks
	cleanPath := filepath.Clean(uri)
	fmt.Println("Clean path:", cleanPath)

	if strings.Contains(cleanPath, "..") || strings.HasPrefix(cleanPath, ".") || strings.Contains(cleanPath, "/.") {
		http.NotFound(writer, request)
		return
	}

	fullPath := filepath.Join(root, cleanPath)
	fmt.Println("Full path:", fullPath)


	// Check if the path is a directory, if so, serve index.html from that directory
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(writer, request)
		} else {
			http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// If the path is a directory, try to serve index.html from that directory
	if info.IsDir() {
		indexPath := filepath.Join(fullPath, "index.html")
		indexInfo, err := os.Stat(indexPath)
		if err == nil && !indexInfo.IsDir() {
			serveFile(writer, request, indexPath, indexInfo)
			return
		}
		// Directory exists but no index.html, return 404
		http.NotFound(writer, request)
		return
	}

	// If the request doesn't have an extension (like /about), try adding .html
	if filepath.Ext(fullPath) == "" {
		htmlPath := fullPath + ".html"
		htmlInfo, err := os.Stat(htmlPath)
		if err == nil && !htmlInfo.IsDir() {
			serveFile(writer, request, htmlPath, htmlInfo)
			return
		}
	}

	// Serve the file if it exists and is not a directory
	if !info.IsDir() {
		serveFile(writer, request, fullPath, info)
		return
	}

	http.NotFound(writer, request)
}

func serveFile(writer http.ResponseWriter, request *http.Request, path string, fileInfo fs.FileInfo) {
	file, err := os.Open(path)
	if err != nil {
		http.NotFound(writer, request)
		return
	}
	defer file.Close()

	contentType := mime.TypeByExtension(filepath.Ext(path))
	if contentType == "" {
		contentType = "text/plain"
	}

	writer.Header().Set("Content-Type", contentType)
	http.ServeContent(writer, request, path, fileInfo.ModTime(), file)
}
