package cache

type Cache interface {
	Set(key, value interface{}) error
	Get(key interface{}) (interface{}, error)
}
