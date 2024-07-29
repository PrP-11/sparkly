package posts

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"prp.com/sparkly/internal/pkg"
	postsMock "prp.com/sparkly/mocks/connectors/repository/mongo/posts"
	cacheMock "prp.com/sparkly/mocks/connectors/services/cache"
	kafkaMock "prp.com/sparkly/mocks/connectors/services/kafka"
)

type testServiceDependencies struct {
	cacheService    *cacheMock.CacheService
	producerService *kafkaMock.ProducerService
	postsRepository *postsMock.Repository
}

// Define your test suite
type serviceTestSuite struct {
	suite.Suite
	testServiceDependencies testServiceDependencies
	testService             Service
}

func (suite *serviceTestSuite) SetupTest() {
	suite.testServiceDependencies = testServiceDependencies{
		cacheService:    cacheMock.NewCacheService(suite.T()),
		producerService: kafkaMock.NewProducerService(suite.T()),
		postsRepository: postsMock.NewRepository(suite.T()),
	}

	suite.testService = NewService(
		suite.testServiceDependencies.postsRepository,
		suite.testServiceDependencies.producerService,
		suite.testServiceDependencies.cacheService,
	)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (suite *serviceTestSuite) TestPushToQueue() {
	ctx := context.Background()

	activity := pkg.PostInteraction{
		UserID: "user-1",
		PostID: "post-1",
		Action: "like",
	}

	suite.testServiceDependencies.producerService.
		EXPECT().PushMessage(ctx, pkg.TopicPostsLogins, activity)

	suite.testService.PushToQueue(ctx, activity)

}

func (suite *serviceTestSuite) TestLog_Successful() {
	ctx := context.Background()

	activity := pkg.PostInteraction{
		UserID: "user-1",
		PostID: "post-1",
		Action: "like",
	}

	suite.testServiceDependencies.postsRepository.
		EXPECT().Insert(ctx, activity).Return(nil)

	for key, duration := range pkg.AnalyticsTimeFrames {
		cacheKey := fmt.Sprintf(pkg.CachekeyPopularPosts, key)

		suite.testServiceDependencies.cacheService.
			EXPECT().ZIncrBy(ctx, cacheKey, activity.PostID, 1.00)

		suite.testServiceDependencies.cacheService.
			EXPECT().Expire(ctx, cacheKey, 2*duration)
	}

	err := suite.testService.Log(ctx, activity)

	// Assertions
	assert.Nil(suite.T(), err)
}

func (suite *serviceTestSuite) TestLog_ErrorFromDB() {
	ctx := context.Background()

	activity := pkg.PostInteraction{
		UserID: "user-1",
		PostID: "post-1",
		Action: "like",
	}

	suite.testServiceDependencies.postsRepository.
		EXPECT().Insert(ctx, activity).Return(assert.AnError)

	err := suite.testService.Log(ctx, activity)

	// Assertions
	assert.NotNil(suite.T(), err)
}

func (suite *serviceTestSuite) TestGetPopularPostsByDuration_SuccessfulFromCache() {

	ctx := context.Background()

	limit := 2
	durationKey := "last_minute"
	cacheKey := fmt.Sprintf(pkg.CachekeyPopularPosts, durationKey)

	mapFromCache := map[string]float64{
		"post-2": 9876,
		"post-3": 8008,
	}

	suite.testServiceDependencies.cacheService.
		EXPECT().ZRevRangeWithScores(ctx, cacheKey, int64(limit)).Return(mapFromCache, nil)

	expectedResponse := []pkg.PopularPost{
		{
			PostID: "post-2",
			Count:  int64(9876),
		},
		{
			PostID: "post-3",
			Count:  int64(8008),
		},
	}

	results, err := suite.testService.GetPopularPostsByDuration(ctx, durationKey, limit)

	// Assertions
	assert.Equal(suite.T(), expectedResponse, results)
	assert.Nil(suite.T(), err)
}

func (suite *serviceTestSuite) TestGetPopularPostsByDuration_SuccessfulFromDB() {

	ctx := context.Background()

	limit := 2
	durationKey := "last_minute"
	cacheKey := fmt.Sprintf(pkg.CachekeyPopularPosts, durationKey)

	responseFromDB := []bson.M{
		{"_id": "post-2", "count": 9999},
		{"_id": "post-3", "count": 8008},
	}

	suite.testServiceDependencies.cacheService.
		EXPECT().ZRevRangeWithScores(ctx, cacheKey, int64(limit)).Return(nil, assert.AnError)

	suite.testServiceDependencies.postsRepository.
		EXPECT().GetPopularPosts(ctx, pkg.AnalyticsTimeFrames[durationKey], limit).Return(responseFromDB, nil)

	results, err := suite.testService.GetPopularPostsByDuration(ctx, durationKey, limit)

	expectedResponse := []pkg.PopularPost{
		{
			PostID: "post-2",
			Count:  int64(9999),
		},
		{
			PostID: "post-3",
			Count:  int64(8008),
		},
	}

	// Assertions
	assert.Equal(suite.T(), expectedResponse, results)
	assert.Nil(suite.T(), err)
}

func (suite *serviceTestSuite) TestGetPopularPostsByDuration_ErrorFromDB() {

	ctx := context.Background()

	limit := 2
	durationKey := "last_minute"
	cacheKey := fmt.Sprintf(pkg.CachekeyPopularPosts, durationKey)

	suite.testServiceDependencies.cacheService.
		EXPECT().ZRevRangeWithScores(ctx, cacheKey, int64(limit)).Return(nil, assert.AnError)

	suite.testServiceDependencies.postsRepository.
		EXPECT().GetPopularPosts(ctx, pkg.AnalyticsTimeFrames[durationKey], limit).Return(nil, assert.AnError)

	results, err := suite.testService.GetPopularPostsByDuration(ctx, durationKey, limit)

	// Assertions
	assert.Nil(suite.T(), results)
	assert.NotNil(suite.T(), err)
}

func (suite *serviceTestSuite) TestGetPopularPostsByDuration_InvalidDurationKey() {

	ctx := context.Background()

	limit := 2
	durationKey := "last_minutes"

	results, err := suite.testService.GetPopularPostsByDuration(ctx, durationKey, limit)

	// Assertions
	assert.Nil(suite.T(), results)
	assert.Nil(suite.T(), err)
}

func (suite *serviceTestSuite) TestGetPopularPosts_Successful() {

	ctx := context.Background()

	limit := 2

	timeFrames := map[string]time.Duration{
		"last_minute": time.Hour,
	}

	mapFromCache := map[string]float64{
		"post-2": 9876,
		"post-3": 8008,
	}

	durationKey := "last_minute"
	cacheKey := fmt.Sprintf(pkg.CachekeyPopularPosts, durationKey)

	suite.testServiceDependencies.cacheService.
		EXPECT().ZRevRangeWithScores(ctx, cacheKey, int64(limit)).Return(mapFromCache, nil)

	results, err := suite.testService.GetPopularPosts(ctx, timeFrames, limit)

	expectedResults := map[string][]pkg.PopularPost{
		durationKey: {
			{
				PostID: "post-2",
				Count:  int64(9876),
			},
			{
				PostID: "post-3",
				Count:  int64(8008),
			},
		},
	}

	// Assertions
	assert.Equal(suite.T(), expectedResults, results)
	assert.Nil(suite.T(), err)
}

func (suite *serviceTestSuite) TestGetPopularPosts_Error() {

	ctx := context.Background()

	limit := 2

	timeFrames := map[string]time.Duration{
		"last_minute": time.Hour,
	}

	durationKey := "last_minute"
	cacheKey := fmt.Sprintf(pkg.CachekeyPopularPosts, durationKey)

	suite.testServiceDependencies.cacheService.
		EXPECT().ZRevRangeWithScores(ctx, cacheKey, int64(limit)).Return(nil, assert.AnError)

	suite.testServiceDependencies.postsRepository.
		EXPECT().GetPopularPosts(ctx, pkg.AnalyticsTimeFrames[durationKey], limit).Return(nil, assert.AnError)

	results, err := suite.testService.GetPopularPosts(ctx, timeFrames, limit)

	// Assertions
	assert.Nil(suite.T(), results)
	assert.NotNil(suite.T(), err)
}
