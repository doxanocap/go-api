package cache

import (
	"auth-api/internal/interfaces"
	"auth-api/logger"
	"context"
	"encoding/json"
	"github.com/doxanocap/pkg/errs"
	"log/slog"
	"time"
)

const (
	slogGroupKey = "cache"
)

type Cache struct {
	provider interfaces.ICacheProvider
	log      *logger.Logger
}

func NewCacheProcessor(provider interfaces.ICacheProvider, log *logger.Logger) *Cache {
	return &Cache{
		provider: provider,
		log:      log,
	}
}

func (c *Cache) Set(ctx context.Context, key string, value []byte) error {
	log := c.log.WithAttrs(slog.Group(slogGroupKey,
		slog.String("key", key),
		slog.String("value", string(value))))

	err := c.provider.Set(ctx, key, value)
	if err != nil {
		log.Error(err.Error())
		return errs.Wrap("cache.proc.Set", err)
	}

	log.Info("set")
	return nil
}

func (c *Cache) SetJSON(ctx context.Context, key string, value interface{}) error {
	raw, err := json.Marshal(value)
	log := c.log.WithAttrs(slog.Group(slogGroupKey,
		slog.String("key", key),
		slog.String("value", string(raw))))

	if err != nil {
		log.Error(err.Error())
		return errs.Wrap("marshal", err)
	}

	err = c.provider.Set(ctx, key, raw)
	if err != nil {
		log.Error(err.Error())
		return errs.Wrap("cache.proc.SetJSON", err)
	}

	log.Info("setJSON")
	return nil
}

func (c *Cache) SetJSONWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	raw, err := json.Marshal(value)
	log := c.log.WithAttrs(slog.Group(slogGroupKey,
		slog.String("key", key),
		slog.String("value", string(raw)),
		slog.Duration("ttl", ttl)))

	if err != nil {
		log.Error(err.Error())
		return errs.Wrap("marshal", err)
	}

	err = c.provider.SetWithTTL(ctx, key, raw, ttl)
	if err != nil {
		log.Error(err.Error())
		return errs.Wrap("cache.proc.SetJSONWithTTL", err)
	}

	log.Info("setJSONWithTTL")
	return nil
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	log := c.log.WithAttrs(slog.Group(slogGroupKey,
		slog.String("key", key)))

	raw, err := c.provider.Get(ctx, key)
	if err != nil {
		log.Error(err.Error())
		return nil, errs.Wrap("cache.proc.Get", err)
	}

	log.Info("get")
	return raw, nil
}

func (c *Cache) GetJSON(ctx context.Context, key string, v interface{}) error {
	log := c.log.WithAttrs(slog.Group(slogGroupKey,
		slog.String("key", key)))

	raw, err := c.provider.Get(ctx, key)
	if err != nil {
		log.Error(err.Error())
		return errs.Wrap("cache.proc.GetJSON", err)
	}

	err = json.Unmarshal(raw, v)
	if err != nil {
		log.Error(err.Error())
		return errs.Wrap("unmarshal", err)
	}

	log.Info("getJSON")
	return nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	log := c.log.WithAttrs(slog.Group(slogGroupKey,
		slog.String("key", key)))

	err := c.provider.Delete(ctx, key)
	if err != nil {
		log.Error(err.Error())
		return errs.Wrap("cache.proc.Delete", err)
	}

	log.Info("delete")
	return nil
}

func (c *Cache) FlushAll(ctx context.Context) error {
	err := c.provider.FlushAll(ctx)
	if err != nil {
		c.log.Error(err.Error())
		return errs.Wrap("cache.proc.FlushAll", err)
	}
	c.log.Info("flushAll")
	return nil
}
