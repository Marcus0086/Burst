package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/fs"
	"log"
	"math/big"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/hcl"
	"golang.org/x/net/http2"
)

const maxWorkers = 10

func main() {

	configs, err := loadConfigs("./config")
	if err != nil {
		log.Fatal(err)
	}
	workers := configs[0].Global.Workers
	if workers <= 0 {
		workers = maxWorkers
	}

	jobs := make(chan *Config, workers)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(jobs, &wg)
	}

	for _, config := range configs {
		jobs <- config
	}
	close(jobs)
	wg.Wait()
}

func worker(jobs <-chan *Config, wg *sync.WaitGroup) {
	defer wg.Done()
	for config := range jobs {
		startServer(config)
	}
}

func startServer(config *Config) {

	if config.Server.Listen == "" {
		config.Server.Listen = ":80"
	}

	addr := config.Server.Listen
	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleConnection(w, r, config)
		}),
	}

	if strings.HasSuffix(addr, ":443") || config.Server.HTTPS {
		cert, err := generateSelfSignedCert()
		if err != nil {
			log.Printf("Error loading certificate and key for %s: %v", addr, err)
			return
		}
		server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
		http2.ConfigureServer(server, &http2.Server{})
		err = server.ListenAndServeTLS("", "")
		fmt.Println("Server started on", config.Server.Listen)
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting server on %s: %v", addr, err)
		}
	} else {
		err := server.ListenAndServe()
		fmt.Println("Server started on", config.Server.Listen)
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting server on %s: %v", addr, err)
		}
	}
	
}

func handleConnection(writer http.ResponseWriter, request *http.Request, config *Config) {
	path := request.URL.Path
	method := request.Method
	headers := request.Header

	fmt.Printf("Method: %s\nURI: %s\nHeaders: %v\n", method, path, headers)

	var matchedRoute *RouteConfig
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
	default:
		http.Error(writer, "501 Not Implemented", http.StatusNotImplemented)
	}

}

func serveStaticFile(writer http.ResponseWriter, request *http.Request, root, uri string) {
	cleanPath := filepath.Clean(uri)
	fmt.Println("Clean path:", cleanPath)
	if strings.Contains(cleanPath, "..") || strings.HasPrefix(cleanPath, ".") || strings.Contains(cleanPath, "/.") {
		http.NotFound(writer, request)
		return
	}

	fullPath := filepath.Join(root, cleanPath)

	fmt.Println("Full path:", fullPath)

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(writer, request)
		} else {
			http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if info.IsDir() {
		indexPath := filepath.Join(fullPath, "index.html")
		indexInfo, err := os.Stat(indexPath)
		if err == nil && !indexInfo.IsDir() {
			serveFile(writer, request, indexPath, indexInfo)
			return
		}
	}
	serveFile(writer, request, fullPath, info)

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

type Config struct {
	Global GlobalConfig `hcl:"global"`
	Server ServerConfig `hcl:"server"`
}

type GlobalConfig struct {
	Workers int `hcl:"workers"`
}

type ServerConfig struct {
	Listen     string            `hcl:"listen"`
	HTTPS      bool              `hcl:"https"`
	Root       string            `hcl:"root"`
	Routes     []RouteConfig     `hcl:"routes"`
	ErrorPages []ErrorPageConfig `hcl:"error_pages"`
	Middleware []string          `hcl:"middleware"`
}

type RouteConfig struct {
	Path       string   `hcl:"path"`
	Handler    string   `hcl:"handler"`
	Script     string   `hcl:"script"`
	Template   string   `hcl:"template"`
	Middleware []string `hcl:"middleware"`
}

type ErrorPageConfig struct {
	StatusCode int    `hcl:"status_code"`
	Path       string `hcl:"path"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = hcl.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}


func loadConfigs(path string) ([]*Config, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var configs []*Config

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".hcl" {
			continue
		}
		config, err := loadConfig(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}
const cacheDir = ".cache"

func generateSelfSignedCert() (tls.Certificate, error) {
	certFile := filepath.Join(cacheDir, "cert.pem")
	keyFile := filepath.Join(cacheDir, "key.pem")

	if _, err := os.Stat(certFile); err == nil {
		if _, err := os.Stat(keyFile); err == nil {
			return tls.LoadX509KeyPair(certFile, keyFile)
		}
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	tmpl := x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		Subject:      pkix.Name{CommonName: "localhost"},
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	certDer, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDer})
	privPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return tls.Certificate{}, err
	}

	if err := os.WriteFile(certFile, certPem, 0644); err != nil {
		return tls.Certificate{}, err
	}

	if err := os.WriteFile(keyFile, privPem, 0644); err != nil {
		return tls.Certificate{}, err
	}

	cert, err := tls.X509KeyPair(certPem, privPem)
	if err != nil {
		return tls.Certificate{}, err
	}

	return cert, nil
}
