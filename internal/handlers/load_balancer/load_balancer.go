package loadbalancer

import (
	"Burst/pkg/models"
	"fmt"
	"net/url"
)

func CreateLoadBalancer(route *models.RouteConfig) (models.LoadBalancer, error) {
	var targets []*url.URL

	for _, target := range route.Targets {
		u, err := url.Parse(target)
		if err != nil {
			return nil, fmt.Errorf("invalid target URL %s: %v", target, err)
		}
		targets = append(targets, u)
	}

	switch route.LoadBalancing {
	case "round_robin":
		return &RoundRobin{Targets: targets}, nil
	default:
		return nil, fmt.Errorf("unknown load balancing method %s", route.LoadBalancing)
	}
}
