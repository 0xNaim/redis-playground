package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RunPubSub demonstrates Redis Pub/Sub functionality
func RunPubSub(rdb *redis.Client) {
	fmt.Println("\n Pub/Sub Example")
	fmt.Println("===================")

	ctx := context.Background()
	channel := "chat:room1"

	// Simple publisher and subscriber demo using goroutines
	fmt.Println("1. Subscribing to channel and publishing messages:")

	pubsub := rdb.Subscribe(ctx, channel)
	defer pubsub.Close()

	// Start subscriber
	done := make(chan struct{})
	go func() {
		for i := 0; i < 3; i++ {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				fmt.Println("   Subscriber error:", err)
				continue
			}
			fmt.Printf("   Subscriber received: %s\n", msg.Payload)
		}
		close(done)
	}()

	// Publisher
	time.Sleep(200 * time.Millisecond)
	for i := 1; i <= 3; i++ {
		payload := fmt.Sprintf("Hello %d from publisher!", i)
		rdb.Publish(ctx, channel, payload)
		fmt.Printf("   Publisher sent: %s\n", payload)
		time.Sleep(100 * time.Millisecond)
	}
	<-done

	// Practical example: Real-time notifications
	fmt.Println("\n2. Practical example - Real-time notification system:")
	notifyChan := "notifications"
	notifyPubSub := rdb.Subscribe(ctx, notifyChan)
	defer notifyPubSub.Close()
	go func() {
		msg, err := notifyPubSub.ReceiveMessage(ctx)
		if err == nil {
			fmt.Printf("   Notification received: %s\n", msg.Payload)
		}
	}()
	rdb.Publish(ctx, notifyChan, "You have a new follower!")
	time.Sleep(200 * time.Millisecond)

	// Cleanup (no actual cleanup needed for pub/sub channels)
	fmt.Println("\n3. Pub/Sub demo complete ✓")
}
