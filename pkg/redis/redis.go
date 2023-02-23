package redis

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "",
		DB:       0,
	})
}

func SaveToRedis(key string, value string) error {
	err := redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetFromRedis(key string) (string, error) {
	val, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func KeyExists(key string) (bool, error) {
	val, err := redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val == 1, nil
}
