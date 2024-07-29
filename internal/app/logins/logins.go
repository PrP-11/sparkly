package logins

import (
	"context"
	"log"
	"sync"
	"time"

	"prp.com/sparkly/internal/connectors/repository/mongo/logins"
	"prp.com/sparkly/internal/connectors/services/cache"
	"prp.com/sparkly/internal/connectors/services/clock"
	"prp.com/sparkly/internal/connectors/services/kafka"
	"prp.com/sparkly/internal/pkg"
)

type Service interface {
	BackfillActiveUsers(ctx context.Context)
	PushToQueue(ctx context.Context, activity pkg.LoginActivity)
	Log(ctx context.Context, activity pkg.LoginActivity) error
	GetActiveUsersByDuration(ctx context.Context, key string) (int, error)
	GetActiveUsers(ctx context.Context, timeFrames map[string]time.Duration) (map[string]int, error)
}

type service struct {
	loginsRepository logins.Repository
	producerService  kafka.ProducerService
	cacheService     cache.CacheService
	clock            clock.Service
}

func NewService(
	loginsRepository logins.Repository,
	producerService kafka.ProducerService,
	cacheService cache.CacheService,
	clock clock.Service,
) Service {
	return &service{
		loginsRepository: loginsRepository,
		producerService:  producerService,
		cacheService:     cacheService,
		clock:            clock,
	}
}

func (service *service) PushToQueue(ctx context.Context, activity pkg.LoginActivity) {
	service.producerService.PushMessage(ctx, pkg.TopicLogsLogins, activity)
}

func (service *service) Log(ctx context.Context, activity pkg.LoginActivity) error {

	err := service.loginsRepository.Insert(ctx, activity)
	if err != nil {
		return err
	}

	service.addUserToCache(ctx, activity.UserID)

	return nil
}

func (service *service) GetActiveUsersByDuration(ctx context.Context, key string) (int, error) {

	if duration, exists := pkg.AnalyticsTimeFrames[key]; exists {

		count, err := service.getActiveUsersFromCache(ctx, duration)
		if err == nil {

			log.Println("active users from cache")
			return int(count), nil
		}

		return service.loginsRepository.GetActiveUsersCountByDuration(ctx, duration)
	}

	return 0, nil
}

func (service *service) GetActiveUsers(ctx context.Context, timeFrames map[string]time.Duration) (map[string]int, error) {

	mutex := sync.Mutex{}
	results := make(map[string]int)
	errChannel := make(chan error, len(timeFrames))

	for key := range timeFrames {

		go func(key string) {

			count, err := service.GetActiveUsersByDuration(ctx, key)
			if err == nil {
				mutex.Lock()
				results[key] = count
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
