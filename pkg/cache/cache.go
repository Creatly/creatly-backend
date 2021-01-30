package cache

type Cache interface {
	Set(key, value interface{}, ttl int64) error
	Get(key interface{}) (interface{}, error)
}
