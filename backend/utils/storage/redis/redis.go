package redis

import (
	"fmt"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constans/config"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func init() {
	addrs := fmt.Sprintf("%s:%s", config.Conf.Redis.Host, config.Conf.Redis.Port)
	Client = redis.NewClient(&redis.Options{
		Addr:     addrs,
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.Db,
	})

	if err := redisotel.InstrumentTracing(Client); err != nil {
		panic(err)
	}

	if err := redisotel.InstrumentMetrics(Client); err != nil {
		panic(err)
	}
}
