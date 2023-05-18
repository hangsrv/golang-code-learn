package main

import (
	"context"
	"fmt"
	pb "golang-code-learn/pb/hello_service"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("===== 第一种 =====")
	cc, err := grpc.Dial(":18080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	shc := pb.NewSayHelloClient(cc)
	resp, err2 := shc.Hello(context.Background(), &pb.Request{Name: "hang"})
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(resp.String())

	fmt.Println("===== 第二种 =====")
	cc2, err := grpc.Dial(":18080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	resp2 := &pb.Response{}
	err3 := cc2.Invoke(context.Background(), "/hello_service.SayHello/Hello", &pb.Request{Name: "lp"}, resp2)
	if err3 != nil {
		panic(err3)
	}
	fmt.Println(resp2.String())
}
