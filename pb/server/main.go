package main

import (
	"context"
	"log"
	"net"

	pb "golang-code-learn/pb/hello_service"

	"google.golang.org/grpc"
)

type HelloServer struct {
	pb.UnimplementedSayHelloServer
}

func (hs *HelloServer) Hello(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	dl, ok := ctx.Deadline()
	log.Printf("req[%+v] deadline[%+v] ok[%+v]", req.Name, dl, ok)
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
