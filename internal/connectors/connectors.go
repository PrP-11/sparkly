package connectors

import (
	"go.mongodb.org/mongo-driver/mongo"
	"prp.com/sparkly/internal/connectors/repository/mongo/logins"
	"prp.com/sparkly/internal/connectors/repository/mongo/posts"
	"prp.com/sparkly/internal/connectors/services/cache"
	"prp.com/sparkly/internal/connectors/services/clock"
	"prp.com/sparkly/internal/connectors/services/kafka"
)

type Connector struct {
	LoginsRepository logins.Repository
	PostsRepository  posts.Repository
	ProducerService  kafka.ProducerService
	CacheService     cache.CacheService
}

func NewConnector(
	mongoDB *mongo.Database,
	clock clock.Service,
) Connector {
	return Connector{
		LoginsRepository: logins.NewRepository(mongoDB, clock),
		PostsRepository:  posts.NewRepository(mongoDB, clock),
		ProducerService:  kafka.NewProducerService(),
		CacheService:     cache.NewCacheService(),
	}
}
