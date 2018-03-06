package driver

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type MemoryStore struct {
	instance        *cache.Cache
	cleanupInterval time.Duration
}

func NewMemoryStore(instance *cache.Cache) *MemoryStore {
	return &MemoryStore{
		instance: instance,
	}
}

func (memory *MemoryStore) Get(key string) interface{} {
	if value, ok := memory.instance.Get(key); ok {
		return value
	}
	return nil
}
func (memory *MemoryStore) GetScan(key string, scanner interface{}) error {
	if value := memory.Get(key); value != nil {
		switch v := scanner.(type) {
		case nil:
			return fmt.Errorf("redis: Scan(nil)")
		case *string:
			*v = value.(string)
		case *[]byte:
			*v = value.([]byte)
		case *int:
			*v = value.(int)
		case *int8:
			*v = value.(int8)
		case *int16:
			*v = value.(int16)
		case *int32:
			*v = value.(int32)
		case *int64:
			*v = value.(int64)
		case *uint:
			*v = value.(uint)
		case *uint8:
			*v = value.(uint8)
		case *uint16:
			*v = value.(uint16)
		case *uint32:
			*v = value.(uint32)
		case *uint64:
			*v = value.(uint64)
		case *float32:
			*v = value.(float32)
		case *float64:
			*v = value.(float64)
		case *bool:
			*v = value.(bool)
		case *time.Time:
			*v = value.(time.Time)
		}
		return nil
	}
	return fmt.Errorf("get the %s fail", key)
}

func (memory *MemoryStore) Put(key string, value interface{}, expiration time.Duration) bool {
	memory.instance.Set(key, value, expiration)
	return true
}

func (memory *MemoryStore) Increment(key string, step int64) error {
	err := memory.instance.Increment(key, step)
	return err
}

func (memory *MemoryStore) Decrement(key string, step int64) error {
	return memory.instance.Decrement(key, step)
}

func (memory *MemoryStore) Forever(key, value string) bool {
	return memory.Put(key, value, -1)
}

func (memory *MemoryStore) Forget(key string) bool {
	memory.instance.Delete(key)
	return true
}

func (memory *MemoryStore) Has(key string) bool {
	value := memory.Get(key)
	if value != nil {
		return true
	}
	return false
}

func (memory *MemoryStore) Flush() bool {
	memory.instance.Flush()
	return true
}
