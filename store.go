package store

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/patrickmn/go-cache"
	"github.com/urionz/store/driver"
)

const (
	// 默认存储时间为1小时
	DefaultExpiration = time.Hour
	// 默认清理任务周期为1小时
	DefaultCleanupInterval = time.Hour
	// map存储类型
	TypeMemory = "memory"
	// redis存储类型
	TypeRedis     = "redis"
	DefaultPrefix = "cache_"
	DefaultAddr   = "localhost:6379"
)

type Store interface {
	// 通过key从存储中检索出对象
	Get(key string) interface{}
	// 通过key从存储中检索出对象并赋值给指针
	GetScan(key string, scanner interface{}) error
	// 存储一个对象在设定的时间内过期
	Put(key string, value interface{}, minutes time.Duration) bool
	// 判断是否存在此key的对象
	Has(key string) bool
	// 永久存储一个对象
	Forever(key, value string) bool
	// 删除一个对象
	Forget(key string) bool
	// 递增一个对象的值
	Increment(key string, step int64) error
	// 递减一个对象的值
	Decrement(key string, step int64) error
	// 清空存储
	Flush() bool
}

type Container struct {
	Driver          Store
	Prefix          string
	Expiration      time.Duration
	CleanupInterval time.Duration
	Addr            string
	DB              int
	Password        string
}

func (container *Container) GetStorage() Store {
	return container.Driver
}

func (container *Container) Get(key string) interface{} {
	return container.Driver.Get(container.getPrefixKey(key))
}

func (container *Container) GetScan(key string, scanner interface{}) error {
	return container.Driver.GetScan(container.getPrefixKey(key), scanner)
}

func (container *Container) GetDefault(key string, def interface{}) interface{} {
	if value := container.Get(key); value != nil {
		return value
	}
	return def
}

func (container *Container) GetScanDefault(key string, scanner interface{}, def interface{}) interface{} {
	if value := container.GetScan(key, scanner); value != nil {
		return value
	}
	return def
}

func (container *Container) Many(keys []string) map[string]interface{} {
	data := make(map[string]interface{})
	for _, key := range keys {
		value := container.Get(key)
		if value == nil {
			value = nil
		}
		data[key] = value
	}
	return data
}

func (container *Container) Put(key string, value interface{}, expiration time.Duration) bool {
	return container.Driver.Put(container.getPrefixKey(key), value, expiration)
}

func (container *Container) PutDefault(key string, value interface{}) bool {
	return container.Put(key, value, container.Expiration)
}

func (container *Container) PutMany(kv map[string]interface{}, expiration time.Duration) bool {
	for key, value := range kv {
		if container.Put(key, value, expiration) {
			continue
		} else {
			return false
		}
	}
	return true
}

func (container *Container) PutManyDefault(kv map[string]interface{}) bool {
	return container.PutMany(kv, container.Expiration)
}

func (container *Container) Add(key string, value interface{}, expiration time.Duration) bool {
	if !container.Has(key) {
		return container.Put(key, value, expiration)
	}
	return false
}

func (container *Container) AddDefault(key string, value interface{}) bool {
	return container.Add(key, value, container.Expiration)
}

func (container *Container) Increment(key string, step int64) error {
	return container.Driver.Increment(container.getPrefixKey(key), step)
}

func (container *Container) Decrement(key string, step int64) error {
	return container.Driver.Decrement(container.getPrefixKey(key), step)
}

func (container *Container) Forever(key, value string) bool {
	return container.Driver.Forever(container.getPrefixKey(key), value)
}

func (container *Container) Forget(key string) bool {
	return container.Driver.Forget(container.getPrefixKey(key))
}

func (container *Container) Has(key string) bool {
	return container.Driver.Has(container.getPrefixKey(key))
}

func (container *Container) Flush() bool {
	return container.Driver.Flush()
}

func (container *Container) GetPrefix() string {
	return container.Prefix
}

func (container *Container) getPrefixKey(key string) string {
	return container.Prefix + key
}

// 通过存储类型创建存储实例
func New(storeType string, container Container) *Container {
	if container.Prefix == "" {
		container.Prefix = DefaultPrefix
	}
	switch storeType {
	case TypeMemory:
		if container.Expiration == 0 {
			container.Expiration = DefaultExpiration
		}
		if container.CleanupInterval == 0 {
			container.CleanupInterval = DefaultCleanupInterval
		}
		if container.Driver == nil {
			container.Driver = driver.NewMemoryStore(cache.New(container.Expiration, container.CleanupInterval))
		}
	case TypeRedis:
		if container.Addr == "" {
			container.Addr = DefaultAddr
		}
		if container.Driver == nil {
			container.Driver = driver.NewRedisStore(redis.NewClient(&redis.Options{
				Addr:     container.Addr,
				Password: container.Password,
				DB:       container.DB,
			}))
		}
	}
	return &container
}

// 创建存储实例通过默认设置
//func NewDefault(instance interface{}) Store {
//	return New(TypeRedis, DefaultExpiration, DefaultCleanupInterval, instance)
//}
