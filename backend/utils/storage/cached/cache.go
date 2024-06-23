package cached

import (
	"github.com/patrickmn/go-cache"
	"sync"
)

// 表示redis随机缓存的时间范围
const redisRandomScope = 1

var cacheMaps = make(map[string]*cache.Cache)

var m = new(sync.Mutex)

type cachedItem interface {
	GetId() string
	IsDirty() bool
}
