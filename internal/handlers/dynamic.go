package handlers

import (
	"Burst/internal/config/middleware"
	"Burst/internal/utils"
	"Burst/pkg/models"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func renderTemplate(writer http.ResponseWriter, request *http.Request, root string, route *models.RouteConfig) {
	middleware.ServerHeaderMiddleware(writer)
	utils.InitDynamicCache()
	cacheKey := request.URL.Path + "?" + request.URL.RawQuery

	if cachedContent, ok := utils.DynamicCache.Get(cacheKey); ok {
		content := cachedContent.([]byte)
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
		writer.WriteHeader(http.StatusOK)
		_, err := writer.Write(content)
		if err != nil {
			log.Printf("Error writing cached response: %v", err)
		}
		return
	}

	cleanPath := filepath.Clean(route.Template)
	if strings.Contains(cleanPath, "..") || strings.HasPrefix(cleanPath, ".") || strings.Contains(cleanPath, "/.") {
		http.NotFound(writer, request)
		return
	}

	fullPath := filepath.Join(root, cleanPath)
	fmt.Println("Full path:", fullPath)
	tmpl, err := template.ParseFiles(fullPath)
	if err != nil {
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Config":  route,
		"Request": request,
		"Data":    route.Data,
		"Time":    time.Now().Format("2006-01-02 15:04:05"),
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Printf("Error executing template %s: %v", fullPath, err)
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	renderedContent := buf.Bytes()
	utils.DynamicCache.Add(cacheKey, renderedContent)
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(renderedContent)))
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(renderedContent)
	if err != nil {
		log.Printf("Error writing rendered response: %v", err)
	}
}
