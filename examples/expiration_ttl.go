package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RunExpirationTTLExamples demonstrates Redis expiration and TTL operations
func RunExpirationTTLExamples(rdb *redis.Client) {
	fmt.Println("\n⏳ Expiration & TTL Operations")
	fmt.Println("==============================")

	ctx := context.Background()
	key := "temp:data"
	value := "This is a temporary value"

	// SET with expiration
	fmt.Println("1. Setting key with expiration (10 seconds):")
	err := rdb.Set(ctx, key, value, 10*time.Second).Err()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Key '%s' set with value '%s' and TTL 10s\n", key, value)

	// Check TTL
	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   TTL for key '%s': %v\n", key, ttl)

	// Get value before expiration
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Value before expiration: %s\n", val)

	// Wait for 11 seconds to expire
	fmt.Println("   Waiting for key to expire...")
	time.Sleep(11 * time.Second)

	// Try to get value after expiration
	val, err = rdb.Get(ctx, key).Result()
	if err != nil {
		fmt.Printf("   Value after expiration: (expired or missing)\n")
	} else {
		fmt.Printf("   Value after expiration: %s\n", val)
	}

	// EXPIRE/PEXPIRE - Set or update expiration
	fmt.Println("\n2. Using EXPIRE to set/update expiration:")
	rdb.Set(ctx, key, value, 0)
	err = rdb.Expire(ctx, key, 5*time.Second).Err()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("   Expiration updated to 5 seconds")

	ttl, _ = rdb.TTL(ctx, key).Result()
	fmt.Printf("   New TTL: %v\n", ttl)

	// PERSIST - Remove expiration from a key
	fmt.Println("\n3. Using PERSIST to make key permanent:")
	err = rdb.Persist(ctx, key).Err()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	ttl, _ = rdb.TTL(ctx, key).Result()
	fmt.Printf("   TTL after PERSIST: %v (should be -1 for permanent)\n", ttl)

	// Practical example: Session expiration
	fmt.Println("\n4. Practical example - Session expiration:")
	sessionKey := "session:xyz"
	rdb.Set(ctx, sessionKey, "user_data", 3*time.Second)
	fmt.Println("   Session created with 3s TTL")
	time.Sleep(4 * time.Second)
	_, err = rdb.Get(ctx, sessionKey).Result()
	if err != nil {
		fmt.Println("   Session expired and key deleted!")
	} else {
		fmt.Println("   Session still exists (unexpected)")
	}

	// Cleanup
	rdb.Del(ctx, key, sessionKey)
	fmt.Println("\n5. Cleanup: Cleaned up expiration examples ✓")
}
