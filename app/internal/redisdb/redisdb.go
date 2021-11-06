package redisdb

import (
	"context"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go/log"

	"github.com/opentracing/opentracing-go"

	"github.com/evgeniyv6/go_backend_shorturl/app/internal/logger"

	"github.com/evgeniyv6/go_backend_shorturl/app/internal/hasher"

	"github.com/gomodule/redigo/redis"
)

type (
	redisConn struct {
		logger logger.ZapWrapper
		pool   *redis.Pool
		tracer opentracing.Tracer
	}
	DBRecord struct {
		ID   uint64 `json:"id"`
		Link string `json:"link"`
		Stat int    `json:"stat"`
	}
	DBAction interface {
		Save(context.Context, string) (string, error)
		Close() error
		GetLink(context.Context, string) (string, error)
		GetInfo(context.Context, string) (*DBRecord, error)
	}
	errorString struct {
		s string
	}
)

func (e *errorString) Error() string {
	return e.s + "\n"
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewPool(addr, port string, logger logger.ZapWrapper, tracer opentracing.Tracer) (DBAction, error) {
	p := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", net.JoinHostPort(addr, port))
		},
	}
	return &redisConn{logger, p, tracer}, nil
}

func (r *redisConn) used(ctx context.Context, num uint64) bool {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, "Check redis EXISTS")
	defer span.Finish()

	span.LogFields(
		log.Uint64("redis id", num),
	)

	conn := r.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			r.logger.Errorw("Couldnot close redis connection.", "err", err)
			span.LogFields(log.Error(err))
		}
	}()

	exists, err := redis.Bool(conn.Do("EXISTS", "go:shorted:"+strconv.FormatUint(num, 10)))
	if err != nil {
		return false
	}
	return exists
}

func (r *redisConn) Save(ctx context.Context, link string) (string, error) {
	conn := r.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			r.logger.Errorw("Couldnot close redis connection.", "err", err)
		}
	}()

	var randNum uint64
	for pres := true; pres; pres = r.used(ctx, randNum) {
		randNum = rand.Uint64()
	}

	dbRec := DBRecord{randNum, link, 0}
	_, err := conn.Do("HMSET", redis.Args{"go:shorted:" + strconv.FormatUint(randNum, 10)}.AddFlat(dbRec)...)
	if err != nil {
		return "", err
	}
	return string(hasher.GenHash(randNum)), nil
}

func (r *redisConn) Close() error {
	return r.pool.Close()
}

func (r *redisConn) GetLink(ctx context.Context, hash string) (string, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, "redis get info")
	defer span.Finish()

	conn := r.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			r.logger.Errorw("Couldnot close redis connection.", "err", err)
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
		log.Uint64("get link hash", clearedRandNum),
	)

	dbRecLink, err := redis.String(conn.Do("HGET", "go:shorted:"+strconv.FormatUint(clearedRandNum, 10), "Link"))
	if err != nil {
		return "", err
	} else if len(dbRecLink) == 0 {
		return "", &errorString{"Short link not found in DB."}
	}

	_, err = conn.Do("HINCRBY", "go:shorted:"+strconv.FormatUint(clearedRandNum, 10), "Stat", 1)
	if err != nil {
		return "", err
	}
	return dbRecLink, nil
}

func (r *redisConn) GetInfo(ctx context.Context, hash string) (*DBRecord, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, "redis get info")
	defer span.Finish()

	var shortLink DBRecord
	conn := r.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			r.logger.Errorw("Couldnot close redis connection.", "err", err)
			span.LogFields(
				log.Error(err),
			)
		}
	}()

	clearedRandNum, err := hasher.GenClear(hash)
	if err != nil {
		return nil, err
	}

	span.LogFields(
		log.Uint64("get info hash", clearedRandNum),
	)

	val, err := redis.Values(conn.Do("HGETALL", "go:shorted:"+strconv.FormatUint(clearedRandNum, 10)))
	if err != nil {
		return nil, err
	} else if len(val) == 0 {
		return nil, &errorString{"Short link not found in DB."}
	}
	err = redis.ScanStruct(val, &shortLink)
	if err != nil {
		return nil, err
	}

	return &shortLink, nil
}
