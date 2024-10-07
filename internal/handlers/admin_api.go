package handlers

import (
	"Burst/internal/config/middleware"
	"Burst/pkg/models"
	"encoding/json"
	"net/http"
	"sync"
)

var mu sync.Mutex

func AdminAPIHandler(writer http.ResponseWriter, request *http.Request, unifiedConfig *models.ConfigJSON, matchedRoute *models.RouteConfig) {
	middleware.ServerHeaderMiddleware(writer)
	switch request.Method {
	case http.MethodGet:
		if matchedRoute.Path == "/config" {
			getAdminAPI(writer, unifiedConfig)
		} else {
			http.Error(writer, "404 Not Found", http.StatusNotFound)
		}
	// case http.MethodPost:
	// 	postAdminAPI(writer, request, config, matchedRoute)
	// case http.MethodPut:
	// 	putAdminAPI(writer, request, config, matchedRoute)
	// case http.MethodDelete:
	// 	deleteAdminAPI(writer, request, config, matchedRoute)
	default:
		http.Error(writer, "501 Not Implemented", http.StatusNotImplemented)
	}
}

func getAdminAPI(writer http.ResponseWriter, unifiedConfig *models.ConfigJSON) {
	mu.Lock()
	defer mu.Unlock()
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(unifiedConfig)
}

func postAdminAPI(writer http.ResponseWriter, request *http.Request, config *models.Config) {
	mu.Lock()
	defer mu.Unlock()
	var newConfig models.ConfigJSON
	err := json.NewDecoder(request.Body).Decode(&newConfig)
	if err != nil {
		http.Error(writer, "400 Bad Request", http.StatusBadRequest)
		return
	}

}
