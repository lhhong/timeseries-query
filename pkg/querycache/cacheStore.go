package querycache

import (
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/lhhong/timeseries-query/pkg/config"
)

type CacheStore struct {
	env       string
	redisPool *redis.Pool
}

func InitCacheStore(conf *config.RedisConfig) *CacheStore {

	return NewCacheStore(conf.Env, conf.Hostname, conf.Port)
}

func NewCacheStore(env string, hostname string, port int) *CacheStore {

	cs := &CacheStore{
		env:       env,
		redisPool: initConnPool(hostname, port),
	}
	return cs
}

func initConnPool(hostname string, port int) *redis.Pool {

	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
		},
	}

	// Test connection
	conn := pool.Get()
	_, err := conn.Do("PING")
	if err != nil {
		log.Println("Cannot connect to redis")
		log.Panicln(err)
	}
	return pool
}

func (cs CacheStore) formatKey(key string) string {
	return fmt.Sprintf("%s/%s", cs.env, key)
}
