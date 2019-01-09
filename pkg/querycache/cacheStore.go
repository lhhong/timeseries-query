package querycache

type CacheStore interface {
	Init(...[]interface{})
	Get(string) string
	Put(string, string)
}
