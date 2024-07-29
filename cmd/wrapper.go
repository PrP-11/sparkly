package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"prp.com/sparkly/internal/app"
	"prp.com/sparkly/internal/connectors"
	"prp.com/sparkly/internal/connectors/services/clock"
)

var appServices app.Services

func setup(cmd *cobra.Command, args []string) error {

	// Initialize MongoDB client
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(cmd.Context(), clientOptions)
	if err != nil {
		panic(err)
	}

	clock := clock.NewClock()
	database := client.Database(os.Getenv("MONGODB_NAME"))

	connector := connectors.NewConnector(database, clock)
	appServices = app.NewServices(connector, clock)

	return nil

}

func cleanup(cmd *cobra.Command, args []string) error {
	return nil
}
