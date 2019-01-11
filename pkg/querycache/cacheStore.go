package querycache

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/lhhong/timeseries-query/pkg/config"
)

type CacheStore struct {
	env   string
	redis redis.Conn
}

func NewCacheStore(conf *config.RedisConfig) *CacheStore {

	cs := &CacheStore{
		env: conf.Env,
	}
	cs.InitConn(conf.Hostname, conf.Port)
	return cs
}

func (cs *CacheStore) InitConn(hostname string, port int) {

	conn, err := redis.DialURL(fmt.Sprintf("redis://%s:%d", hostname, port))
	if err != nil {
		log.Println("Cannot connect to redis")
		log.Panicln(err)
	}
	cs.redis = conn
}
