package data

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func listPopPush(ctx context.Context, rc *redis.Client) {
	i, err := rc.RPush(context.TODO(), "l", 1, 2, 3, 4, 5).Result()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(i)
	log.Println("=============================")
	for {
		s, err := rc.LPop(context.Background(), "l").Result()
		if err != nil {
			if err == redis.Nil {
				break
			}
			log.Fatalln(err)
		}
		log.Println(s)
	}
}

func listPopPushWithLua(ctx context.Context, rc *redis.Client) {
	rc.FlushAll(ctx)

	script := `
		local key = KEYS[1]
		
		local member, err = redis.call('lpop', key)  
		if member then  
			redis.call('rpush', key, member)
		else  
			return err  
		end  
		
		return member
	`

	scr := redis.NewScript(script)

	i1, err := rc.LPush(ctx, "l1", "1", "2", "3", "4", "5", "hphang").Result()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(i1)

	i2, err := scr.Run(ctx, rc, []string{"l1"}).Result()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(i2)

	i3, err := rc.RPop(ctx, "l1").Result()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(i3)
}
