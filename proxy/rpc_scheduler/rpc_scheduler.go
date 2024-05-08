package rpc_scheduler

import (
	"errors"
)

var (
	NoHostError                = errors.New("no host")
	AlgorithmNotSupportedError = errors.New("algorithm not supported")
)


type RPCScheduler interface {
	// Add 增加主机
	Add(string)
	// Remove 删除主机
	Remove(string)
	// Schedule 进行调度
	Schedule(string, Req) (string, error)
	// Inc 增加一个负载
	Inc(string)
	// Done 负载完成
	Done(string)
	ToString() string
}

type RPCFactory func([]string) RPCScheduler

var factories = make(map[string]RPCFactory)

func Build(algorithm string, workers []string) (RPCScheduler, error) {
	factory, ok := factories[algorithm]
	if !ok {
		return nil, AlgorithmNotSupportedError
	}
	return factory(workers), nil
}
