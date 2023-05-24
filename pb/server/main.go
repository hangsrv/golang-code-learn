package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "golang-code-learn/pb/hello_service"

	"google.golang.org/grpc"
)

type HelloServer struct {
	pb.UnimplementedSayHelloServer
}

func (hs *HelloServer) Hello(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Printf("req[%s]", req.Name)
	time.Sleep(2 * time.Second)
	return &pb.Response{Message: "hello " + req.Name}, nil
}

func main() {
	server := grpc.NewServer()

	pb.RegisterSayHelloServer(server, &HelloServer{})

	l, err := net.Listen("tcp", ":18080")
	if err != nil {
		panic(err)
	}

	server.Serve(l)
}
