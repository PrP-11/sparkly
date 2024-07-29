package posts

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"prp.com/sparkly/internal/connectors/services/clock"
	"prp.com/sparkly/internal/pkg"
)

const collectionPostInteractions = "post_interactions"

type Repository interface {
	Insert(ctx context.Context, activity pkg.PostInteraction) error
	GetPopularPosts(ctx context.Context, duration time.Duration, limit int) ([]bson.M, error)
}

type repository struct {
	collection *mongo.Collection
	clock      clock.Service
}

func NewRepository(mongoDB *mongo.Database, clock clock.Service) Repository {
	collection := mongoDB.Collection(collectionPostInteractions)

	collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{"timestamp", 1}, {"postId", 1}},
	})

	return &repository{
		collection: collection,
		clock:      clock,
	}
}

func (repository *repository) Insert(ctx context.Context, activity pkg.PostInteraction) error {
	activity.Timestamp = repository.clock.Now()

	_, err := repository.collection.InsertOne(ctx, activity)

	return err
}

func (repository *repository) GetPopularPosts(ctx context.Context, duration time.Duration, limit int) ([]bson.M, error) {
	filter := bson.M{"timestamp": bson.M{"$gte": repository.clock.Now().Add(-duration)}}

	matchStage := bson.D{{"$match", filter}}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$postId"},
			{"count", bson.D{{"$sum", 1}}},
		}},
	}

	sortStage := bson.D{{"$sort", bson.D{{"count", -1}}}}
	limitStage := bson.D{{"$limit", limit}}

	cursor, err := repository.collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, sortStage, limitStage})
	if err != nil {
		return nil, err
	}

	var posts []bson.M
	err = cursor.All(ctx, &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
