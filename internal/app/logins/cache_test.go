package logins

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"prp.com/sparkly/internal/pkg"
)

func (suite *serviceTestSuite) TestBackfillActiveUsers_Successful() {
	ctx := context.Background()

	now, _ := time.Parse("2006-01-02 15:04:05", "2024-07-28 00:00:00")

	responseFromDB := []bson.M{
		{"_id": "user-1", "latestTimestamp": primitive.NewDateTimeFromTime(now)},
	}

	for _, duration := range pkg.AnalyticsTimeFrames {

		suite.testServiceDependencies.loginsRepository.
			EXPECT().GetActiveUsersByDuration(ctx, duration).Return(responseFromDB, nil)

		suite.testServiceDependencies.cacheService.
			EXPECT().DeleteKeys(ctx, pkg.CacheKeyActiveUsers).Return(nil)

		userID := responseFromDB[0]["_id"].(string)
		timestamp := responseFromDB[0]["latestTimestamp"].(primitive.DateTime).Time().Unix()

		suite.testServiceDependencies.cacheService.
			EXPECT().ZAdd(ctx, pkg.CacheKeyActiveUsers, float64(timestamp), userID, now.Unix()-25*3600).Return()

	}

	suite.testService.BackfillActiveUsers(ctx)

}
