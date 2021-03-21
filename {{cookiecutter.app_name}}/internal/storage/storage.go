package storage

import (
	"time"
)

var (
	// StorageDB ...
	StorageDB Storage
	{% if cookiecutter.use_redis == "y" %}
	// CacheRedis ...
	CacheRedis Cache
	{% endif %}
)

// Storage interface
type Storage interface {
	Init() error
	ShowVersion() string
	Close() error
}


{% if cookiecutter.use_redis == "y" %}
// Cache interface
type Cache interface {
	Init() error
	Get(string) (interface{}, error)
	Set(string, interface{}, time.Duration) error
	Close() error
}
{% endif %}