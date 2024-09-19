package loadbalancer

import (
	"net/url"
	"sync"
)

type RoundRobin struct {
	Targets []*url.URL
	mutex   sync.Mutex
	index   int
}

func (rr *RoundRobin) NextTarget() *url.URL {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()

	target := rr.Targets[rr.index%len(rr.Targets)]
	rr.index++
	return target
}
