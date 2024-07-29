package cache

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheService interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) error
	DeleteKeys(ctx context.Context, keys ...string) error
	Expire(ctx context.Context, key string, duration time.Duration)
	ZAdd(ctx context.Context, key string, score float64, member interface{}, remMax int64)
	ZCount(ctx context.Context, key string, cutoff int64) (int64, error)
	ZIncrBy(ctx context.Context, key, member string, increment float64)
	ZRevRangeWithScores(ctx context.Context, key string, limit int64) (map[string]float64, error)
}

type cacheService struct {
	client *redis.Client
}

func NewCacheService() CacheService {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Printf("Connected to Redis: %s", pong)

	return &cacheService{
		client: rdb,
	}
}

func (service *cacheService) Get(ctx context.Context, key string) (string, error) {
	return service.client.Get(ctx, key).Result()
}

func (service *cacheService) Set(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	return service.client.Set(ctx, key, value, duration).Err()
}

func (service *cacheService) DeleteKeys(ctx context.Context, keys ...string) error {
	_, err := service.client.Del(ctx, keys...).Result()
	return err
}

func (service *cacheService) Expire(ctx context.Context, key string, duration time.Duration) {
	service.client.Expire(ctx, key, duration)
}

func (service *cacheService) ZAdd(ctx context.Context, key string, score float64, member interface{}, remMax int64) {
	service.client.ZAdd(ctx, key, &redis.Z{
		Score:  score,
		Member: member,
	})

	// Remove entries older than the remMax
	service.client.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(remMax, 10))
}

func (service *cacheService) ZCount(ctx context.Context, key string, cutoff int64) (int64, error) {
	return service.client.ZCount(ctx, key, strconv.FormatInt(cutoff, 10), "+inf").Result()
}

func (service *cacheService) ZIncrBy(ctx context.Context, key, member string, increment float64) {
	service.client.ZIncrBy(ctx, key, increment, member)
}

func (service *cacheService) ZRevRangeWithScores(ctx context.Context, key string, limit int64) (map[string]float64, error) {
	result := make(map[string]float64)

	posts, err := service.client.ZRevRangeWithScores(ctx, key, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		result[post.Member.(string)] = post.Score
	}

	return result, nil
}
