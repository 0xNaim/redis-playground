package examples

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RunHashesExamples demonstrates Redis hash operations
func RunHashesExamples(rdb *redis.Client) {
	fmt.Println("\n  Hash Operations")
	fmt.Println("===================")

	ctx := context.Background()

	// HSET - Set hash field values
	fmt.Println("1. Creating user profile with HSET:")
	err := rdb.HSet(ctx, "user:123", map[string]interface{}{
		"name":     "John Doe",
		"email":    "john@example.com",
		"age":      "30",
		"location": "San Francisco",
		"role":     "Developer",
	}).Err()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("   User profile created ✓")

	// HGET - Get specific field value
	fmt.Println("\n2. Getting specific fields with HGET:")
	name, err := rdb.HGet(ctx, "user:123", "name").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Name: %s\n", name)

	email, err := rdb.HGet(ctx, "user:123", "email").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Email: %s\n", email)

	// HGETALL - Get all fields and values
	fmt.Println("\n3. Getting all fields with HGETALL:")
	userProfile, err := rdb.HGetAll(ctx, "user:123").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("   Complete profile:")
	for field, value := range userProfile {
		fmt.Printf("     %s: %s\n", field, value)
	}

	// HMGET - Get multiple fields at once
	fmt.Println("\n4. Getting multiple fields with HMGET:")
	fields, err := rdb.HMGet(ctx, "user:123", "name", "role", "location").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fieldNames := []string{"name", "role", "location"}
	fmt.Println("   Selected fields:")
	for i, field := range fieldNames {
		fmt.Printf("     %s: %v\n", field, fields[i])
	}

	// HEXISTS - Check if field exists
	fmt.Println("\n5. Checking field existence with HEXISTS:")
	exists, err := rdb.HExists(ctx, "user:123", "age").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Field 'age' exists: %t\n", exists)

	exists, err = rdb.HExists(ctx, "user:123", "salary").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Field 'salary' exists: %t\n", exists)

	// HKEYS - Get all field names
	fmt.Println("\n6. Getting all field names with HKEYS:")
	keys, err := rdb.HKeys(ctx, "user:123").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Available fields: %v\n", keys)

	// HVALS - Get all values
	fmt.Println("\n7. Getting all values with HVALS:")
	values, err := rdb.HVals(ctx, "user:123").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   All values: %v\n", values)

	// HLEN - Get number of fields
	fmt.Println("\n8. Getting field count with HLEN:")
	fieldCount, err := rdb.HLen(ctx, "user:123").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Number of fields: %d\n", fieldCount)

	// HINCRBY - Increment numeric field
	fmt.Println("\n9. Incrementing numeric fields with HINCRBY:")
	newAge, err := rdb.HIncrBy(ctx, "user:123", "age", 1).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Age after increment: %d\n", newAge)

	// HDEL - Delete specific fields
	fmt.Println("\n10. Deleting fields with HDEL:")
	deleted, err := rdb.HDel(ctx, "user:123", "location").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Deleted %d field(s)\n", deleted)

	// Verify deletion
	remainingFields, err := rdb.HKeys(ctx, "user:123").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Remaining fields: %v\n", remainingFields)

	// Practical example: Session management
	fmt.Println("\n11. Practical example - Session management:")
	sessionID := "session:abc123"
	err = rdb.HSet(ctx, sessionID, map[string]interface{}{
		"user_id":    "123",
		"username":   "johndoe",
		"login_time": "2024-01-01T10:00:00Z",
		"ip_address": "192.168.1.1",
		"user_agent": "Mozilla/5.0...",
	}).Err()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Session %s created ✓\n", sessionID)

	// Get session info
	sessionData, err := rdb.HGetAll(ctx, sessionID).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("   Session data:")
	for k, v := range sessionData {
		fmt.Printf("     %s: %s\n", k, v)
	}

	// Cleanup
	fmt.Println("\n12. Cleanup:")
	rdb.Del(ctx, "user:123", sessionID)
	fmt.Println("   Cleaned up hash examples ✓")
}
