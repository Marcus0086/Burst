package handlers

import (
	"fmt"
	"net/http"
	"strings"

	configTypes "Burst/pkg/models"
)

func HandleConnection(writer http.ResponseWriter, request *http.Request, config *configTypes.Config) {
	path := request.URL.Path
	method := request.Method
	headers := request.Header

	fmt.Printf("Method: %s\nURI: %s\nHeaders: %v\n", method, path, headers)

	var matchedRoute *configTypes.RouteConfig
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
	case "static":
		serveStaticFile(writer, request, config.Server.Root, path)
	case "static_content":
		serveStaticContent(writer, matchedRoute.Body)
	case "dynamic":
		renderTemplate(writer, request, config.Server.Root, matchedRoute)
	case "reverse_proxy":
		proxyHandler(writer, request, matchedRoute)
	default:
		http.Error(writer, "501 Not Implemented", http.StatusNotImplemented)
	}

}
