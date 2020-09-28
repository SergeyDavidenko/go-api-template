package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/config"
	log "github.com/sirupsen/logrus"
)

// Cache is interface structure
type Cache struct {
	rdb *redis.Client
	ctx context.Context
}

// New func implements the cache interface
func New() *Cache {
	return &Cache{}
}

// Init client Cache.
func (c *Cache) Init() error {
	redisHost := fmt.Sprintf("%s:%d", config.Conf.Storage.RedisHost, config.Conf.Storage.RedisPort)
	log.Debug("redis host is - ", redisHost)
	c.ctx = context.Background()
	c.rdb = redis.NewClient(&redis.Options{
		Addr: redisHost,
		Password: config.Conf.Storage.RedisPassrod,
		DB:   config.Conf.Storage.RedisDB,
	})
	_, err := c.rdb.Ping(c.ctx).Result()
	return err
}

// Close Cache.
func (c *Cache) Close() error {
	err := c.rdb.Close()
	return err
}

// Get ...
func (c *Cache) Get(key string) (interface{}, error) {
	val, err := c.rdb.Get(c.ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// Set ...
func (c *Cache) Set(key string, value interface{}, durations time.Duration) error {
	err := c.rdb.Set(c.ctx, key, "value", durations).Err()
	if err != nil {
		panic(err)
	}
}
