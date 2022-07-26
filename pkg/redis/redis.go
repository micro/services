package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/micro/micro/v3/service/config"
	log "github.com/micro/micro/v3/service/logger"
)

type Counter struct {
	sync.RWMutex

	prefix string
	client *redis.Client
}

type Config struct {
	Address  string `json:"address"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Secure   bool   `json:"secure,omitempty"`
}

const Nil = redis.Nil

func Key(args ...string) string {
	return strings.Join(args, ":")
}

func (c *Counter) Incr(ctx context.Context, key, path string, delta int64) (int64, error) {
	return c.client.IncrBy(ctx, fmt.Sprintf("%s:%s:%s", c.prefix, key, path), delta).Result()
}

func (c *Counter) Decr(ctx context.Context, key, path string, delta int64) (int64, error) {
	return c.client.DecrBy(ctx, fmt.Sprintf("%s:%s:%s", c.prefix, key, path), delta).Result()
}

func (c *Counter) Read(ctx context.Context, key, path string) (int64, error) {
	ret, err := c.client.Get(ctx, fmt.Sprintf("%s:%s:%s", c.prefix, key, path)).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return ret, err
}

func (c *Counter) Reset(ctx context.Context, key, path string) error {
	return c.client.Set(ctx, fmt.Sprintf("%s:%s:%s", c.prefix, key, path), 0, 0).Err()
}

func (c *Counter) Delete(ctx context.Context, key string) error {
	keys, err := c.client.Keys(ctx, fmt.Sprintf("%s:%s:*", c.prefix, key)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	if err := c.client.Del(ctx, keys...).Err(); err != nil && err != redis.Nil {
		return err
	}

	return nil
}

func NewCounter(prefix string) *Counter {
	var redisConfig Config

	val, err := config.Get("micro.redis")
	if err != nil {
		log.Fatalf("No redis config found %s", err)
	}

	if err := val.Scan(&redisConfig); err != nil {
		log.Fatalf("Error parsing redis config %s", err)
	}

	log.Infof("Using %v", redisConfig)

	opts := &redis.Options{
		Addr:     redisConfig.Address,
		Username: redisConfig.Username,
		Password: redisConfig.Password,
	}
	if redisConfig.Secure {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
		}
	}

	rc := redis.NewClient(opts)

	_, err = rc.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to ping redis: %v", err)
	}

	return &Counter{
		prefix: Key(prefix, "counter"),
		client: rc,
	}
}
