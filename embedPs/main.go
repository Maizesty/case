package main

import (
	"fmt"
	"net"
	"context"
	"embed_server/conf"
	pb "embed_server/embed"  
	"google.golang.org/grpc"
)



type server struct{ 
	pb.EmbServerServer
}

func (s *server) Lookup(ctx context.Context, in *pb.EmbReqData) (*pb.EmbResp, error) {
	keys := in.Keys
	out := make([]*pb.EmbVector,len(keys))
	for i := 0; i < len(keys);i++{
		out[i] = &pb.EmbVector{Element: Local_emb[keys[i]%conf.EmbMax32]}
	}
	return &pb.EmbResp{EmbVectors: out}, nil
}

func main(){
	port := conf.Port
	fmt.Println("binding to port: ", port)
	
	// 设置监听地址和端口
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v",port))
	if err != nil {
		fmt.Println("failed to listen: ", err)
	}

		// 实例化1个服务器程序
	s := grpc.NewServer()

		// 调用服务注册函数
	pb.RegisterEmbServerServer(s, &server{})
	fmt.Println("server listening at", lis.Addr())

		// 在监听端口上运行服务器程序
	if err := s.Serve(lis); err != nil {
		fmt.Println("failed to serve: ", err)
	}
}