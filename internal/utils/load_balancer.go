package utils

import (
	loadbalancer "Burst/internal/handlers/load_balancer"
	"Burst/pkg/models"
	"sync"
)

var (
    loadBalancers = make(map[string]models.LoadBalancer)
    lbMutex       sync.Mutex
)


func GetLoadBalancer(route *models.RouteConfig) (models.LoadBalancer, error) {
    lbMutex.Lock()
    defer lbMutex.Unlock()

    key := route.Path

    if lb, exists := loadBalancers[key]; exists {
        return lb, nil
    }

	lb, err := loadbalancer.CreateLoadBalancer(route)
	if err != nil {
		return nil, err
	}

	loadBalancers[key] = lb

    return lb, nil
}
