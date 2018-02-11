package store

import (
	"time"

	"github.com/urionz/store/driver"
)

const (
	// 默认存储时间为1小时
	DefaultExpiration = time.Hour
	// 默认清理任务周期为1小时
	DefaultCleanupInterval = time.Second
	// map存储类型
	TypeMemory = "memory"
	// redis存储类型
	TypeRedis = "redis"
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
	driver     Store
	prefix     string
	expiration time.Duration
}

func (container *Container) GetStorage() Store {
	return container.driver
}

func (container *Container) Get(key string) interface{} {
	return container.driver.Get(container.getPrefixKey(key))
}

func (container *Container) GetScan(key string, scanner interface{}) error {
	return container.driver.GetScan(container.getPrefixKey(key), scanner)
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
	return container.driver.Put(container.getPrefixKey(key), value, expiration)
}

func (container *Container) PutDefault(key string, value interface{}) bool {
	return container.Put(key, value, container.expiration)
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
	return container.PutMany(kv, container.expiration)
}

func (container *Container) Add(key string, value interface{}, expiration time.Duration) bool {
	if !container.Has(key) {
		return container.Put(key, value, expiration)
	}
	return false
}

func (container *Container) AddDefault(key string, value interface{}) bool {
	return container.Add(key, value, container.expiration)
}

func (container *Container) Increment(key string, step int64) error {
	return container.driver.Increment(container.getPrefixKey(key), step)
}

func (container *Container) Decrement(key string, step int64) error {
	return container.driver.Decrement(container.getPrefixKey(key), step)
}

func (container *Container) Forever(key, value string) bool {
	return container.driver.Forever(container.getPrefixKey(key), value)
}

func (container *Container) Forget(key string) bool {
	return container.driver.Forget(container.getPrefixKey(key))
}

func (container *Container) Has(key string) bool {
	return container.driver.Has(container.getPrefixKey(key))
}

func (container *Container) Flush() bool {
	return container.driver.Flush()
}

func (container *Container) GetPrefix() string {
	return container.prefix
}

func (container *Container) getPrefixKey(key string) string {
	return container.prefix + key
}

// 通过存储类型创建存储实例
func New(storeType string, expiration, cleanupInterval time.Duration, instance interface{}) Store {
	switch storeType {
	case TypeMemory:
		return &Container{
			driver:     driver.NewMemoryStore(cleanupInterval, instance),
			prefix:     "memory_",
			expiration: expiration,
		}
	case TypeRedis:
		return &Container{
			driver:     driver.NewRedisStore(instance),
			prefix:     "redis_",
			expiration: expiration,
		}
	}
	return nil
}

// 创建存储实例通过默认设置
func NewDefault(instance interface{}) Store {
	return New(TypeRedis, DefaultExpiration, DefaultCleanupInterval, instance)
}
