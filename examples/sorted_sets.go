package examples

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RunSortedSetsExamples demonstrates Redis sorted set operations
func RunSortedSetsExamples(rdb *redis.Client) {
	fmt.Println("\n Sorted Set Operations")
	fmt.Println("========================")

	ctx := context.Background()

	// ZADD - Add members with scores
	fmt.Println("1. Adding members with scores using ZADD:")

	leaderboard := "game:leaderboard"

	// Add players with their scores
	players := []redis.Z{
		{Score: 1500, Member: "alice"},
		{Score: 2300, Member: "bob"},
		{Score: 1800, Member: "charlie"},
		{Score: 2100, Member: "diana"},
		{Score: 1200, Member: "eve"},
	}

	added, err := rdb.ZAdd(ctx, leaderboard, players...).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Added %d players to leaderboard\n", added)

	// ZRANGE - Get members by rank (ascending order)
	fmt.Println("\n2. Getting members by rank with ZRANGE:")

	// Get all players (lowest to highest score)
	allPlayers, err := rdb.ZRangeWithScores(ctx, leaderboard, 0, -1).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("   All players (ascending):")
	for i, player := range allPlayers {
		fmt.Printf("     %d. %s: %.0f points\n", i+1, player.Member, player.Score)
	}

	// ZREVRANGE - Get members by rank (descending order)
	fmt.Println("\n3. Getting top players with ZREVRANGE:")

	// Get top 3 players
	topPlayers, err := rdb.ZRevRangeWithScores(ctx, leaderboard, 0, 2).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("   Top 3 players:")
	for i, player := range topPlayers {
		fmt.Printf("     %d. %s: %.0f points\n", i+1, player.Member, player.Score)
	}

	// ZSCORE - Get score of specific member
	fmt.Println("\n4. Getting specific scores with ZSCORE:")
	aliceScore, err := rdb.ZScore(ctx, leaderboard, "alice").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Alice's score: %.0f\n", aliceScore)

	// ZRANK - Get rank of member (0-based, ascending)
	fmt.Println("\n5. Getting player ranks with ZRANK and ZREVRANK:")
	aliceRank, err := rdb.ZRank(ctx, leaderboard, "alice").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Alice's rank (ascending): %d\n", aliceRank)

	// ZREVRANK - Get rank of member (0-based, descending)
	aliceRevRank, err := rdb.ZRevRank(ctx, leaderboard, "alice").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Alice's rank (descending): %d (position from top)\n", aliceRevRank)

	// ZCARD - Get number of members
	fmt.Println("\n6. Getting leaderboard size with ZCARD:")
	playerCount, err := rdb.ZCard(ctx, leaderboard).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Total players: %d\n", playerCount)

	// ZINCRBY - Increment member score
	fmt.Println("\n7. Updating scores with ZINCRBY:")
	newScore, err := rdb.ZIncrBy(ctx, leaderboard, 300, "alice").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Alice's new score after +300: %.0f\n", newScore)

	// Check new rankings
	newTopPlayers, err := rdb.ZRevRangeWithScores(ctx, leaderboard, 0, 2).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("   Updated top 3:")
	for i, player := range newTopPlayers {
		fmt.Printf("     %d. %s: %.0f points\n", i+1, player.Member, player.Score)
	}

	// ZRANGEBYSCORE - Get members by score range
	fmt.Println("\n8. Getting players by score range with ZRANGEBYSCORE:")

	// Players with scores between 1500 and 2000
	midRangePlayers, err := rdb.ZRangeByScoreWithScores(ctx, leaderboard, &redis.ZRangeBy{
		Min: "1500",
		Max: "2000",
	}).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("   Players with scores 1500-2000:")
	for _, player := range midRangePlayers {
		fmt.Printf("     %s: %.0f points\n", player.Member, player.Score)
	}

	// ZCOUNT - Count members in score range
	fmt.Println("\n9. Counting players in score range with ZCOUNT:")
	count, err := rdb.ZCount(ctx, leaderboard, "1500", "2000").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Players with scores 1500-2000: %d\n", count)

	// ZREM - Remove members
	fmt.Println("\n10. Removing players with ZREM:")
	removed, err := rdb.ZRem(ctx, leaderboard, "eve").Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Removed %d player(s)\n", removed)

	// ZREMRANGEBYRANK - Remove by rank range
	fmt.Println("\n11. Removing bottom players with ZREMRANGEBYRANK:")
	removedByRank, err := rdb.ZRemRangeByRank(ctx, leaderboard, 0, 0).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("   Removed %d bottom player(s)\n", removedByRank)

	// Final leaderboard
	finalLeaderboard, err := rdb.ZRevRangeWithScores(ctx, leaderboard, 0, -1).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("   Final leaderboard:")
	for i, player := range finalLeaderboard {
		fmt.Printf("     %d. %s: %.0f points\n", i+1, player.Member, player.Score)
	}

	// Practical example: Time-series data (using timestamps as scores)
	fmt.Println("\n12. Practical example - Time-series data:")
	timeSeriesKey := "sensor:temperature"

	// Add temperature readings with timestamps as scores
	readings := []redis.Z{
		{Score: 1640995200, Member: "22.5"},
		{Score: 1640995260, Member: "23.1"},
		{Score: 1640995320, Member: "22.8"},
		{Score: 1640995380, Member: "23.4"},
		{Score: 1640995440, Member: "23.0"},
	}

	rdb.ZAdd(ctx, timeSeriesKey, readings...)

	// Get latest 3 readings
	latest, err := rdb.ZRevRangeWithScores(ctx, timeSeriesKey, 0, 2).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("   Latest 3 temperature readings:")
	for _, reading := range latest {
		fmt.Printf("     Timestamp %.0f: %s°C\n", reading.Score, reading.Member)
	}

	// Get readings in time range
	timeRange, err := rdb.ZRangeByScoreWithScores(ctx, timeSeriesKey, &redis.ZRangeBy{
		Min: "1640995200",
		Max: "1640995320",
	}).Result()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("   Readings in first 2 minutes:")
	for _, reading := range timeRange {
		fmt.Printf("     Timestamp %.0f: %s°C\n", reading.Score, reading.Member)
	}

	// Practical example: Priority queue
	fmt.Println("\n13. Practical example - Priority queue:")
	priorityQueue := "task:priority_queue"

	// Add tasks with priority scores (higher score = higher priority)
	tasks := []redis.Z{
		{Score: 1, Member: "backup_database"},
		{Score: 5, Member: "fix_critical_bug"},
		{Score: 3, Member: "deploy_feature"},
		{Score: 2, Member: "update_documentation"},
		{Score: 4, Member: "security_patch"},
	}

	rdb.ZAdd(ctx, priorityQueue, tasks...)

	// Process tasks by priority (highest first)
	fmt.Println("   Processing tasks by priority:")
	for i := 0; i < 3; i++ {
		// Get highest priority task
		highestPriority, err := rdb.ZRevRangeWithScores(ctx, priorityQueue, 0, 0).Result()
		if err != nil || len(highestPriority) == 0 {
			break
		}

		task := highestPriority[0]
		fmt.Printf("     Processing (priority %.0f): %s\n", task.Score, task.Member)

		// Remove processed task
		rdb.ZRem(ctx, priorityQueue, task.Member)
	}

	// Cleanup
	fmt.Println("\n14. Cleanup:")
	rdb.Del(ctx, leaderboard, timeSeriesKey, priorityQueue)
	fmt.Println("   Cleaned up sorted set examples ✓")
}
