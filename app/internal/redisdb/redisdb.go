package redisdb

import (
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/evgeniyv6/go_backend_shorturl/app/internal/hasher"

	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

type (
	redisConn struct{ pool *redis.Pool }
	DBRecord  struct {
		ID   uint64 `json:"id"`
		Link string `json:"link"`
		Stat int    `json:"stat"`
	}
	DBAction interface {
		Save(string) (string, error)
		Close() error
		GetLink(string) (string, error)
		GetInfo(string) (*DBRecord, error)
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

func NewPool(addr, port string) (DBAction, error) {
	p := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", net.JoinHostPort(addr, port), redis.DialPassword(os.Getenv("REDIS_PASSWORD")))
		},
	}
	return &redisConn{p}, nil
}

func (r *redisConn) used(num uint64) bool {
	conn := r.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			zap.S().Errorw("Couldnot close redis connection.", "err", err)
		}
	}()

	exists, err := redis.Bool(conn.Do("EXISTS", "go:shorted:"+strconv.FormatUint(num, 10)))
	if err != nil {
		return false
	}
	return exists
}

func (r *redisConn) Save(link string) (string, error) {
	conn := r.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			zap.S().Errorw("Couldnot close redis connection.", "err", err)
		}
	}()

	var randNum uint64
	for pres := true; pres; pres = r.used(randNum) {
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

func (r *redisConn) GetLink(hash string) (string, error) {
	conn := r.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			zap.S().Errorw("Couldnot close redis connection.", "err", err)
		}
	}()

	clearedRandNum, err := hasher.GenClear(hash)
	if err != nil {
		return "", err
	}

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

func (r *redisConn) GetInfo(hash string) (*DBRecord, error) {
	var shortLink DBRecord
	conn := r.pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			zap.S().Errorw("Couldnot close redis connection.", "err", err)
		}
	}()

	clearedRandNum, err := hasher.GenClear(hash)
	if err != nil {
		return nil, err
	}

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
