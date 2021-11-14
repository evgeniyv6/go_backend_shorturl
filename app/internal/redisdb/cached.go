package redisdb

import (
	"context"
	"time"

	"github.com/evgeniyv6/go_backend_shorturl/app/internal/hasher"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type CachedRedis struct {
	mainDB redisConn
	cache  *cache.Cache
}

func NewCachedRedis(action redisConn) (DBAction, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       1,
	})

	rCache := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &CachedRedis{
		mainDB: action,
		cache:  rCache,
	}, nil
}

func (c *CachedRedis) Save(ctx context.Context, uri string) (string, error) {
	return c.mainDB.Save(ctx, uri)
}

func (c *CachedRedis) Close() error {
	return c.mainDB.Close()
}

func (c *CachedRedis) GetLink(ctx context.Context, hash string) (string, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, c.mainDB.tracer, "cache get info")
	defer span.Finish()

	conn := c.mainDB.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			c.mainDB.logger.Errorw("Couldnot close redis connection.", "err", err)
			span.LogFields(
				log.Error(err),
			)
		}
	}()

	clearedRandNum, err := hasher.GenClear(hash)
	if err != nil {
		return "", err
	}

	span.LogFields(
		log.Uint64("get info hash", clearedRandNum),
	)

	var link string

	err = c.cache.Get(ctx, hash, link)

	switch err {
	case nil:
		c.mainDB.logger.Info("Get info from cache")
		return link, nil

	case cache.ErrCacheMiss:
		dbLink, dbErr := c.mainDB.GetLink(ctx, hash)
		if dbErr != nil {
			return "", dbErr
		}

		err = c.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   hash,
			Value: dbLink,
			TTL:   5 * time.Second,
		})

		if err != nil {
			return "", err
		}
		return dbLink, nil

	}
	return "", err
}

func (c *CachedRedis) GetInfo(ctx context.Context, hash string) (*DBRecord, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, c.mainDB.tracer, "cache get info")
	defer span.Finish()

	var shortLink DBRecord
	conn := c.mainDB.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			c.mainDB.logger.Errorw("Couldnot close redis connection.", "err", err)
			span.LogFields(
				log.Error(err),
			)
		}
	}()

	err := c.cache.Get(ctx, hash, &shortLink)

	switch err {
	case nil:
		c.mainDB.logger.Info("Get info from cache")
		return &shortLink, nil

	case cache.ErrCacheMiss:
		rec, dbErr := c.mainDB.GetInfo(ctx, hash)
		if dbErr != nil {
			return nil, dbErr
		}

		err = c.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   hash,
			Value: rec,
			TTL:   5 * time.Second,
		})

		if err != nil {
			return nil, err
		}
		return rec, nil

	}
	return nil, err
}
