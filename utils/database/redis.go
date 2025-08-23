package database

import (
	"fmt"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	*redis.Client
}

var store *RedisStore
var once2 sync.Once

func GetRedisStore() *RedisStore {
	once2.Do(func() {
		store = &RedisStore{
			redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
				Password: "", // no password set
				DB:       0,  // use default DB
			}),
		}
	})
	return store
}
