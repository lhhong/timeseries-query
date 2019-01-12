package querycache

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func (cs CacheStore) GetBytes(key string) ([]byte, error) {
	conn := cs.redisPool.Get()
	defer conn.Close()
	res, err := redis.Bytes(conn.Do("GET", fmt.Sprintf("%s/%s", cs.env, key)))
	if err != nil {
		return res, err
	}
	return res, nil
}

func (cs CacheStore) SetBytes(key string, val []byte) error {

	//TODO export to params
	cacheExpiry := 30 * 60 // 30 min

	conn := cs.redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", fmt.Sprintf("%s/%s", cs.env, key), val, "EX", cacheExpiry)
	if err != nil {
		return err
	}
	return nil
}

func (cs CacheStore) Delete(key string) (int, error) {
	conn := cs.redisPool.Get()
	defer conn.Close()

	res, err := redis.Int(conn.Do("DEL", fmt.Sprintf("%s/%s", cs.env, key)))
	if err != nil {
		return res, err
	}
	return res, nil
}

func (cs CacheStore) GetsetBytes(key string, val []byte) ([]byte, error) {

	conn := cs.redisPool.Get()
	defer conn.Close()

	res, err := redis.Bytes(conn.Do("GETSET", fmt.Sprintf("%s/%s", cs.env, key), val))
	if err != nil {
		return res, err
	}
	return res, nil
}
