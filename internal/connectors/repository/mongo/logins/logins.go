package logins

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"prp.com/sparkly/internal/connectors/services/clock"
	"prp.com/sparkly/internal/pkg"
)

const collectionUserLogins = "user_logins"

type Repository interface {
	Insert(ctx context.Context, activity pkg.LoginActivity) error
	GetActiveUsersCountByDuration(ctx context.Context, duration time.Duration) (int, error)
	GetActiveUsersByDuration(ctx context.Context, duration time.Duration) ([]bson.M, error)
}

type repository struct {
	collection *mongo.Collection
	clock      clock.Service
}

func NewRepository(mongoDB *mongo.Database, clock clock.Service) Repository {
	collection := mongoDB.Collection(collectionUserLogins)

	// Creating indexes in MongoDB
	collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{"timestamp", 1}, {"userId", 1}},
	})

	return &repository{
		collection: collection,
		clock:      clock,
	}
}

func (repository *repository) Insert(ctx context.Context, activity pkg.LoginActivity) error {
	activity.Timestamp = repository.clock.Now()

	_, err := repository.collection.InsertOne(ctx, activity)

	return err
}

func (repository *repository) GetActiveUsersByDuration(ctx context.Context, duration time.Duration) ([]bson.M, error) {
	filter := bson.M{"timestamp": bson.M{"$gte": repository.clock.Now().Add(-duration)}}

	pipeline := mongo.Pipeline{
		{{"$match", filter}},
		{{"$group", bson.M{
			"_id":             "$userId",
			"latestTimestamp": bson.M{"$max": "$timestamp"},
		}}},
	}

	cursor, err := repository.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (repository *repository) GetActiveUsersCountByDuration(ctx context.Context, duration time.Duration) (int, error) {

	users, err := repository.GetActiveUsersByDuration(ctx, duration)

	return len(users), err
}
