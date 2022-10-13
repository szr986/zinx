package db

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
)

var ctx = context.Background()
var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: "",
		DB:       0,
	})
	fmt.Println("Redis init success")
	RedisClient.SAdd("test", "ok")
	Res := RedisClient.SMembers("test").Val()
	fmt.Println("test  redis:", Res)

}

func GetRedisClient() *redis.Client {
	return RedisClient
}
