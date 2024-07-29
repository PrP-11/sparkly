package consumers

import (
	"prp.com/sparkly/internal/app"
	"prp.com/sparkly/internal/app/logins"
	"prp.com/sparkly/internal/app/posts"
	"prp.com/sparkly/internal/pkg"
)

type Consumer interface {
	Register() pkg.TopicConsumerMap
}

type consumer struct {
	loginsService logins.Service
	postsService  posts.Service
}

func NewConsumer(app app.Services) Consumer {
	return &consumer{
		loginsService: app.LoginsService,
		postsService:  app.PostsService,
	}
}

func (consumer *consumer) Register() pkg.TopicConsumerMap {

	return pkg.TopicConsumerMap{
		pkg.TopicLogsLogins:  consumer.LogLogins,
		pkg.TopicPostsLogins: consumer.LogPosts,
	}
}
