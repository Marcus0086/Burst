package loadbalancer

import (
	"errors"
	"log"
	"net"
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

func healthCheck(backends []*Backend, interval time.Duration) {
	for {
		var wg sync.WaitGroup
		for _, backend := range backends {
			wg.Add(1)
			go func(backend *Backend) {
				defer wg.Done()
				alive := isBackendAlive(backend.URL)
				backend.SetAlive(alive)
			}(backend)
		}
		wg.Wait()
		time.Sleep(interval)
	}
}

func isBackendAlive(target *url.URL) bool {
	timeout := time.Duration(2 * time.Second)
	conn, err := net.DialTimeout("tcp", target.Host, timeout)
	if err != nil {
		log.Println("Error:", err)
		return false
	}
	defer conn.Close()
	return true
}