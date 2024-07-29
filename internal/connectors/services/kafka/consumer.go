package kafka

import (
	"context"
	"log"
	"os"

	"github.com/IBM/sarama"
)

type ConsumerService interface {
	Consume(ctx context.Context, consumerFunc func(ctx context.Context, message *sarama.ConsumerMessage) error)
}

type consumerService struct {
	consumer          sarama.Consumer
	partitionConsumer sarama.PartitionConsumer
}

func NewConsumerService(topic string) ConsumerService {
	broker := os.Getenv("KAFKA_BROKERS")
	log.Println("BROKER: ", broker)
	consumer, err := sarama.NewConsumer([]string{broker}, nil)
	if err != nil {
		log.Fatal("Failed to start Sarama consumer:", err)
	}

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Failed to start partition consumer:", err)
	}

	return &consumerService{
		consumer:          consumer,
		partitionConsumer: partitionConsumer,
	}
}

func (service *consumerService) Consume(ctx context.Context, consumerFunc func(ctx context.Context, message *sarama.ConsumerMessage) error) {
	for msg := range service.partitionConsumer.Messages() {
		err := consumerFunc(ctx, msg)
		if err != nil {
			log.Printf("Unable to consume message offset %d: %s = %s\n", msg.Offset, string(msg.Key), string(msg.Value))
			continue
		}

		log.Printf("Consumed message offset %d: %s = %s\n", msg.Offset, string(msg.Key), string(msg.Value))
	}
}
