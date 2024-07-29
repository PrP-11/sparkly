package consumers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"prp.com/sparkly/internal/pkg"
)

func (consumer *consumer) LogPosts(ctx context.Context, message *sarama.ConsumerMessage) error {

	var payload = new(pkg.PostInteraction)
	if err := json.Unmarshal(message.Value, payload); err != nil {
		log.Println("Error unmarshalling logins data", err)
		return err
	}

	return consumer.postsService.Log(ctx, *payload)
}
