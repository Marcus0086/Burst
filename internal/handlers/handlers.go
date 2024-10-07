package handlers

import (
	"Burst/pkg/models"
	"fmt"
	"net/http"
	"strings"
)

func HandleConnection(writer http.ResponseWriter, request *http.Request, config *models.Config, unifiedConfig *models.ConfigJSON) {
	path := request.URL.Path
	method := request.Method
	headers := request.Header

	fmt.Printf("Method: %s\nURI: %s\nHeaders: %v\n", method, path, headers)

	var matchedRoute *models.RouteConfig
	for _, route := range config.Server.Routes {
		if strings.HasSuffix(route.Path, "*") {
			basePath := strings.TrimSuffix(route.Path, "*")
			fmt.Println("Base path:", basePath)
			if strings.HasPrefix(path, basePath) {
				matchedRoute = &route
				break
			}
		} else if route.Path == path {
			matchedRoute = &route
			break
		}
	}

	if matchedRoute == nil {
		http.NotFound(writer, request)
		return
	}

	switch matchedRoute.Handler {
	case models.StaticHandler:
		serveStaticFile(writer, request, config.Server.Root, path)
	case models.StaticContentHandler:
		serveStaticContent(writer, matchedRoute.Body)
	case models.DynamicHandler:
		renderTemplate(writer, request, config.Server.Root, matchedRoute)
	case models.ReverseProxyHandler:
		proxyHandler(writer, request, matchedRoute)
	case models.AdminAPIHandler:
		AdminAPIHandler(writer, request, unifiedConfig, matchedRoute)
	default:
		http.Error(writer, "501 Not Implemented", http.StatusNotImplemented)
	}

}
