package examples

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RunSetsExamples demonstrates Redis set operations
func RunSetsExamples(rdb *redis.Client) {
	fmt.Println("\n Set Operations")
	fmt.Println("=================")

	ctx := context.Background()

	// SADD - Add members to a set
	fmt.Println("1. Adding members with SADD:")

	// Create user interests
	interestSet := "user:123:interests"
	added, err := rdb.SAdd(ctx, interestSet, "programming", "music", "travel", "photography").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Added %d interests to user:123\n", added)

	// Try adding duplicate (won't be added)
	added, err = rdb.SAdd(ctx, interestSet, "programming", "reading").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Added %d new interests (duplicates ignored)\n", added)

	// SMEMBERS - Get all members
	fmt.Println("\n2. Getting all members with SMEMBERS:")
	interests, err := rdb.SMembers(ctx, interestSet).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   User interests: %v\n", interests)

	// SCARD - Get set size
	fmt.Println("\n3. Getting set size with SCARD:")
	size, err := rdb.SCard(ctx, interestSet).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Number of interests: %d\n", size)

	// SISMEMBER - Check if member exists
	fmt.Println("\n4. Checking membership with SISMEMBER:")
	isMember, err := rdb.SIsMember(ctx, interestSet, "programming").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Is 'programming' an interest? %t\n", isMember)

	isMember, err = rdb.SIsMember(ctx, interestSet, "cooking").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Is 'cooking' an interest? %t\n", isMember)

	// Create another user's interests for set operations
	fmt.Println("\n5. Creating another user's interests:")
	otherInterestSet := "user:456:interests"
	rdb.SAdd(ctx, otherInterestSet, "programming", "gaming", "travel", "cooking")

	otherInterests, err := rdb.SMembers(ctx, otherInterestSet).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   User 456 interests: %v\n", otherInterests)

	// SINTER - Set intersection (common interests)
	fmt.Println("\n6. Finding common interests with SINTER:")
	commonInterests, err := rdb.SInter(ctx, interestSet, otherInterestSet).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Common interests: %v\n", commonInterests)

	// SUNION - Set union (all unique interests)
	fmt.Println("\n7. Finding all unique interests with SUNION:")
	allInterests, err := rdb.SUnion(ctx, interestSet, otherInterestSet).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   All unique interests: %v\n", allInterests)

	// SDIFF - Set difference (interests only in first set)
	fmt.Println("\n8. Finding unique interests with SDIFF:")
	uniqueToUser123, err := rdb.SDiff(ctx, interestSet, otherInterestSet).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Interests unique to user 123: %v\n", uniqueToUser123)

	uniqueToUser456, err := rdb.SDiff(ctx, otherInterestSet, interestSet).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Interests unique to user 456: %v\n", uniqueToUser456)

	// SPOP - Remove and return random member
	fmt.Println("\n9. Random operations with SPOP and SRANDMEMBER:")
	randomInterest, err := rdb.SPop(ctx, interestSet).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Randomly removed interest: %s\n", randomInterest)

	// SRANDMEMBER - Get random member without removing
	randomMember, err := rdb.SRandMember(ctx, interestSet).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Random interest (not removed): %s\n", randomMember)

	// Get multiple random members
	randomMembers, err := rdb.SRandMemberN(ctx, interestSet, 2).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   2 random interests: %v\n", randomMembers)

	// SREM - Remove specific members
	fmt.Println("\n10. Removing specific members with SREM:")
	removed, err := rdb.SRem(ctx, interestSet, "music").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Removed %d member(s)\n", removed)

	remainingInterests, err := rdb.SMembers(ctx, interestSet).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Remaining interests: %v\n", remainingInterests)

	// Practical example: Tagging system
	fmt.Println("\n11. Practical example - Article tagging system:")

	// Article tags
	rdb.SAdd(ctx, "article:1:tags", "redis", "database", "nosql", "performance")
	rdb.SAdd(ctx, "article:2:tags", "golang", "programming", "performance", "backend")
	rdb.SAdd(ctx, "article:3:tags", "redis", "golang", "tutorial", "backend")

	// Find articles with common tags
	fmt.Println("   Articles tagged with 'redis':")
	// In a real system, you'd maintain reverse indexes
	// For demo, we'll check each article
	articles := []string{"article:1:tags", "article:2:tags", "article:3:tags"}
	for i, article := range articles {
		hasRedis, _ := rdb.SIsMember(ctx, article, "redis").Result()
		if hasRedis {
			fmt.Printf("     Article %d has 'redis' tag\n", i+1)
		}
	}

	// Find articles with multiple tags (intersection example)
	fmt.Println("   Articles tagged with both 'performance' AND 'backend':")
	for i, article := range articles {
		hasPerf, _ := rdb.SIsMember(ctx, article, "performance").Result()
		hasBackend, _ := rdb.SIsMember(ctx, article, "backend").Result()
		if hasPerf && hasBackend {
			fmt.Printf("     Article %d has both tags\n", i+1)
		}
	}

	// Practical example: Online users tracking
	fmt.Println("\n12. Practical example - Online users tracking:")
	onlineUsers := "online_users"

	// Users come online
	rdb.SAdd(ctx, onlineUsers, "user:123", "user:456", "user:789")

	// Check who's online
	online, err := rdb.SMembers(ctx, onlineUsers).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Online users: %v\n", online)

	// User goes offline
	rdb.SRem(ctx, onlineUsers, "user:456")

	// Check online count
	onlineCount, err := rdb.SCard(ctx, onlineUsers).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Online user count: %d\n", onlineCount)

	// Cleanup
	fmt.Println("\n13. Cleanup:")
	rdb.Del(ctx, interestSet, otherInterestSet, "article:1:tags", "article:2:tags", "article:3:tags", onlineUsers)
	fmt.Println("   Cleaned up set examples ✓")
}
