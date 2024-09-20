package loadbalancer

import (
	"Burst/pkg/models"
	"fmt"
	"net/url"
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




