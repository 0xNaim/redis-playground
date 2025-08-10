package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"redis-playground/config"
	"redis-playground/examples"
	"strings"
)

func main() {
	// Initialize Redis client
	rdb := config.InitRedis()
	defer rdb.Close()

	// Test connection
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	fmt.Println("Welcome to Redis Playground with Go!")
	fmt.Println("=====================================")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		showMenu()
		fmt.Print("Enter your choice: ")
		if !scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			examples.RunStringExamples(rdb)
		case "2":
			examples.RunListExamples(rdb)
		case "3":
			examples.RunSetsExamples(rdb)
		case "4":
			examples.RunSortedSetsExamples(rdb)
		case "5":
			examples.RunHashesExamples(rdb)
		case "0":
			fmt.Println("Exiting Redis Playground. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}

		fmt.Println("\nPress Enter to continue...")
		scanner.Scan()
	}
}

func showMenu() {
	fmt.Println("\n Choose an option:")
	fmt.Println("1. Run String Examples")
	fmt.Println("2. Run List Examples")
	fmt.Println("3. Run Set Examples")
	fmt.Println("4. Run Sorted Set Examples")
	fmt.Println("5. Run Hash Examples")
	fmt.Println("0. Exit")
}
