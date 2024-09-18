package worker

import (
	"Burst/internal/server"
	"Burst/pkg/models"
	"sync"
)

func Worker(jobs <-chan *models.Config, wg *sync.WaitGroup) {
	defer wg.Done()
	for config := range jobs {
		server.StartServer(config)
	}
}