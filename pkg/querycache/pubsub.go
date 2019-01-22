package querycache

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

func (cs CacheStore) Subscribe(channel string, onStart func(redis.Conn, chan []byte)) {
	conn := cs.redisPool.Get()
	psc := redis.PubSubConn{Conn: conn}
	psc.Subscribe(cs.formatKey(channel))

	dataChan := make(chan []byte)
	go func(conn redis.Conn, dataChan chan []byte) {
		for {
			switch m := psc.Receive().(type) {
			case error:
				log.Println("error", m)
			case redis.Message:
				dataChan <- m.Data
			case redis.Subscription:
				if m.Count == 0 {
					// Unsubscribed
					return
				}
			}
		}
	}(conn, dataChan)
	go onStart(conn, dataChan)
}

func (cs CacheStore) Unsubscribe(conn redis.Conn) {
	conn.Do("UNSUBSCRIBE")
	conn.Close()
}

func (cs CacheStore) Publish(channel string, data []byte) {
	conn := cs.redisPool.Get()
	defer conn.Close()

	conn.Do("PUBLISH", cs.formatKey(channel), data)
}
