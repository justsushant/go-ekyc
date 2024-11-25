package store

type CacheStore interface {
	GetObject(key string) (string, error)
	SetObject(key, val string) error
}
