package logins

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"prp.com/sparkly/internal/pkg"
	clockMock "prp.com/sparkly/mocks/connectors/services/clock"
)

type mockClock struct {
	currentTime time.Time
}

func (m *mockClock) Now() time.Time {
	return m.currentTime
}

func TestRepository_Insert(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	clock := clockMock.NewService(t)

	now, _ := time.Parse("2006-01-02 15:04:05", "2024-07-28 00:00:00")

	mt.Run("Insert login activity", func(mt *mtest.T) {

		repo := NewRepository(mt.DB, clock)

		clock.EXPECT().Now().Return(now)

		activity := pkg.LoginActivity{
			UserID: "user1",
			Action: "login",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.Insert(context.Background(), activity)
		assert.NoError(t, err)

		mt.AddMockResponses(mtest.CreateCursorResponse(1, mt.DB.Name()+"."+collectionUserLogins, mtest.FirstBatch, bson.D{
			{"userId", activity.UserID},
			{"action", activity.Action},
			{"timestamp", now},
		}))

		var result pkg.LoginActivity
		err = mt.DB.Collection(collectionUserLogins).FindOne(context.Background(), bson.M{"userId": activity.UserID}).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, activity.UserID, result.UserID)
		assert.Equal(t, activity.Action, result.Action)
		assert.Equal(t, now, result.Timestamp)
	})
}
