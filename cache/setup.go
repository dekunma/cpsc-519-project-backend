package cache

import (
	"github.com/redis/go-redis/v9"
	"os"
)

var RDB *redis.Client

func ConnectRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ENDPOINT"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
