package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"prp.com/sparkly/internal/connectors/services/kafka"
	"prp.com/sparkly/internal/ports/kafka/consumers"
)

var kafkaCommand = &cobra.Command{
	Use:   "kafka",
	Short: "Starts kafka consumers",
	RunE: func(cmd *cobra.Command, args []string) error {

		for topic, consumer := range consumers.NewConsumer(appServices).Register() {
			service := kafka.NewConsumerService(topic)
			go service.Consume(cmd.Context(), consumer)
		}

		// Wait until shutdown signal is received
		s := make(chan os.Signal, 1)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-s
		return nil
	},
}
