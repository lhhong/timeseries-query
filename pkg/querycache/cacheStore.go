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

func NewCacheStore(conf *config.RedisConfig) *CacheStore {

	cs := &CacheStore{
		env: conf.Env,
	}
	cs.InitConn(conf.Hostname, conf.Port)
	return cs
}

func (cs *CacheStore) InitConn(hostname string, port int) {

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
	cs.redisPool = pool
}
