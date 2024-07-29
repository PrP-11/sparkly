package app

import (
	"prp.com/sparkly/internal/app/logins"
	"prp.com/sparkly/internal/app/posts"
	"prp.com/sparkly/internal/connectors"
	"prp.com/sparkly/internal/connectors/services/clock"
)

type Services struct {
	LoginsService logins.Service
	PostsService  posts.Service
}

func NewServices(
	connector connectors.Connector,
	clock clock.Service,
) Services {

	return Services{

		LoginsService: logins.NewService(
			connector.LoginsRepository,
			connector.ProducerService,
			connector.CacheService,
			clock,
		),

		PostsService: posts.NewService(
			connector.PostsRepository,
			connector.ProducerService,
			connector.CacheService,
		),
	}

}
