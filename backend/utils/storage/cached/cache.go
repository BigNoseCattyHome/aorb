package cached

import (
	"context"
	"errors"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constans/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/redis"
	"github.com/patrickmn/go-cache"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// 表示redis随机缓存的时间范围
const redisRandomScope = 1

var cacheMaps = make(map[string]*cache.Cache)

var m = new(sync.Mutex)

type cachedItem interface {
	GetId() uint32
	IsDirty() bool
}

func ScanGet(ctx context.Context, key string, obj interface{}, collectionName string) (bool, error) {
	ctx, span := tracing.Tracer.Start(ctx, "Cached-GetFromScanCache")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("Cached.GetFromScanCache").WithContext(ctx)
	key = config.Conf.Redis.Prefix + key

	c := getOrCreateCache(key)
	wrappedObj := obj.(cachedItem)
	key = key + strconv.FormatUint(uint64(wrappedObj.GetId()), 10)
	if x, found := c.Get(key); found {
		dstVal := reflect.ValueOf(obj)
		dstVal.Elem().Set(x.(reflect.Value))
		return true, nil
	}

	// 缓存没有命中
	logger.WithFields(logrus.Fields{
		"key": key,
	}).Infof("Missed local memory cached")

	if err := redis.Client.HGetAll(ctx, key).Scan(obj); err != nil {
		if err != redis2.Nil {
			logger.WithFields(logrus.Fields{
				"key": key,
				"err": err,
			}).Errorf("Redis error when find struct")
			logging.SetSpanError(span, err)
			return false, err
		}
	}

	// 如果redis命中，那么存储到本地缓存然后返回
	if wrappedObj.IsDirty() {
		logger.WithFields(logrus.Fields{
			"key": key,
		}).Infof("Redis hit the key")
		c.Set(key, reflect.ValueOf(obj).Elem(), cache.DefaultExpiration)
		return true, nil
	}

	// redis没有命中，回调到数据库
	logger.WithFields(logrus.Fields{
		"key": key,
	}).Warnf("Missed Redis Cached")

	collections := database.MongoDbClient.Database("aorb").Collection(collectionName)
	result := collections.FindOne(ctx, obj)
	if result == nil {
		logger.WithFields(logrus.Fields{
			"key": key,
		}).Warnf("Missed DB obj, seems wrong key")
		return false, errors.New("Missed DB obj, seems wrong key")
	}

	if result := redis.Client.HSet(ctx, key, obj); result.Err() != nil {
		logger.WithFields(logrus.Fields{
			"err": result.Err(),
			"key": key,
		}).Errorf("Redis error when set struct info")
		logging.SetSpanError(span, result.Err())
		return false, nil
	}

	c.Set(key, reflect.ValueOf(obj).Elem(), cache.DefaultExpiration)
	return true, nil
}

// ScanTagDelete 将缓存值标记为删除，下次从cache读取时会回调到数据库
func ScanTagDelete(ctx context.Context, key string, obj interface{}) {
	ctx, span := tracing.Tracer.Start(ctx, "Cached-ScanTagDelete")
	defer span.End()
	logging.SetSpanWithHostname(span)
	key = config.Conf.Redis.Prefix + key

	redis.Client.HDel(ctx, key)

	c := getOrCreateCache(key)
	wrappedObj := obj.(cachedItem)
	key = key + strconv.FormatUint(uint64(wrappedObj.GetId()), 10)
	c.Delete(key)
}

// ScanWriteCache 将数据写入缓存，如果state为false，那么只写入本地缓存，否则同时写入redis
func ScanWriteCache(ctx context.Context, key string, obj interface{}, state bool) (err error) {
	ctx, span := tracing.Tracer.Start(ctx, "Cached-ScanWriteCache")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("Cached.ScanWriteCache").WithContext(ctx)
	key = config.Conf.Redis.Prefix + key

	wrappedObj := obj.(cachedItem)
	key = key + strconv.FormatUint(uint64(wrappedObj.GetId()), 10)
	c := getOrCreateCache(key)
	c.Set(key, reflect.ValueOf(obj).Elem(), cache.DefaultExpiration)

	if state {
		if err = redis.Client.HGetAll(ctx, key).Scan(obj); err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
				"key": key,
			}).Errorf("Redis error when find struct info")
			logging.SetSpanError(span, err)
			return
		}
	}
	return
}

// 读取字符串缓存，其中找到了返回true，没找到或者异常返回false
func Get(ctx context.Context, key string) (string, bool, error) {
	ctx, span := tracing.Tracer.Start(ctx, "Cached-Get")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("Cached.GetFromStringCache").WithContext(ctx)
	key = config.Conf.Redis.Prefix + key

	c := getOrCreateCache("strings")
	if x, found := c.Get(key); found {
		return x.(string), true, nil
	}

	// 缓存没有命中，回调数据库
	logger.WithFields(logrus.Fields{
		"key": key,
	}).Infof("Missed local memory cached")

	var result *redis2.StringCmd
	if result = redis.Client.Get(ctx, key); result.Err() != nil && result.Err() != redis2.Nil {
		logger.WithFields(logrus.Fields{
			"err":    result.Err(),
			"string": key,
		}).Errorf("Redis error when find string")
		logging.SetSpanError(span, result.Err())
		return "", false, nil
	}

	value, err := result.Result()
	switch {
	case err == redis2.Nil:
		return "", false, nil
	case err != nil:
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Err when write Redis")
		logging.SetSpanError(span, err)
		return "", false, err
	default:
		c.Set(key, value, cache.DefaultExpiration)
		return value, true, nil
	}
}

// GetWithFunc 从缓存中获取字符串，如果没有就调用func函数获取
func GetWithFunc(ctx context.Context, key string, f func(ctx context.Context, key string) (string, error)) (string, error) {
	ctx, span := tracing.Tracer.Start(ctx, "Cached-GetFromStringCacheWithFunc")
	defer span.End()
	logging.SetSpanWithHostname(span)
	value, ok, err := Get(ctx, key)

	if err != nil {
		return "", err
	}
	if ok {
		return value, nil
	}

	// 如果不存在，调用函数
	value, err = f(ctx, key)
	if err != nil {
		return "", err
	}
	Write(ctx, key, value, true)
	return value, nil
}

// Write 写入字符串缓存，如果state为false，那么只写入本地缓存，否则同时写入redis
func Write(ctx context.Context, key string, value string, state bool) {
	ctx, span := tracing.Tracer.Start(ctx, "Cached-SetStringCache")
	defer span.End()
	logging.SetSpanWithHostname(span)
	key = config.Conf.Redis.Prefix + key

	c := getOrCreateCache("strings")
	c.Set(key, value, cache.DefaultExpiration)

	if state {
		redis.Client.Set(ctx, key, value, 120*time.Hour+time.Duration(rand.Intn(redisRandomScope))*time.Second)
	}
}

// TagDelete 删除字符串缓存
// 与ScanTagDelete不一样，这里是直接删除键值对，而ScanTagDelete是打标记然后等待下一次访问的时候更新
func TagDelete(ctx context.Context, key string) {
	ctx, span := tracing.Tracer.Start(ctx, "Cached-DeleteStringCache")
	defer span.End()
	logging.SetSpanWithHostname(span)
	key = config.Conf.Redis.Prefix + key

	c := getOrCreateCache("strings")
	c.Delete(key)

	redis.Client.Del(ctx, key)
}

// 从本地缓存中获取指定的键值对的值(一个缓存实例)，如果没有就新建一个
func getOrCreateCache(name string) *cache.Cache {
	cc, ok := cacheMaps[name]
	if !ok {
		m.Lock()
		defer m.Unlock()
		cc, ok := cacheMaps[name]
		if !ok {
			cc = cache.New(5*time.Minute, 10*time.Minute)
			cacheMaps[name] = cc
			return cc
		}
		return cc
	}
	return cc
}

// CacheAndRedisGet 从内存缓存和redis中获取数据
func CacheAndRedisGet(ctx context.Context, key string, obj interface{}) (bool, error) {
	ctx, span := tracing.Tracer.Start(ctx, "CacheAndRedisGet")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("CacheAndRedisGet").WithContext(ctx)
	key = config.Conf.Redis.Prefix + key

	c := getOrCreateCache(key)
	wrappedObj := obj.(cachedItem)
	key = key + strconv.FormatUint(uint64(wrappedObj.GetId()), 10)
	if x, found := c.Get(key); found {
		dstVal := reflect.ValueOf(obj)
		dstVal.Elem().Set(x.(reflect.Value))
		return true, nil
	}

	// 缓存没有命中，回调redis
	logger.WithFields(logrus.Fields{
		"key": key,
	}).Infof("Missed local memory cached")

	if err := redis.Client.HGetAll(ctx, key).Scan(obj); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
			"key": key,
		}).Errorf("Redis error when find struct")
		logging.SetSpanError(span, err)
		return false, err
	}

	// redis命中，存储到本地缓存然后返回
	if wrappedObj.IsDirty() {
		logger.WithFields(logrus.Fields{
			"key": key,
		}).Infof("Redis hit the key")
		c.Set(key, reflect.ValueOf(obj).Elem(), cache.DefaultExpiration)
		return true, nil
	}

	logger.WithFields(logrus.Fields{
		"key": key,
	}).Warnf("Missed Redis Cached")

	return false, nil
}

func ActionRedisSync(time time.Duration, f func(client redis2.UniversalClient) error) {
	go func() {
		daemon := NewTick(time, f)
		daemon.Start()
	}()
}
