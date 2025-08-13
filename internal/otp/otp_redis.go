package otp

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	rdb       *redis.Client
	expiresIn time.Duration
	maxReq    int
	window    time.Duration
}

type RedisConfig struct {
	RDB         *redis.Client
	ExpiresIn   time.Duration
	MaxRequests int
	Window      time.Duration
}

func NewRedis(c RedisConfig) *RedisService {
	return &RedisService{ rdb: c.RDB, expiresIn: c.ExpiresIn, maxReq: c.MaxRequests, window: c.Window }
}

func (s *RedisService) Request(ctx context.Context, phone string) (string, time.Time, error) {
	rlKey := fmt.Sprintf("otp:rl:%s", phone)
	cnt, err := s.rdb.Incr(ctx, rlKey).Result()
	if err != nil { return "", time.Time{}, err }
	if cnt == 1 { s.rdb.Expire(ctx, rlKey, s.window) }
	if cnt > int64(s.maxReq) {
		return "", time.Time{}, fmt.Errorf("rate limit exceeded: max %d per %s", s.maxReq, s.window)
	}

	code := gen6()
	key := fmt.Sprintf("otp:code:%s", phone)
	if err := s.rdb.Set(ctx, key, code, s.expiresIn).Err(); err != nil {
		return "", time.Time{}, err
	}
	return code, time.Now().Add(s.expiresIn), nil
}

func (s *RedisService) Verify(ctx context.Context, phone, code string) (bool, error) {
	key := fmt.Sprintf("otp:code:%s", phone)
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil { return false, nil }
	if err != nil { return false, err }
	if val != code { return false, nil }
	_ = s.rdb.Del(ctx, key).Err() // one-time
	return true, nil
}

func gen6() string {
	var b [3]byte
	_, _ = rand.Read(b[:])
	n := (int(b[0])<<16 | int(b[1])<<8 | int(b[2])) % 1000000
	return fmt.Sprintf("%06d", n)
}
