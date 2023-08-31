package common

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

func NewClient() *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		Password: "123456",
	})
	sc := cli.Ping(context.Background())
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}
	return cli
}
