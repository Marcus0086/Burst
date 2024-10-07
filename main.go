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

const maxWorkers = 100

func main() {
	var unifiedConfig models.ConfigJSON
	allConfigs := make([]*models.Config, 0)
	rootConfig, err := config.LoadRootConfig()
	if err != nil {
		log.Fatal(err)
	}
	unifiedConfig = mergeConfig(unifiedConfig, rootConfig)
	allConfigs = append(allConfigs, rootConfig)
	if _, err := os.Stat("./config"); !os.IsNotExist(err) {
		configs, err := config.LoadConfigs("./config")
		if err != nil {
			log.Default().Println(err)
		} else {
			allConfigs = append(allConfigs, configs...)
			for _, c := range configs {
				unifiedConfig = mergeConfig(unifiedConfig, c)
			}
		}
	}

	if rootConfig.Global.AdminAPI {
		adminConfig := initializeAdminAPIConfig()
		allConfigs = append(allConfigs, adminConfig)
	}

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
		go worker.Worker(ctx, &wg, &unifiedConfig)
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

func mergeConfig(unifiedConfig models.ConfigJSON, newConfig *models.Config) models.ConfigJSON {
	server := models.ServerConfig{
		Listen: newConfig.Server.Listen,
		Routes: newConfig.Server.Routes,
	}
	if unifiedConfig.Apps.HTTP.Servers == nil {
		unifiedConfig.Apps.HTTP.Servers = make(map[string]models.ServerConfig)
	}
	unifiedConfig.Apps.HTTP.Servers[newConfig.Server.Listen] = server

	return unifiedConfig
}

func initializeAdminAPIConfig() *models.Config {
	adminRoutes := make([]models.RouteConfig, 0)
	adminRoutes = append(adminRoutes, models.RouteConfig{
		Path:    "/config",
		Method:  "GET",
		Handler: "admin_api",
	})
	adminRoutes = append(adminRoutes, models.RouteConfig{
		Path:    "/load",
		Method:  "POST",
		Handler: "admin_api",
	})
	adminServerConfig := models.ServerConfig{
		Listen: ":3001",
		Routes: adminRoutes,
	}
	adminConfig := models.Config{
		Server: adminServerConfig,
	}
	return &adminConfig
}
