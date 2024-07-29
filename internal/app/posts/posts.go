package posts

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"prp.com/sparkly/internal/connectors/repository/mongo/posts"
	"prp.com/sparkly/internal/connectors/services/cache"
	"prp.com/sparkly/internal/connectors/services/kafka"
	"prp.com/sparkly/internal/pkg"
)

type Service interface {
	BackfillPolularPosts(ctx context.Context)
	PushToQueue(ctx context.Context, activity pkg.PostInteraction)
	Log(ctx context.Context, activity pkg.PostInteraction) error
	GetPopularPostsByDuration(ctx context.Context, key string, limit int) ([]pkg.PopularPost, error)
	GetPopularPosts(ctx context.Context, timeFrames map[string]time.Duration, limit int) (map[string][]pkg.PopularPost, error)
}

type service struct {
	postsRepository posts.Repository
	producerService kafka.ProducerService
	cacheService    cache.CacheService
}

func NewService(
	postsRepository posts.Repository,
	producerService kafka.ProducerService,
	cacheService cache.CacheService,
) Service {
	return &service{
		postsRepository: postsRepository,
		producerService: producerService,
		cacheService:    cacheService,
	}
}

func (service *service) PushToQueue(ctx context.Context, activity pkg.PostInteraction) {
	service.producerService.PushMessage(ctx, pkg.TopicPostsLogins, activity)
}

func (service *service) Log(ctx context.Context, activity pkg.PostInteraction) error {

	err := service.postsRepository.Insert(ctx, activity)
	if err != nil {
		return err
	}

	service.incrementPostsCache(ctx, activity.PostID, 1)

	return nil
}

func (service *service) GetPopularPostsByDuration(ctx context.Context, key string, limit int) ([]pkg.PopularPost, error) {

	if duration, exists := pkg.AnalyticsTimeFrames[key]; exists {

		results, err := service.getPostsCache(ctx, key, limit)
		if err == nil {
			log.Println("active users from cache")
			return results, err
		}

		return service.getPopularPostsFromDB(ctx, duration, limit)

	}

	return nil, nil
}

func (service *service) getPopularPostsFromDB(ctx context.Context, duration time.Duration, limit int) ([]pkg.PopularPost, error) {

	posts, err := service.postsRepository.GetPopularPosts(ctx, duration, limit)
	if err != nil {
		return nil, err
	}

	results := make([]pkg.PopularPost, 0)
	for _, item := range posts {
		id := item["_id"].(string)

		count := item["count"]

		var countInt64 int64
		switch v := count.(type) {
		case int:
			countInt64 = int64(v)
		case int32:
			countInt64 = int64(v)
		case int64:
			countInt64 = v
		default:
			fmt.Printf("unsupported type for count: %T\n", v)
			continue
		}

		results = append(results, pkg.PopularPost{
			PostID: id,
			Count:  countInt64,
		})
	}

	return results, nil
}

func (service *service) GetPopularPosts(ctx context.Context, timeFrames map[string]time.Duration, limit int) (map[string][]pkg.PopularPost, error) {

	mutex := sync.Mutex{}
	results := make(map[string][]pkg.PopularPost)
	errChannel := make(chan error, len(timeFrames))

	for key := range timeFrames {

		go func(key string) {

			posts, err := service.GetPopularPostsByDuration(ctx, key, limit)
			if err == nil {
				mutex.Lock()
				results[key] = posts
				mutex.Unlock()
			}

			errChannel <- err

		}(key)

	}

	for i := 0; i < len(timeFrames); i++ {

		err := <-errChannel
		if err != nil {
			return nil, err
		}

	}

	return results, nil

}
