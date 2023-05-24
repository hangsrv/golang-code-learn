package main

import (
	"context"
	pb "golang-code-learn/pb/hello_service"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	// 连接设置 timeout = 1s 的超时时间
	cc, err := grpc.Dial(":18080", grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		log.Fatal(err)
	}
	// 请求设置 timeout = 3s 的超时时间
	ctx, cf := context.WithTimeout(context.Background(), time.Second*3)
	defer cf()
	resp := &pb.Response{}
	if err := cc.Invoke(ctx, "/hello_service.SayHello/Hello", &pb.Request{Name: "lp"}, resp); err != nil {
		log.Fatal(err)
	}
	log.Print(resp)
}
