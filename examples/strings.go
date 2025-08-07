package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RunStringExamples demonstrates Redis string operations
func RunStringExamples(rdb *redis.Client) {
	ctx := context.Background()
	fmt.Println("\n🔤 String Operations Examples")
	fmt.Println("============================")

	// Basic SET and GET
	fmt.Println("1. Basic SET and GET:")
	err := rdb.Set(ctx, "user:1", "Naim Islam", 0).Err()
	if err != nil {
		panic("Failed to set value: " + err.Error())
	}
	val, err := rdb.Get(ctx, "user:1").Result()
	if err != nil {
		panic("Failed to get value: " + err.Error())
	}
	fmt.Printf("user:1 = %s\n", val)

	// SET with expiration
	fmt.Println("\n2. SET with expiration (5 seconds):")
	err = rdb.Set(ctx, "temp:session", "12345", 5*time.Second).Err()
	if err != nil {
		panic("Failed to set value with expiration: " + err.Error())
	}

	ttl, err := rdb.TTL(ctx, "temp:session").Result()
	if err != nil {
		panic("Failed to get TTL: " + err.Error())
	}
	fmt.Printf(" temp:session will expire in %s\n", ttl)

	// INCR and DECR
	fmt.Println("\n3. Increment and Decrement:")
	err = rdb.Set(ctx, "counter", 10, 0).Err()
	if err != nil {
		panic("Failed to set initial counter value: " + err.Error())
	}

	newVal, err := rdb.Incr(ctx, "counter").Result()
	if err != nil {
		panic("Failed to increment counter: " + err.Error())
	}
	fmt.Printf("Counter after increment: %d\n", newVal)

	newVal, err = rdb.Decr(ctx, "counter").Result()
	if err != nil {
		panic("Failed to decrement counter: " + err.Error())
	}
	fmt.Printf("Counter after decrement: %d\n", newVal)

	// Append
	fmt.Println("\n4. APPEND operation:")
	err = rdb.Set(ctx, "message", "Hello", 0).Err()
	if err != nil {
		panic("Failed to set initial message: " + err.Error())
	}

	length, err := rdb.Append(ctx, "message", " World!").Result()
	if err != nil {
		panic("Failed to append to message: " + err.Error())
	}

	finalMsg, _ := rdb.Get(ctx, "message").Result()
	fmt.Printf(" Appended message: %s (length: %d)\n", finalMsg, length)

	// MSET and MGET (Multiple operations)
	fmt.Println("\n5. Multiple SET and GET:")
	err = rdb.MSet(ctx, "key1", "value1", "key2", "value2", "key3", "value3").Err()
	if err != nil {
		panic("Failed to set multiple values: " + err.Error())
	}

	values, err := rdb.MGet(ctx, "key1", "key2", "key3").Result()
	if err != nil {
		panic("Failed to get multiple values: " + err.Error())
	}

	for i, val := range values {
		if val == nil {
			fmt.Printf(" key%d = <nil>\n", i+1)
		} else {
			fmt.Printf(" key%d = %s\n", i+1, val)
		}
	}

	// Clean up
	rdb.Del(ctx, "user:1", "temp:session", "counter", "message", "key1", "key2", "key3")
}
