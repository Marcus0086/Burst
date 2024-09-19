package handlers

import (
	"Burst/internal/utils"
	"Burst/pkg/models"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func proxyHandler(writer http.ResponseWriter, request *http.Request, route *models.RouteConfig) {
    var targetURL *url.URL
    var err error
    if route.LoadBalancing != "" && route.Targets != nil {
        loadBalancer, err := utils.GetLoadBalancer(route)
        if err != nil {
            http.Error(writer, "Failed to create load balancer", http.StatusInternalServerError)
            return
        }
        targetURL = loadBalancer.NextTarget()
    } else {
        targetURL, err = url.Parse(route.Target)
        if err != nil {
            http.Error(writer, "Invalid proxy target URL", http.StatusInternalServerError)
            return
        }
    }

    fmt.Println("Target URL is:", targetURL)
    proxy := httputil.NewSingleHostReverseProxy(targetURL)
    originalDirector := proxy.Director

    proxy.Director = func(req *http.Request) {
        originalDirector(req)

        // Compute base path by removing "/*" from route.Path if present
        basePath := strings.TrimSuffix(route.Path, "/*")
        // Ensure basePath does not end with a slash
        basePath = strings.TrimSuffix(basePath, "/")

        // Log the base path and original request path for debugging
        log.Println("Base path:", basePath)
        log.Println("Original request path:", req.URL.Path)

        if strings.HasPrefix(req.URL.Path, basePath) {
            // Remove the base path from the request URL path
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

        // Log the adjusted request path
        log.Println("Adjusted request path:", req.URL.Path)

        // Set headers
        req.Header.Set("X-Forwarded-For", request.RemoteAddr)
        req.Header.Set("X-Forwarded-Host", request.Host)
        req.Header.Set("X-Forwarded-Proto", request.URL.Scheme)
    }

    proxy.ModifyResponse = func(resp *http.Response) error {
        // Remove sensitive headers
        resp.Header.Del("Server")
        return nil
    }

    proxy.Transport = &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
        log.Printf("Proxy error: %v", err)
        http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
    }

    proxy.ServeHTTP(writer, request)
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
