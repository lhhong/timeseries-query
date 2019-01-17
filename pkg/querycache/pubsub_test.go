package querycache

import (
	"testing"

	"github.com/gomodule/redigo/redis"
)

func TestPubsub(t *testing.T) {
	cs := NewCacheStore("test", "localhost", 6379)

	publishData := "testData"

	onStart := func(conn redis.Conn, dataChan chan []byte) {
		cs.Publish("testChan", []byte("testData"))
		res := string(<-dataChan)
		if res != publishData {
			t.Errorf("res should be %s, got %s", publishData, res)
		}
		cs.Unsubscribe(conn)
	}

	cs.Subscribe("testChan", onStart)

}
