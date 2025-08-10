package examples

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RunListExamples demonstrates Redis list operations
func RunListExamples(rdb *redis.Client) {
	fmt.Println("\n List Operations")
	fmt.Println("==================")

	ctx := context.Background()

	// LPUSH/RPUSH - Add elements to the left/right of the list
	fmt.Println("1. Adding elements with LPUSH and RPUSH:")

	// Create a task queue
	listKey := "task_queue"

	// Add tasks to the right (end) of the queue
	length, err := rdb.RPush(ctx, listKey, "task1", "task2", "task3").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   RPUSH task_queue task1 task2 task3: length = %d\n", length)

	// Add urgent task to the left (beginning) of the queue
	length, err = rdb.LPush(ctx, listKey, "urgent_task").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   LPUSH task_queue urgent_task: length = %d\n", length)

	// LRANGE - Get elements from the list
	fmt.Println("\n2. Viewing list contents with LRANGE:")
	tasks, err := rdb.LRange(ctx, listKey, 0, -1).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Current queue: %v\n", tasks)

	// Get first 2 elements
	firstTwo, err := rdb.LRange(ctx, listKey, 0, 1).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   First 2 tasks: %v\n", firstTwo)

	// LPOP/RPOP - Remove and return elements
	fmt.Println("\n3. Processing tasks with LPOP and RPOP:")

	// Process from the left (FIFO - First In, First Out)
	task, err := rdb.LPop(ctx, listKey).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   LPOP (processed): %s\n", task)

	// Check remaining tasks
	remaining, err := rdb.LRange(ctx, listKey, 0, -1).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Remaining tasks: %v\n", remaining)

	// LLEN - Get list length
	fmt.Println("\n4. Checking queue size with LLEN:")
	queueSize, err := rdb.LLen(ctx, listKey).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Queue size: %d\n", queueSize)

	// LINDEX - Get element at specific index
	fmt.Println("\n5. Getting specific elements with LINDEX:")
	firstTask, err := rdb.LIndex(ctx, listKey, 0).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   First task (index 0): %s\n", firstTask)

	lastTask, err := rdb.LIndex(ctx, listKey, -1).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Last task (index -1): %s\n", lastTask)

	// LSET - Set element at specific index
	fmt.Println("\n6. Updating elements with LSET:")
	err = rdb.LSet(ctx, listKey, 0, "updated_task1").Err()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("   LSET task_queue 0 'updated_task1' ✓")

	updated, err := rdb.LRange(ctx, listKey, 0, -1).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Updated queue: %v\n", updated)

	// LREM - Remove elements
	fmt.Println("\n7. Removing specific elements with LREM:")

	// Add some duplicate elements first
	rdb.RPush(ctx, listKey, "duplicate", "duplicate", "unique")

	// Remove 2 occurrences of "duplicate" from the list
	removed, err := rdb.LRem(ctx, listKey, 2, "duplicate").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   LREM task_queue 2 'duplicate': removed %d elements\n", removed)

	afterRemoval, err := rdb.LRange(ctx, listKey, 0, -1).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   After removal: %v\n", afterRemoval)

	// Practical example: Activity feed
	fmt.Println("\n8. Practical example - Activity feed:")
	feedKey := "user:123:activity_feed"

	// Add activities (newest first)
	activities := []string{
		"User logged in",
		"User updated profile",
		"User posted a comment",
		"User liked a post",
		"User shared an article",
	}

	for _, activity := range activities {
		rdb.LPush(ctx, feedKey, activity)
	}

	// Get recent activities (limit to 3)
	recentActivities, err := rdb.LRange(ctx, feedKey, 0, 2).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("   Recent activities:")
	for i, activity := range recentActivities {
		fmt.Printf("     %d. %s\n", i+1, activity)
	}

	// Maintain feed size (keep only last 10 activities)
	err = rdb.LTrim(ctx, feedKey, 0, 9).Err() // Keep elements from index 0 to 9
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("   Activity feed trimmed to last 10 items ✓")

	// Stack example (LIFO - Last In, First Out)
	fmt.Println("\n9. Stack example (LIFO):")
	stackKey := "operation_stack"

	// Push operations
	rdb.LPush(ctx, stackKey, "operation1", "operation2", "operation3")

	// Pop operations (LIFO order)
	fmt.Println("   Popping from stack:")
	for i := 0; i < 3; i++ {
		op, err := rdb.LPop(ctx, stackKey).Result()
		if err == redis.Nil {
			break
		} else if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		}
		fmt.Printf("     Popped: %s\n", op)
	}

	// Cleanup
	fmt.Println("\n10. Cleanup:")
	rdb.Del(ctx, listKey, feedKey, stackKey)
	fmt.Println("   Cleaned up list examples ✓")
}
