package handlers

import (
	"Burst/internal/utils"
	"Burst/pkg/models"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

func proxyHandler(writer http.ResponseWriter, request *http.Request, route *models.RouteConfig) {
	utils.InitLoadBalancer(route)
	lb := utils.GetLoadBalancer(route)
	maxAttemps := len(lb.Backends()) * 2

	var lastErr error

	for attemps := 0; attemps < maxAttemps; attemps++ {
		backend, err := lb.NextTarget()
		if err != nil {
			http.Error(writer, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}
		targetURL := backend.URL
		fmt.Println("Target URL is:", targetURL)
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		
        originalDirector := proxy.Director
		ctx := request.Context()

		proxy.Director = func(req *http.Request) {
            originalDirector(req)
            req = req.WithContext(ctx)

            basePath := strings.TrimSuffix(route.Path, "/*")
            if basePath == "" && route.Path == "/*" {
                basePath = "/"
            }
            basePath = strings.TrimSuffix(basePath, "/")

            // Adjust the request URL path
            if strings.HasPrefix(req.URL.Path, basePath) {
                req.URL.Path = strings.TrimPrefix(req.URL.Path, basePath)
                if req.URL.Path == "" || !strings.HasPrefix(req.URL.Path, "/") {
                    req.URL.Path = "/" + req.URL.Path
                }
            }

            req.URL.RawPath = req.URL.Path

            // Prepend the target path if specified
            if targetURL.Path != "" {
                req.URL.Path = singleJoiningSlash(targetURL.Path, req.URL.Path)
                req.URL.RawPath = req.URL.Path
            }

            // Set headers
            req.Header.Set("X-Forwarded-For", request.RemoteAddr)
            req.Header.Set("X-Forwarded-Host", request.Host)
            req.Header.Set("X-Forwarded-Proto", request.URL.Scheme)
        }

		proxy.ModifyResponse = func(resp *http.Response) error {
			// Remove sensitive headers
			resp.Header.Set("Server", "Burst")
			return nil
		}

		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		var proxyError error
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error: %v", err)
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			proxyError = err
		}

		proxy.ServeHTTP(writer, request)
		if proxyError != nil {
			backend.SetAlive(false)
			lastErr = proxyError
			continue
		} else {
			return
		}
	}

	log.Printf("All backends failed. Last error: %v", lastErr)
	http.Error(writer, "Service Unavailable", http.StatusServiceUnavailable)

}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
