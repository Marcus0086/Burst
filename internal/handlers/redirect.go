package handlers

import (
	"Burst/pkg/models"
	"net/http"
	"net/url"
	"strings"
)

func replaceWildcard(from string, to string, path string) string {
	basePath := strings.TrimSuffix(from, "/*")
	remainingPath := strings.TrimPrefix(path, basePath)
	return strings.Replace(to, "*", remainingPath, 1)
}

func redirectHandler(writer http.ResponseWriter, request *http.Request, route *models.RouteConfig) {
	if route.Redirects != nil {
		for _, redirect := range route.Redirects {
			if redirect.From == request.URL.Path {
				if _, err := url.ParseRequestURI(redirect.To); err != nil {
					http.Error(writer, "Invalid Redirect URL", http.StatusInternalServerError)
					return
				}
				// Perform the exact redirect
				http.Redirect(writer, request, redirect.To, redirect.Status)
				return
			}

			if strings.HasSuffix(redirect.From, "/*") {
				basePath := strings.TrimSuffix(redirect.From, "/*")
				if strings.HasPrefix(request.URL.Path, basePath) {
					remainingPath := strings.TrimPrefix(request.URL.Path, basePath)
					if remainingPath != "" {
						newURL := replaceWildcard(redirect.From, redirect.To, request.URL.Path)
						if _, err := url.ParseRequestURI(newURL); err != nil {
							http.Error(writer, "Invalid Redirect URL", http.StatusInternalServerError)
							return
						}
						http.Redirect(writer, request, newURL, redirect.Status)
						return
					}
				}
			}
		}
	}

}