package cache

import (
	"camphor-willow-god/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func init() {
	redisConfig := config.ApplicationConfig.Redis
	options, err := redis.ParseURL(redisConfig.URL)
	if err != nil {
		panic(err)
	}
	RedisClient = redis.NewClient(options)
}
