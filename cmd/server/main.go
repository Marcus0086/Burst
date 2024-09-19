package main

import (
	"Burst/internal/config"
	"Burst/internal/worker"
	"Burst/pkg/models"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const maxWorkers = 10

func main() {
	allConfigs := make([]*models.Config, 0)
	rootConfig, err := config.LoadRootConfig()
	if err != nil {
		log.Fatal(err)
	}
	allConfigs = append(allConfigs, rootConfig)
	configs, err := config.LoadConfigs("./config")
	if err != nil {
		log.Fatal(err)
	}
	allConfigs = append(allConfigs, configs...)
	workers := rootConfig.Global.Workers
	if workers <= 0 {
		workers = maxWorkers
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker.Worker(ctx, &wg)
	}

	go func() {
		for _, config := range allConfigs {
			worker.AddJob(config)
		}
		worker.CloseJobs()
	}()

	<-stop
	cancel()
	wg.Wait()
}
