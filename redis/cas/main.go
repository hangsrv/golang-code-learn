package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456",
	})

	wg := &sync.WaitGroup{}
	key := "key"
	ctx := context.Background()
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// cget方法
			cgetValue, cgetVersion, err := cget(client, ctx, key)
			if err != nil {
				fmt.Println("Error in cget:", err)
				return
			}
			fmt.Printf("cget result: key=%s value=%s, version=%d\n", key, cgetValue, cgetVersion)

			time.Sleep(2 * time.Second)

			// cset方法
			newValue := "value"
			err = cset(client, ctx, key, cgetVersion, newValue)
			if err != nil {
				fmt.Println("Error in cset:", err)
				return
			}
			fmt.Println("cset succeeded")

			// 再次执行cget
			cgetValue, cgetVersion, err = cget(client, ctx, key)
			if err != nil {
				fmt.Println("Error in cget:", err)
				return
			}
			fmt.Printf("cget result after cset: key=%s value=%s, version=%d\n", key, cgetValue, cgetVersion)
		}()
	}
	wg.Wait()

	cgetValue, cgetVersion, err := cget(client, ctx, key)
	if err != nil {
		fmt.Println("Error in cget:", err)
		return
	}
	fmt.Printf("cget result: key=%s value=%s, version=%d\n", key, cgetValue, cgetVersion)
}

// cget 获取值和版本号
func cget(client *redis.Client, ctx context.Context, key string) (string, int, error) {
	script := `
		local key = KEYS[1]
		local versionKey = key .. "_version"
		local value = redis.call("GET", key)
		local version = tonumber(redis.call("GET", versionKey))
		if value == false or version == nil then
			value = ""
			version = 0
		end
		return {value, version}
	`
	key = fmt.Sprintf("{%s}", key)
	result, err := client.Eval(ctx, script, []string{key}).Result()
	if err != nil {
		return "", 0, err
	}

	data := result.([]interface{})
	value := data[0].(string)
	version := int(data[1].(int64))
	return value, version, nil
}

// cset CAS设置值
func cset(client *redis.Client, ctx context.Context, key string, expectedVersion int, newValue string) error {
	script := `
		local key = KEYS[1]
		local versionKey = key .. "_version"
		local expectedVersion = tonumber(ARGV[1])
		local newValue = ARGV[2]

		local currentVersion = tonumber(redis.call("GET", versionKey))

		if currentVersion == nil then
			currentVersion = 0
		end

		if currentVersion == expectedVersion then
		    redis.call("SET", key, newValue)
		    redis.call("INCR", versionKey)
		    return "OK"
		else
			return redis.error_reply("CAS operation failed")
		end
	`
	key = fmt.Sprintf("{%s}", key)
	_, err := client.Eval(ctx, script, []string{key}, expectedVersion, newValue).Result()
	if err != nil {
		return err
	}

	return nil
}
