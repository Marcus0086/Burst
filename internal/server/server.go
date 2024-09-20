package server

import (
	"Burst/internal/handlers"
	"Burst/internal/utils"
	"Burst/pkg/models"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/http2"
)

func StartServer(ctx context.Context, config *models.Config) {

	if config.Server.Listen == "" {
		log.Fatal("Listen address is not set")
	}

	addr := config.Server.Listen
	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.HandleConnection(w, r, config)
		}),
	}

	errChan := make(chan error, 1)

	go func() {
		var err error
		if strings.HasSuffix(addr, ":443") || config.Server.HTTPS {
			cert, err := utils.GenerateSelfSignedCert()
			if err != nil {
				log.Printf("Error loading certificate and key for %s: %v", addr, err)
				errChan <- err
				return
			}
			server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
			http2.ConfigureServer(server, &http2.Server{})
			err = server.ListenAndServeTLS("", "")
			if err != nil && err != http.ErrServerClosed {
				log.Printf("Error starting server on %s: %v", addr, err)
			}
		} else {
			err = server.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				log.Printf("Error starting server on %s: %v", addr, err)
			}
			if err != nil && err != http.ErrServerClosed {
				log.Printf("Error starting server on %s: %v", addr, err)
				errChan <- err

			}
		}
		close(errChan)
	}()
	fmt.Println("Server started on", config.Server.Listen)
	select {
	case <-ctx.Done():
		log.Printf("Shutting down server on %s", addr)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server forced to shutdown: %v", err)
		}
		log.Printf("Server on %s gracefully shut down", addr)
	case err := <-errChan:
		if err != nil {
			log.Printf("Server on %s encountered an error: %v", addr, err)
		}
	}
}
