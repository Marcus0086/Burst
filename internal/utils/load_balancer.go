package utils

import (
	loadbalancer "Burst/internal/handlers/load_balancer"
	"Burst/pkg/models"
	"log"
	"sync"
)

var (
    loadBalancers = make(map[string]loadbalancer.LoadBalancer)
    lbMutex       sync.RWMutex
)


func InitLoadBalancer(route *models.RouteConfig) {
    lbMutex.Lock()
    defer lbMutex.Unlock()

    if _, exists := loadBalancers[route.Path]; exists {
        return // Load balancer already initialized for this route
    }

    lb, err := loadbalancer.CreateLoadBalancer(route)
    if err != nil {
        log.Fatalf("Error initializing load balancer: %v", err)
    }

    loadBalancers[route.Path] = lb
}

func GetLoadBalancer(route *models.RouteConfig) loadbalancer.LoadBalancer {
    lbMutex.RLock()
    defer lbMutex.RUnlock()

    if lb, ok := loadBalancers[route.Path]; ok {
        return lb
    }

    return nil
}
