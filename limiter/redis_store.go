package limiter

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type Store interface {
	Allow(key string, limit int, duration time.Duration) (bool, error)
}

type RedisStore struct {
	rdb redis.Client
}

func NewRedisStore(rdb redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func (s *RedisStore) Allow(key string, limit int, duration time.Duration) (bool, error) {
	ctx := context.Background()

	countStr, err := s.rdb.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}

	var count int
	if err == redis.Nil {
		count = 0
	} else {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			return false, err
		}
	}

	if count > limit {
		s.rdb.Expire(ctx, key, duration)
		return false, nil
	}

	countIncr, err := s.rdb.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if countIncr == 1 {
		s.rdb.Expire(ctx, key, time.Second)
	}

	if countIncr > int64(limit) {
		s.rdb.Expire(ctx, key, duration)
		return false, nil
	}

	return true, nil
}
