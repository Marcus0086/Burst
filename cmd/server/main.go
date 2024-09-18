package main

import (
	"Burst/internal/config"
	"Burst/internal/worker"
	"Burst/pkg/models"
	"log"
	"sync"
)

const maxWorkers = 10

func main() {

	configs, err := config.LoadConfigs("./config")
	if err != nil {
		log.Fatal(err)
	}
	workers := configs[0].Global.Workers
	if workers <= 0 {
		workers = maxWorkers
	}

	jobs := make(chan *models.Config, workers)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker.Worker(jobs, &wg)
	}

	for _, config := range configs {
		jobs <- config
	}
	close(jobs)
	wg.Wait()
}
