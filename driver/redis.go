package driver

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisStore struct {
	instance *redis.Client
}

func NewRedisStore(instance *redis.Client) *RedisStore {
	return &RedisStore{
		instance: instance,
	}
}

func (redis *RedisStore) Get(key string) interface{} {
	value, err := redis.instance.Get(key).Result()
	if err != nil {
		return nil
	}
	return value
}

func (redis *RedisStore) GetScan(key string, scanner interface{}) error {
	aaa := redis.instance.Get(key).Scan(scanner)
	return aaa
}

func (redis *RedisStore) Put(key string, value interface{}, expiration time.Duration) bool {
	if err := redis.instance.Set(key, value, expiration).Err(); err != nil {
		return false
	}
	return true
}

func (redis *RedisStore) Increment(key string, step int64) error {
	return redis.instance.IncrBy(key, step).Err()
}

func (redis *RedisStore) Decrement(key string, step int64) error {
	return redis.instance.DecrBy(key, step).Err()
}

func (redis *RedisStore) Forever(key, value string) bool {
	if result, _ := redis.instance.Set(key, value, 0).Result(); result == "" {
		return false
	}
	return true
}

func (redis *RedisStore) Forget(key string) bool {
	if result, _ := redis.instance.Del(key).Result(); result == 0 {
		return false
	}
	return true
}

func (redis *RedisStore) Has(key string) bool {
	if result, _ := redis.instance.Exists(key).Result(); result == 0 {
		return false
	}
	return true
}

func (redis *RedisStore) Flush() bool {
	if err := redis.instance.FlushDB().Err(); err != nil {
		return false
	}
	return true
}
