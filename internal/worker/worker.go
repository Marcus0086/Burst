package worker

import (
	"Burst/internal/server"
	"Burst/pkg/models"
	"context"
	"log"
	"sync"
	"time"
)

var (
	jobs = make(chan *models.Config)
	jobsOnce sync.Once
)

func AddJob(config *models.Config) {
	jobs <- config
}

func CloseJobs() {	
	jobsOnce.Do(func() {
		close(jobs)
	})
}

func Worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			server.StartServer(ctx, job)
		case <-ctx.Done():
			return
		case <-time.After(time.Second * 5):
			log.Println("Worker still running")
			return
		}

	}
}