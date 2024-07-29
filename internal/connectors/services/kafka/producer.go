package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"
)

type ProducerService interface {
	PushMessage(ctx context.Context, topic string, body any)
}

type producerService struct {
	producer sarama.SyncProducer
	config   *sarama.Config
}

func NewProducerService() ProducerService {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Net.DialTimeout = 20 * time.Second
	config.Net.ReadTimeout = 20 * time.Second
	config.Net.WriteTimeout = 20 * time.Second

	broker := os.Getenv("KAFKA_BROKERS")
	log.Println("BROKER: ", broker)

	producer, err := sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		log.Fatal("Failed to start Sarama producer:", err)
	} else {
		log.Println("Kafka producer started successfully")
	}

	return &producerService{
		config:   config,
		producer: producer,
	}

}

func (service *producerService) PushMessage(ctx context.Context, topic string, body any) {
	bodyBytes, _ := json.Marshal(body)
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(bodyBytes),
	}

	partition, offset, err := service.producer.SendMessage(msg)
	if err != nil {
		log.Fatal("Failed to send message:", err)
	}

	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
}
