package storage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr, password string) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{Client: rdb}, nil
}

func (r *RedisClient) SetPrice(ctx context.Context, symbol string, price float64) error {
	key := fmt.Sprintf("price:%s", symbol)
	return r.Client.Set(ctx, key, price, 0).Err()
}

func (r *RedisClient) GetPrice(ctx context.Context, symbol string) (float64, error) {
	key := fmt.Sprintf("price:%s", symbol)
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(val, 64)
}

func (r *RedisClient) AddPriceHistory(ctx context.Context, symbol string, price float64, timestamp int64) error {
	key := fmt.Sprintf("history:%s", symbol)
	return r.Client.ZAdd(ctx, key, redis.Z{
		Score:  float64(timestamp),
		Member: price,
	}).Err()
}

func (r *RedisClient) GetOldestPrice(ctx context.Context, symbol string) (float64, error) {
	key := fmt.Sprintf("history:%s", symbol)
	vals, err := r.Client.ZRangeWithScores(ctx, key, 0, 0).Result()
	if err != nil {
		return 0, err
	}
	if len(vals) == 0 {
		return 0, fmt.Errorf("no history for %s", symbol)
	}
	strPrice := fmt.Sprintf("%v", vals[0].Member)
	return strconv.ParseFloat(strPrice, 64)
}

func (r *RedisClient) TrimHistory(ctx context.Context, symbol string, retention int64) error {
	key := fmt.Sprintf("history:%s", symbol)
	min := "-inf"
	max := fmt.Sprintf("%d", retention)
	return r.Client.ZRemRangeByScore(ctx, key, min, max).Err()
}

func (r *RedisClient) CachePnL(ctx context.Context, symbol string, pnl float64, ttl time.Duration) error {
	key := fmt.Sprintf("pnl:%s", symbol)
	return r.Client.Set(ctx, key, pnl, ttl).Err()
}

func (r *RedisClient) GetCachedPnL(ctx context.Context, symbol string) (float64, error) {
	key := fmt.Sprintf("pnl:%s", symbol)
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(val, 64)
}
