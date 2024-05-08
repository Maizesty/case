package rpc_scheduler

import (
	"sync"
	"case_proxy/cnt"
)

type BaseRPCScheduler struct {
	sync.RWMutex
	Workers []string
	Cnt *cnt.Smap
}

func (b *BaseRPCScheduler) Add(worker string) {
	b.Lock()
	defer b.Unlock()
	for _, h := range b.Workers {
		if h == worker {
			return
		}
	}
	b.Workers = append(b.Workers, worker)
}

// Remove new host from the balancer
func (b *BaseRPCScheduler) Remove(worker string) {
	b.Lock()
	defer b.Unlock()
	for i, h := range b.Workers {
		if h == worker {
			b.Workers = append(b.Workers[:i], b.Workers[i+1:]...)
			return
		}
	}
}

// Balance selects a suitable host according
// func (b *BaseRPCScheduler) Balance(key string, r *http.Request) (string, error) {
// 	return "", nil
// }

// Inc .
func (b *BaseRPCScheduler) Inc(_ string) {}

// Done .
func (b *BaseRPCScheduler) Done(w string) {
	b.Cnt.Lock()
	b.Cnt.M[w]++
	b.Cnt.Unlock()
}
