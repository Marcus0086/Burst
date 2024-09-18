package server

import (
	"Burst/internal/handlers"
	"Burst/internal/utils"
	"Burst/pkg/models"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/http2"
)

func StartServer(config *models.Config) {

	if config.Server.Listen == "" {
		config.Server.Listen = ":80"
	}

	addr := config.Server.Listen
	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.HandleConnection(w, r, config)
		}),
	}

	if strings.HasSuffix(addr, ":443") || config.Server.HTTPS {
		cert, err := utils.GenerateSelfSignedCert()
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
