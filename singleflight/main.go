package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/singleflight"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// 同步单飞
	getData(context.Background(), "key", 5)

	fmt.Println("==============================")

	// 超时控制
	ctx, cf := context.WithTimeout(context.Background(), time.Second*1)
	defer cf()
	getDataTimeout(ctx, "key", 5)
}

func getData(ctx context.Context, key string, nums int) {
	wg := sync.WaitGroup{}
	for i := 0; i < nums; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := getDataBySingleFlight(ctx, "key")
			if err != nil {
				log.Printf("get data err:%v", err)
				return
			}
			log.Printf("get data success res:%v", res)
		}()
	}
	wg.Wait()
}

func getDataTimeout(ctx context.Context, key string, nums int) {
	wg := sync.WaitGroup{}
	for i := 0; i < nums; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := getDataBySingleFlightWithTimeOut(ctx, "key")
			if err != nil {
				log.Printf("get data err:%v", err)
				return
			}
			log.Printf("get data success res:%v", res)
		}()
	}
	wg.Wait()
}

var singleFlight singleflight.Group

func getDataBySingleFlight(ctx context.Context, key string) (string, error) {
	// localCache无数据，从redis查
	resData, err := getFromCache(ctx, key)
	if err != nil && err != redis.Nil {
		return resData, err
	} else if redis.Nil == err {
		// redis无数据 ，查db
		data, err, _ := singleFlight.Do(key, func() (interface{}, error) {
			//  建议休眠0.1s 捕获0.1s内的重复请求
			time.Sleep(100 * time.Millisecond)
			// 查db
			data, err := getFromDB(ctx, key)
			if err != nil {
				return data, err
			}
			// 回写 Cache
			err = set2Cache(ctx, data)
			if err != nil {
				return data, err
			}
			// 返回db数据,回写cache的error不上抛
			return data, nil
		})
		if err != nil {
			return resData, err
		}
		resData = data.(string)
	}
	return resData, nil
}

func getDataBySingleFlightWithTimeOut(ctx context.Context, key string) (string, error) {
	// localCache无数据，从redis查
	resData, err := getFromCache(ctx, key)
	if err != nil && err != redis.Nil {
		return resData, err
	} else if redis.Nil == err {
		// redis无数据 ，查db
		ch := singleFlight.DoChan(key, func() (interface{}, error) {
			//  建议休眠0.1s 捕获0.1s内的重复请求
			time.Sleep(100 * time.Millisecond)
			// 查db
			data, err := getFromDB(ctx, key)
			if err != nil {
				return data, err
			}
			// 回写 Cache
			err = set2Cache(ctx, data)
			if err != nil {
				return data, err
			}
			// 返回db数据,回写cache的error不上抛
			return data, nil
		})
		// 超时等待
		select {
		case <-ctx.Done():
			return "", errors.New("request timeout~")
		case data := <-ch:
			resData = data.Val.(string)
		}
	}
	return resData, nil
}

func getFromCache(ctx context.Context, key string) (string, error) {
	log.Print("cache missing~")
	return "", redis.Nil
}
func getFromDB(ctx context.Context, key string) (string, error) {
	time.Sleep(2 * time.Second)
	log.Print("get data from db success~")
	return "data", nil
}
func set2Cache(ctx context.Context, key string) error {
	log.Print("set cache success~")
	return nil
}
