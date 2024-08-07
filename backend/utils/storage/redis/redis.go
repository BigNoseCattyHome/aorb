package redis

import (
	"fmt"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

var RedisCommentClient *redis.Client
var RedisMessageClient *redis.Client

func init() {
	addrs := fmt.Sprintf("%s:%s", config.Conf.Redis.Host, config.Conf.Redis.Port)
	RedisCommentClient = redis.NewClient(&redis.Options{
		Addr:     addrs,
		Password: config.Conf.Redis.Password,
		//DB:       config.Conf.Redis.Db,
		DB: 0,
	})

	RedisMessageClient = redis.NewClient(&redis.Options{
		Addr:     addrs,
		Password: config.Conf.Redis.Password,
		DB:       1,
	})

	if err := redisotel.InstrumentTracing(RedisCommentClient); err != nil {
		panic(err)
	}

	if err := redisotel.InstrumentTracing(RedisMessageClient); err != nil {
		panic(err)
	}
}
