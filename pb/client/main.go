package main

import (
	"context"
	pb "golang-code-learn/pb/hello_service"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	log.Print("===== 第一种 =====")
	cc, err := grpc.Dial(":18080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	shc := pb.NewSayHelloClient(cc)
	resp, err2 := shc.Hello(context.Background(), &pb.Request{Name: "hang"})
	if err2 != nil {
		log.Fatal(err2)
	}
	log.Print(resp.String())

	log.Print("===== 第二种 =====")
	cc2, err := grpc.Dial(":18080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	resp2 := &pb.Response{}
	ctx, cf := context.WithTimeout(context.Background(), time.Second*5)
	defer cf()
	err3 := cc2.Invoke(ctx, "/hello_service.SayHello/Hello", &pb.Request{Name: "lp"}, resp2)
	if err3 != nil {
		log.Fatal(err3)
	}
	log.Print(resp2.String())
}
