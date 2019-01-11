package querycache

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func (cs CacheStore) GetBytes(key string) ([]byte, error) {
	conn := cs.redisPool.Get()
	res, err := redis.Bytes(conn.Do("GET", fmt.Sprintf("%s/%s", cs.env, key)))
	if err != nil {
		return nil, err
	}
	return res, nil
}
