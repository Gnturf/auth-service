package initializer

import (
	"auth-service/config"
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedisConn(config *config.AppConfig) *redis.Client {
	// Create a context for Redis operations
	ctx := context.Background()
	// redisAddr := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)

	// Initialize the Redis client
	client := redis.NewClient(&redis.Options{
		Addr		: "localhost:6379", 				// Redis server address
		Password: "",               // No password (default)
		DB			: config.RedisDB,   // Default DB index
	})

	// Ping Redis to check if it's reachable
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis!")
	return client
}