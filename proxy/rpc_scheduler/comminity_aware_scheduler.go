package rpc_scheduler

import (
	"case_proxy/cnt"
	"case_proxy/load"
	"math/rand"
	"time"
)

func init() {
	factories[CommiunityAwareScheduler] = NewCommiunityAware
}

type CommiunityAware struct {
	BaseRPCScheduler
	rnd *rand.Rand
	rest int
}

// toString implements Scheduler.
func (*CommiunityAware) ToString() string {
	return "CommiunityAware"
}

// Schedule implements Scheduler.
func (c *CommiunityAware) Schedule(_ string, r Req) (string, error) {

	return c.vote(r.Keys), nil
}

func (c *CommiunityAware) vote(keys [][]int32) string {
	counter := make([]int, len(c.Workers))
	maxVal := 0
	maxIndex := 0
	for _, a := range keys {
		for jj:=0; jj<1;jj++{
			for _, aa := range a {
			if _, ok := load.HotCache[aa]; !ok {
				index := load.Partition[aa] % len(c.Workers)
				counter[index]++
				if counter[index] > maxVal {
					maxVal = counter[index]
					maxIndex = index
				}
			}
		}
		}

	}
	maxIndex = maxIndex 
	n1 := c.Workers[c.generateRandomNumber(maxIndex)]
	return n1
}

func NewCommiunityAware(workers []string) RPCScheduler {
	for _,v := range(workers){
		cnt.Cnter.M[v] = 0
	}
	return &CommiunityAware{
		BaseRPCScheduler: BaseRPCScheduler{
			Workers: workers,
			Cnt: cnt.Cnter,
		},
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
		rest : len(workers) / 4 ,
	}
}
func (c *CommiunityAware)generateRandomNumber(input int) int {
	// 初始化随机数生成器
	

	// 根据输入数字生成随机数
	switch input {
	case 0:
		return c.rnd.Intn(6) // 生成 0 到 5 之间的随机数
	case 1:
		return c.rnd.Intn(2) + 6 // 生成 6 到 7 之间的随机数
	case 2:
		return  8 // 生成 8 到 9 之间的随机数
	case 3:
		return c.rnd.Intn(3)  + 9// 生成 10 到 11 之间的随机数
	default:
		return -1 // 输入错误时返回一个负数
	}
}
