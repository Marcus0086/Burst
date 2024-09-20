package loadbalancer

import (
	"errors"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	URL *url.URL
	Alive bool
	mutex sync.RWMutex
}

func (b *Backend) IsAlive() bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.Alive
}

func (b *Backend) SetAlive(alive bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.Alive = alive
}


type RoundRobin struct {
	Targets []*Backend
	mutex   sync.Mutex
	index   int
}

func (rr *RoundRobin) NextTarget() (*Backend, error) {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()

	totalTargets := len(rr.Targets)
	startingIndex := rr.index

	for {
		target := rr.Targets[rr.index%totalTargets]
		rr.index++
		if target.IsAlive() {
			return target, nil
		}
		if rr.index%totalTargets == startingIndex%totalTargets {
			break
		}
	}
	return nil, errors.New("no available targets")
}

func (rr *RoundRobin) Backends() []*Backend {
	return rr.Targets
}

func roundRobin(targets []*Backend) (LoadBalancer, error) {
	lb := &RoundRobin{Targets: targets}
	performInitialHealthCheck(targets)
	go healthCheck(targets, 10*time.Second)
	return lb, nil
}