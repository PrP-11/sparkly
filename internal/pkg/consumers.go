package pkg

import (
	"context"

	"github.com/IBM/sarama"
)

type TopicConsumerMap map[string]func(ctx context.Context, message *sarama.ConsumerMessage) error

const TopicLogsLogins string = "principal.sparkly-services.logins.log"
const TopicPostsLogins string = "principal.sparkly-services.posts.log"
