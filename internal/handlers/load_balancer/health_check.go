package loadbalancer

import (
	"log"
	"net"
	"net/url"
	"sync"
	"time"
)

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