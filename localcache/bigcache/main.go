package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/allegro/bigcache"
)

func main() {
	// 配置参数
	cacheSize := 64     // 64MB
	concurrency := 20  // 并发度
	ratio := 0.8        // 读写占比
	totalOps := 10000000 // 总操作数

	// 初始化缓存
	cache, err := bigcache.NewBigCache(bigcache.Config{
		Shards:             1024, // 分片数量
		LifeWindow:         0,    // 缓存过期时间（不会自动过期）
		CleanWindow:        0,
		MaxEntriesInWindow: 0,
		MaxEntrySize:       500,
		Verbose:            false,
		HardMaxCacheSize:   cacheSize,
	})
	if err != nil {
		panic(err)
	}

	// 统计变量
	var readCount int64
	var writeCount int64
	var readTime int64
	var writeTime int64
	var ops int64
	var wg sync.WaitGroup

	// 并发读写
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()

			rand.Seed(time.Now().UnixNano())
			for {
				key := fmt.Sprintf("key-%d", rand.Intn(10000))
				value := []byte("value")
				if ops >= int64(totalOps) {
					return
				}
				if rand.Float64() < ratio { // 读
					startTime := time.Now()
					cache.Get(key)
					atomic.AddInt64(&readTime, int64(time.Since(startTime)))
					atomic.AddInt64(&readCount, 1)
					atomic.AddInt64(&ops, 1)
				} else { // 写
					startTime := time.Now()
					cache.Set(key, value)
					atomic.AddInt64(&writeTime, int64(time.Since(startTime)))
					atomic.AddInt64(&writeCount, 1)
					atomic.AddInt64(&ops, 1)
				}
			}
		}()
	}

	// 等待所有协程结束
	wg.Wait()
	// 输出统计结果
	fmt.Printf("Cache size: %d MB\n", cacheSize)
	fmt.Printf("Concurrency: %d\n", concurrency)
	fmt.Printf("Read/write ratio: %.2f/%.2f\n", ratio, 1-ratio)
	fmt.Printf("Total operations: %d\n", readCount+writeCount)
	fmt.Printf("Read operations: %d\n", readCount)
	fmt.Printf("Write operations: %d\n", writeCount)
	fmt.Printf("Average read time: %s\n", time.Duration(readTime/readCount))
	fmt.Printf("Average write time: %s\n", time.Duration(writeTime/writeCount))
}
