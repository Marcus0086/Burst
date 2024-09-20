package loadbalancer

import (
	"Burst/pkg/models"
	"fmt"
	"net/url"
	"sync"
	"time"
)

type LoadBalancer interface {
	NextTarget() (*Backend, error)
	Backends() []*Backend
}

func CreateLoadBalancer(route *models.RouteConfig) (LoadBalancer, error) {
	var targets []*Backend
	if (route.Target != "") {
		u, err := url.Parse(route.Target)
		if err != nil {
			return nil, fmt.Errorf("invalid target URL %s: %v", route.Target, err)
		}
		backend := &Backend{URL: u, Alive: false}
		targets = append(targets, backend)
	}
	for _, target := range route.Targets {
		u, err := url.Parse(target)
		if err != nil {
			return nil, fmt.Errorf("invalid target URL %s: %v", target, err)
		}
		backend := &Backend{URL: u, Alive: false}
		targets = append(targets, backend)
	}

	switch route.LoadBalancing {
	case "round_robin":
		lb, err := roundRobin(targets)
		if err != nil {
			return nil, err
		}
		return lb, nil
	default:
		lb, err := roundRobin(targets)
		if err != nil {
			return nil, err
		}
		return lb, nil
	}
}

func roundRobin(targets []*Backend) (LoadBalancer, error) {
	lb := &RoundRobin{Targets: targets}
	performInitialHealthCheck(targets)
	go healthCheck(targets, 10*time.Second)
	return lb, nil
}


func performInitialHealthCheck(backends []*Backend) {
    var wg sync.WaitGroup
    for _, b := range backends {
        wg.Add(1)
        go func(b *Backend) {
            defer wg.Done()
            alive := isBackendAlive(b.URL)
            b.SetAlive(alive)
        }(b)
    }
    wg.Wait()
}