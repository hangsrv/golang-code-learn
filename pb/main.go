package main

import (
	pb "golang-code-learn/pb/hello_service"
	"log"

	"google.golang.org/protobuf/proto"
)

func main() {
	u := &pb.Request{
		Name: "hang",
	}
	log.Print(u.String())

	bs, err := proto.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(bs)

	u2 := &pb.Request{}
	if err := proto.Unmarshal(bs, u2); err != nil {
		log.Fatal(err)
	}
	log.Print(u2.String())
}
