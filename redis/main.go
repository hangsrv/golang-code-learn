package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// docker run -id --name redis -p 6379:6379 redis --requirepass "123456"
func main() {
	cli := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		Password: "123456",
	})
	sc := cli.Ping(context.Background())
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}
	log.Print(sc.Val())
}
