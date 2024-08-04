package redis

import (
	"fmt"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func init() {
	addrs := fmt.Sprintf("%s:%s", config.Conf.Redis.Host, config.Conf.Redis.Port)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addrs,
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.Db,
	})

	if err := redisotel.InstrumentTracing(RedisClient); err != nil {
		panic(err)
	}
}
