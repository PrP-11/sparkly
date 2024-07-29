package logins

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"prp.com/sparkly/internal/pkg"
)

func (service *service) addUserToCache(ctx context.Context, userId string) {
	now := service.clock.Now().Unix()
	service.addUserToCacheWithTime(ctx, userId, now)
}

func (service *service) addUserToCacheWithTime(ctx context.Context, userId string, now int64) {
	service.cacheService.ZAdd(ctx, pkg.CacheKeyActiveUsers, float64(now), userId, now-25*3600) // remove entry older than 25 hours
}

func (service *service) getActiveUsersFromCache(ctx context.Context, duration time.Duration) (int64, error) {
	cutoff := service.clock.Now().Add(-duration).Unix()
	return service.cacheService.ZCount(ctx, pkg.CacheKeyActiveUsers, cutoff)
}

func (service *service) BackfillActiveUsers(ctx context.Context) {

	for _, duration := range pkg.AnalyticsTimeFrames {

		users, err := service.loginsRepository.GetActiveUsersByDuration(ctx, duration)
		if err != nil {
			continue
		}

		service.cacheService.DeleteKeys(ctx, pkg.CacheKeyActiveUsers)

		for _, user := range users {

			userID := user["_id"].(string)
			timestamp := user["latestTimestamp"].(primitive.DateTime).Time().Unix()

			service.addUserToCacheWithTime(ctx, userID, timestamp)

		}
	}

}
