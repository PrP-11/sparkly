package posts

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"prp.com/sparkly/internal/pkg"
)

func (suite *serviceTestSuite) TestBackfillPolularPosts_Successful() {
	ctx := context.Background()

	responseFromDB := []bson.M{
		{"_id": "post-2", "count": 9999},
		{"_id": "post-3", "count": 8008},
	}

	for key, duration := range pkg.AnalyticsTimeFrames {

		cacheKey := fmt.Sprintf(pkg.CachekeyPopularPosts, key)

		suite.testServiceDependencies.cacheService.
			EXPECT().DeleteKeys(ctx, cacheKey).Return(nil)

		suite.testServiceDependencies.postsRepository.
			EXPECT().GetPopularPosts(ctx, duration, 10).Return(responseFromDB, nil)

		suite.testServiceDependencies.cacheService.
			EXPECT().ZIncrBy(ctx, cacheKey, "post-2", float64(9999)).Return()

		suite.testServiceDependencies.cacheService.
			EXPECT().Expire(ctx, cacheKey, 2*duration).Return()

		suite.testServiceDependencies.cacheService.
			EXPECT().ZIncrBy(ctx, cacheKey, "post-3", float64(8008)).Return()

		suite.testServiceDependencies.cacheService.
			EXPECT().Expire(ctx, cacheKey, 2*duration).Return()
	}

	suite.testService.BackfillPolularPosts(ctx)

}
