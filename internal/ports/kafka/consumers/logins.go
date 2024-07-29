package consumers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"prp.com/sparkly/internal/pkg"
)

func (consumer *consumer) LogLogins(ctx context.Context, message *sarama.ConsumerMessage) error {

	var payload = new(pkg.LoginActivity)
	if err := json.Unmarshal(message.Value, payload); err != nil {
		log.Println("Error unmarshalling logins data", err)
		return err
	}

	return consumer.loginsService.Log(ctx, *payload)

}
