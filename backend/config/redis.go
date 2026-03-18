package config

import (
	"haveYouWorkedOutToday/global"
	"log"

	"github.com/go-redis/redis"
)

func InitRedis() {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis, got error: %s", err.Error())
		log.Printf("Continuing without Redis...")
		return
	}

	global.RedisDB = RedisClient
	log.Printf("Redis connected successfully")
}