package rpc_scheduler

import (
	"case_proxy/cnt"
	"case_proxy/conf"
	"math/rand"
	"time"

	"github.com/dgraph-io/ristretto"
)

type NaiveCache struct {
	BaseRPCScheduler
	Caches map[int] *ristretto.Cache
	rnd *rand.Rand
}
func init() {
	factories[NaiveCacheScheduler] = NewNaiveCache
}
func (*NaiveCache) ToString() string {
	panic("NaiveCache")
}

func NewNaiveCache(workers []string) RPCScheduler{
	cacheSize := conf.Conf.CacheSize
	for _,v := range(workers){
		cnt.Cnter.M[v] = 0
	}
	cache := &NaiveCache{
		Caches: make(map[int]*ristretto.Cache),
		BaseRPCScheduler: BaseRPCScheduler{
			Workers: workers,
			Cnt: cnt.Cnter,
		},
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	config := &ristretto.Config{
		NumCounters: 1e7,     // Num keys to track frequency of (10M).
		MaxCost:     cacheSize, // Maximum cost of cache (1GB).
		BufferItems: 64,      // Number of keys per Get buffer.
	}
	for i :=0; i< len(workers);i++{
		rcache, err := ristretto.NewCache(config)
		if err != nil {
			return nil
		}
		cache.Caches[i] = rcache
	}
	return cache

}


func (n *NaiveCache) Schedule(_ string, r Req) (string, error) {
	var hostIndex = -1
	var maxCnt = 0

	keys := r.Keys

	for i :=0; i< len(n.Workers);i++{
			cache := n.Caches[i]
			cnt := 0
			for _, a := range keys{
				for jj:=0; jj<1;jj++{
						for _ ,aa := range a {
					if _,ok := cache.Get(aa); ok{
						cnt  = cnt + 1
					}
				}
				}

			}
			if cnt > maxCnt{
				maxCnt  = cnt
				hostIndex = i
			}
	}
	if hostIndex == -1{
		hostIndex = n.rnd.Intn(len(n.Workers))
	}
	cache := n.Caches[hostIndex]
	for _, a := range keys{
		for _ ,aa := range a {
			cache.Set(aa,aa,1)
		}
	}
	return n.Workers[hostIndex],nil

}
