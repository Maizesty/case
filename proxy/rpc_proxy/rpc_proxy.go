// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc_proxy

import (
	"case_proxy/infer"
	"case_proxy/rpc_scheduler"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PRCProxy struct {
	hostMap      map[string]*infer.MyServiceClient
	lb           rpc_scheduler.RPCScheduler
	Ctx          context.Context
	sync.RWMutex // protect alive
	alive        map[string]bool
}
type Resp struct {
	Output  float32 `json:"output"`
	Perfect bool    `json:"perfect"`
	// Rank	int8 `json:"rank"`
	// KeysEvicted uint64 `json:"keys_evicted"`
	// Ratio				float64	`json:"ratio"`
}

// func test(){
// 	name := "World"
// 	addr := "127.0.0.1:50051"
// 	// 连接gRPC服务器
//         conn, err := grpc.Dial( addr , grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		fmt.Println("connect error to: ",addr)
// 	}
// 	defer conn.Close()

//         // 实例化一个client对象，传入参数conn
// 	c := infer.NewMyServiceClient(conn)

// 	// 初始化上下文，设置请求超时时间为1秒
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	//延迟关闭请求会话
// 	defer cancel()

// 	// 调用SayHello方法，以请求服务，然后得到响应消息
// 	r, err := c.SayHello(ctx, &.HelloRequest{Name: name})
// 	if err != nil {
// 		fmt.Println("can not greet to: ",addr)
// 	} else {
// 		fmt.Println("response from server: ",r.GetMessage())
// 	}
// }

func (h *PRCProxy) Setlb(targetHosts []string, algorithm string) (
	*PRCProxy, error) {
	hosts := make([]string, 0)
	hostMap := make(map[string]*infer.MyServiceClient)
	// alive := make(map[string]bool)
	for _, targetHost := range targetHosts {
		// url, err := url.Parse(targetHost)
		// if err != nil {
		// 	return nil, err
		// }
		conn, err := grpc.Dial(targetHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Println("connect error to: ", targetHost)
		}
		proxy := infer.NewMyServiceClient(conn)
		hostMap[targetHost] = &proxy
		hosts = append(hosts, targetHost)
	}

	lb, err := rpc_scheduler.Build(algorithm, hosts)
	if err != nil {
		return nil, err
	}
	h.lb = lb
	return h, nil
}

// NewHTTPProxy create  new reverse proxy with url and balancer algorithm
func NewPRCProxy(targetHosts []string, algorithm string) (
	*PRCProxy, error) {

	hosts := make([]string, 0)
	hostMap := make(map[string]*infer.MyServiceClient)
	alive := make(map[string]bool)
	for _, targetHost := range targetHosts {
		// url, err := url.Parse(targetHost)
		// if err != nil {
		// 	return nil, err
		// }
		conn, err := grpc.Dial(targetHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Println("connect error to: ", targetHost)
		}
		proxy := infer.NewMyServiceClient(conn)
		hostMap[targetHost] = &proxy
		hosts = append(hosts, targetHost)
	}

	lb, err := rpc_scheduler.Build(algorithm, hosts)
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return &PRCProxy{
		hostMap: hostMap,
		lb:      lb,
		Ctx:     ctx,
		alive:   alive,
	}, nil
}

// ServeHTTP implements a proxy to the http server
func (h *PRCProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("proxy causes panic :%s", err)
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(err.(error).Error()))
		}
	}()
	var rr rpc_scheduler.Req
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &rr)
	if err != nil {
		return
	}
	host, err := h.lb.Schedule("", rr)
	// fmt.Println(h.lb.ToString())
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(fmt.Sprintf("balance error: %s", err.Error())))
		return
	}
	h.lb.Inc(host)
	defer h.lb.Done(host)
	c := h.hostMap[host]
	cc := *c
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := cc.Infer(ctx, &infer.ReqData{Keys: rr.Keys[0]})
	if err != nil {
		fmt.Println("can not greet to: ", host,err)
	} 
	resp := Resp{Output: result.Output, Perfect: result.Perfect}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		// 处理错误
		http.Error(w, "Error serializing struct to JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	// h.hostMap[host].ServeHTTP(w, r)
}
