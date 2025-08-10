package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RunCachingExamples demonstrates Redis caching patterns
func RunCachingExamples(rdb *redis.Client) {
	fmt.Println("\n  Caching Examples")
	fmt.Println("=====================")

	ctx := context.Background()
	cacheKey := "cache:user:42"
	dbValue := "Naim"

	// 1. Cache-aside pattern
	fmt.Println("1. Cache-aside pattern:")
	val, err := rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		fmt.Println("   Cache miss! Fetching from DB...")
		val = dbValue
		rdb.Set(ctx, cacheKey, val, 10*time.Second)
		fmt.Println("   Value cached in Redis")
	} else if err != nil {
		fmt.Printf("   Redis error: %v\n", err)
		return
	} else {
		fmt.Println("   Cache hit!")
	}
	fmt.Printf("   Value: %s\n", val)

	// 2. Expiring cache
	fmt.Println("\n2. Expiring cache:")
	rdb.Set(ctx, "cache:expiring", "temporary", 3*time.Second)
	val, _ = rdb.Get(ctx, "cache:expiring").Result()
	fmt.Printf("   Value before expire: %s\n", val)
	time.Sleep(4 * time.Second)
	val, err = rdb.Get(ctx, "cache:expiring").Result()
	if err == redis.Nil {
		fmt.Println("   Value after expire: (cache expired)")
	}

	// 3. Manual cache invalidation
	fmt.Println("\n3. Manual cache invalidation:")
	rdb.Set(ctx, "cache:invalidate", "stale", 0)
	rdb.Del(ctx, "cache:invalidate")
	val, err = rdb.Get(ctx, "cache:invalidate").Result()
	if err == redis.Nil {
		fmt.Println("   Value after invalidation: (no cache)")
	}

	// Practical example: Caching expensive computation
	fmt.Println("\n4. Practical example - Caching computed result:")
	expensiveKey := "cache:expensive"
	val, err = rdb.Get(ctx, expensiveKey).Result()
	if err == redis.Nil {
		fmt.Println("   Cache miss! Running expensive operation...")
		val = "Expensive Result"
		rdb.Set(ctx, expensiveKey, val, 5*time.Second)
		fmt.Println("   Computed result cached")
	} else {
		fmt.Println("   Cache hit!")
	}
	fmt.Printf("   Expensive operation result: %s\n", val)

	// Cleanup
	rdb.Del(ctx, cacheKey, "cache:expiring", "cache:invalidate", expensiveKey)
	fmt.Println("\n5. Cleanup: Cleaned up caching examples ✓")
}
