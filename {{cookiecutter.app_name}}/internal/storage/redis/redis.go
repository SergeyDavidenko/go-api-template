package redis

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/internal/config"
	log "github.com/sirupsen/logrus"
)

// Cache is interface structure
type Cache struct {
	rdb  *redis.Client
	rdbC *redis.ClusterClient
	ctx  context.Context
}

// New func implements the cache interface
func New() *Cache {
	return &Cache{}
}

// Init client Cache.
func (c *Cache) Init() error {
	log.Debug("redis host is - ", config.Conf.Storage.RedisHost)
	c.ctx = context.Background()
	if config.Conf.Storage.RedisCluster {
		c.rdbC = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: strings.Split(config.Conf.Storage.RedisHost, ","),
		})
		_, err := c.rdbC.Ping(c.ctx).Result()
		log.Info("Start on cluster redis mode")
		return err
	}
	log.Info("Start single redis mode")
	c.rdb = redis.NewClient(&redis.Options{
		Addr:     config.Conf.Storage.RedisHost,
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
