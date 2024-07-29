package logins

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"prp.com/sparkly/internal/pkg"
	loginsMock "prp.com/sparkly/mocks/connectors/repository/mongo/logins"
	cacheMock "prp.com/sparkly/mocks/connectors/services/cache"
	clockMock "prp.com/sparkly/mocks/connectors/services/clock"
	kafkaMock "prp.com/sparkly/mocks/connectors/services/kafka"
)

type testServiceDependencies struct {
	cacheService     *cacheMock.CacheService
	producerService  *kafkaMock.ProducerService
	loginsRepository *loginsMock.Repository
	clock            *clockMock.Service
}

// Define your test suite
type serviceTestSuite struct {
	suite.Suite
	testServiceDependencies testServiceDependencies
	testService             Service
}

func (suite *serviceTestSuite) SetupTest() {
	suite.testServiceDependencies = testServiceDependencies{
		cacheService:     cacheMock.NewCacheService(suite.T()),
		producerService:  kafkaMock.NewProducerService(suite.T()),
		loginsRepository: loginsMock.NewRepository(suite.T()),
		clock:            clockMock.NewService(suite.T()),
	}

	suite.testService = NewService(
		suite.testServiceDependencies.loginsRepository,
		suite.testServiceDependencies.producerService,
		suite.testServiceDependencies.cacheService,
		suite.testServiceDependencies.clock,
	)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (suite *serviceTestSuite) TestPushToQueue() {
	ctx := context.Background()

	activity := pkg.LoginActivity{
		UserID: "user-1",
		Action: "login",
	}

	suite.testServiceDependencies.producerService.
		EXPECT().PushMessage(ctx, pkg.TopicLogsLogins, activity)

	suite.testService.PushToQueue(ctx, activity)

}

func (suite *serviceTestSuite) TestLog_Successful() {
	ctx := context.Background()

	activity := pkg.LoginActivity{
		UserID: "user-1",
		Action: "login",
	}

	now, _ := time.Parse("2006-01-02 15:04:05", "2024-07-28 00:00:00")

	suite.testServiceDependencies.loginsRepository.
		EXPECT().Insert(ctx, activity).Return(nil)

	suite.testServiceDependencies.clock.
		EXPECT().Now().Return(now)

	suite.testServiceDependencies.cacheService.
		EXPECT().ZAdd(
		ctx, pkg.CacheKeyActiveUsers, float64(now.Unix()), activity.UserID, now.Unix()-25*3600,
	)

	err := suite.testService.Log(ctx, activity)

	// Assertions
	assert.Nil(suite.T(), err)
}

func (suite *serviceTestSuite) TestLog_ErrorFromDB() {
	ctx := context.Background()

	activity := pkg.LoginActivity{
		UserID: "user-1",
		Action: "login",
	}

	suite.testServiceDependencies.loginsRepository.
		EXPECT().Insert(ctx, activity).Return(assert.AnError)

	err := suite.testService.Log(ctx, activity)

	// Assertions
	assert.NotNil(suite.T(), err)
}

func (suite *serviceTestSuite) TestGetActiveUsersByDuration_SuccessfulFromCache() {

	ctx := context.Background()

	var count int64 = 5
	durationKey := "last_minute"

	now, _ := time.Parse("2006-01-02 15:04:05", "2024-07-28 00:00:00")
	suite.testServiceDependencies.clock.
		EXPECT().Now().Return(now)

	cutoff := now.Add(-pkg.AnalyticsTimeFrames[durationKey]).Unix()

	suite.testServiceDependencies.cacheService.
		EXPECT().ZCount(ctx, pkg.CacheKeyActiveUsers, cutoff).Return(count, nil)

	results, err := suite.testService.GetActiveUsersByDuration(ctx, durationKey)

	// Assertions
	assert.Equal(suite.T(), int(count), results)
	assert.Nil(suite.T(), err)
}

func (suite *serviceTestSuite) TestGetActiveUsersByDuration_SuccessfulFromDB() {

	ctx := context.Background()

	var count int64 = 5
	durationKey := "last_minute"

	now, _ := time.Parse("2006-01-02 15:04:05", "2024-07-28 00:00:00")
	suite.testServiceDependencies.clock.
		EXPECT().Now().Return(now)

	cutoff := now.Add(-pkg.AnalyticsTimeFrames[durationKey]).Unix()

	suite.testServiceDependencies.cacheService.
		EXPECT().ZCount(ctx, pkg.CacheKeyActiveUsers, cutoff).Return(count, assert.AnError)

	suite.testServiceDependencies.loginsRepository.
		EXPECT().GetActiveUsersCountByDuration(ctx, pkg.AnalyticsTimeFrames[durationKey]).Return(int(count), nil)

	results, err := suite.testService.GetActiveUsersByDuration(ctx, durationKey)

	// Assertions
	assert.Equal(suite.T(), int(count), results)
	assert.Nil(suite.T(), err)
}

func (suite *serviceTestSuite) TestGetActiveUsersByDuration_ErrorFromDB() {

	ctx := context.Background()

	var count int64 = 5
	durationKey := "last_minute"

	now, _ := time.Parse("2006-01-02 15:04:05", "2024-07-28 00:00:00")
	suite.testServiceDependencies.clock.
		EXPECT().Now().Return(now)

	cutoff := now.Add(-pkg.AnalyticsTimeFrames[durationKey]).Unix()

	suite.testServiceDependencies.cacheService.
		EXPECT().ZCount(ctx, pkg.CacheKeyActiveUsers, cutoff).Return(count, assert.AnError)

	suite.testServiceDependencies.loginsRepository.
		EXPECT().GetActiveUsersCountByDuration(ctx, pkg.AnalyticsTimeFrames[durationKey]).Return(0, assert.AnError)

	results, err := suite.testService.GetActiveUsersByDuration(ctx, durationKey)

	// Assertions
	assert.Equal(suite.T(), 0, results)
	assert.NotNil(suite.T(), err)
}

func (suite *serviceTestSuite) TestGetActiveUsers_Successful() {

	ctx := context.Background()

	var count int64 = 5
	durationKey := "last_minute"

	timeFrames := map[string]time.Duration{
		"last_minute": time.Hour,
	}

	now, _ := time.Parse("2006-01-02 15:04:05", "2024-07-28 00:00:00")
	suite.testServiceDependencies.clock.
		EXPECT().Now().Return(now)

	cutoff := now.Add(-pkg.AnalyticsTimeFrames[durationKey]).Unix()

	suite.testServiceDependencies.cacheService.
		EXPECT().ZCount(ctx, pkg.CacheKeyActiveUsers, cutoff).Return(count, nil)

	results, err := suite.testService.GetActiveUsers(ctx, timeFrames)

	expectedResults := map[string]int{
		durationKey: int(count),
	}

	// Assertions
	assert.Equal(suite.T(), expectedResults, results)
	assert.Nil(suite.T(), err)

}

func (suite *serviceTestSuite) TestGetActiveUsers_Error() {

	ctx := context.Background()

	durationKey := "last_minute"

	timeFrames := map[string]time.Duration{
		durationKey: time.Hour,
	}

	now, _ := time.Parse("2006-01-02 15:04:05", "2024-07-28 00:00:00")
	suite.testServiceDependencies.clock.
		EXPECT().Now().Return(now)

	cutoff := now.Add(-pkg.AnalyticsTimeFrames[durationKey]).Unix()

	suite.testServiceDependencies.cacheService.
		EXPECT().ZCount(ctx, pkg.CacheKeyActiveUsers, cutoff).Return(0, assert.AnError)

	suite.testServiceDependencies.loginsRepository.
		EXPECT().GetActiveUsersCountByDuration(ctx, pkg.AnalyticsTimeFrames[durationKey]).Return(0, assert.AnError)

	results, err := suite.testService.GetActiveUsers(ctx, timeFrames)

	// Assertions
	assert.Nil(suite.T(), results)
	assert.NotNil(suite.T(), err)
}
