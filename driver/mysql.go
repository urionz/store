package driver

import (
	"time"

	"github.com/jinzhu/gorm"
)

type MysqlStore struct {
	instance          *gorm.DB
	prefix            string
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
}

func NewMysqlStore(defaultExpiration, cleanupInterval time.Duration) *MysqlStore {
	return &MysqlStore{
		prefix:            "mysql_",
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}
}

func (mysql *MysqlStore) Get(key string) interface{} {
	return nil
}

func (mysql *MysqlStore) GetScan(key string, scanner interface{}) error {
	return nil
}

func (mysql *MysqlStore) Put(key string, value interface{}, minutes time.Duration) bool {
	return false
}

func (mysql *MysqlStore) Increment(key string, step int64) error {
	return nil
}

func (mysql *MysqlStore) Decrement(key string, step int64) error {
	return nil
}

func (mysql *MysqlStore) Forever(key, value string) bool {
	return false
}

func (mysql *MysqlStore) Forget(key string) bool {
	return false
}

func (mysql *MysqlStore) Has(key string) bool {
	return false
}

func (mysql *MysqlStore) Flush() bool {
	return false
}

func (mysql *MysqlStore) GetPrefix() string {
	return mysql.prefix
}

func (mysql *MysqlStore) CleanupExpired() bool {
	return true
}
