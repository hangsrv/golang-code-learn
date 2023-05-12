package main

import (
	"golang-code-learn/pb/user"
	"log"

	"google.golang.org/protobuf/proto"
)

func main() {
	u := &user.User{
		Name:  "hang",
		Age:   23,
		Hobby: []string{"1", "2", "3"},
	}
	log.Print(u.String())

	bs, err := proto.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(bs)

	u2 := &user.User{}
	if err := proto.Unmarshal(bs, u2); err != nil {
		log.Fatal(err)
	}
	log.Print(u2.String())
}
