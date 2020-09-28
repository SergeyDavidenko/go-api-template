package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/config"
	log "github.com/sirupsen/logrus"
)

// Cache is interface structure
type Cache struct {
	rdb *redis.Client
	rdbC *redis.ClusterClient
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
	if config.Conf.Storage.RedisCluster {
		c.rdbC = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: []string{redisHost},
		})
		_, err := c.rdbC.Ping(c.ctx).Result()
		return err
	}
	c.rdb = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: config.Conf.Storage.RedisPassrod,
		DB:       config.Conf.Storage.RedisDB,
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
	if config.Conf.Storage.RedisCluster {
		val, err := c.rdbC.Get(c.ctx, key).Result()
		if err != nil {
			return "", err
		}
		return val, nil
	}
	val, err := c.rdb.Get(c.ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// Set ...
func (c *Cache) Set(key string, value interface{}, durations time.Duration) error {
	if config.Conf.Storage.RedisCluster {
		err := c.rdbC.Set(c.ctx, key, value, durations).Err()
		return err
	}
	err := c.rdb.Set(c.ctx, key, value, durations).Err()
	return err
}
