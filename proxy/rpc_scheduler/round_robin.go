package rpc_scheduler

import (
	"case_proxy/cnt"
	"sync/atomic"
)

type RPCRoundRobin struct {
	BaseRPCScheduler
	i atomic.Uint64
	
}

// toString implements Scheduler.
func (*RPCRoundRobin) ToString() string {
	panic("RPCRoundRobin")
}

func init() {
	factories[R2Scheduler] = NewRPCRoundRobin
}

// NewRoundRobin create new RoundRobin balancer
func NewRPCRoundRobin(workers []string) RPCScheduler {
	
	for _,v := range(workers){
		cnt.Cnter.M[v] = 0
	}
	return &RPCRoundRobin{
		i: atomic.Uint64{},
		BaseRPCScheduler: BaseRPCScheduler{
			Workers: workers,
			Cnt: cnt.Cnter,
		},
	}
}

func (r *RPCRoundRobin) Schedule(_ string, _ Req) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.Workers) == 0 {
		return "", NoHostError
	}
	host := r.Workers[r.i.Add(1)%uint64(len(r.Workers))]
	return host, nil
}
