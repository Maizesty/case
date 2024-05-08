package main

import(
	embed "embed_client/embed"
	is "embed_client/infer"
	"embed_client/conf"
	"fmt"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
	// "sync"
	"github.com/dgraph-io/ristretto"
	"sync/atomic"
	"net"
	tg "github.com/galeone/tfgo"
	// tf "github.com/galeone/tensorflow/tensorflow/go"
)

var Conn *grpc.ClientConn
var Ctx context.Context
var cache *ristretto.Cache
var CntHit atomic.Uint64 = atomic.Uint64{}
var CntTotal atomic.Uint64 = atomic.Uint64{}
var Model                  	*tg.Model
var EmbClient embed.EmbServerClient
func init(){
	loadModel()
}

func loadModel() {
	Model = tg.LoadModel(conf.ModelPath, []string{"serve"}, nil)
}

func main(){

	port := conf.Port
	fmt.Println("binding to port: ", port)
	
	// 设置监听地址和端口
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v",port))
	if err != nil {
		fmt.Println("failed to listen: ", err)
	}

	s := grpc.NewServer()
	is.RegisterMyServiceServer(s, &server{})
	fmt.Println("server listening at", lis.Addr())

	// name := "World"
	addr := conf.EmbServerIps[0]
	// 连接gRPC服务器
  Conn, err := grpc.Dial( addr , grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("connect error to: ",addr)
	}
	defer Conn.Close()

        // 实例化一个client对象，传入参数conn
	EmbClient = embed.NewEmbServerClient(Conn)

	// 初始化上下文，设置请求超时时间为1秒
	Ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	//延迟关闭请求会话
	// defer cancel()
	var cerr error
	cache, cerr = ristretto.NewCache(&ristretto.Config{
		NumCounters:        conf.CacheSize * 1000,
		MaxCost:            conf.CacheSize,
		BufferItems:        64,
		IgnoreInternalCost: true,
		Metrics:            true,
	})
	if cerr != nil {
		return
	}
	if err := s.Serve(lis); err != nil {
		fmt.Println("failed to serve: ", err)
	}


}