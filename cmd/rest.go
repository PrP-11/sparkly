package cmd

import (
	"github.com/spf13/cobra"
	"prp.com/sparkly/internal/ports/rest/handlers"
	"prp.com/sparkly/internal/ports/rest/routes"
)

var restCommand = &cobra.Command{
	Use:   "rest",
	Short: "Starts rest server",
	RunE: func(cmd *cobra.Command, args []string) error {

		handler := handlers.NewHandler(appServices)
		router := routes.Register(cmd.Context(), handler)

		appServices.LoginsService.BackfillActiveUsers(cmd.Context())
		appServices.PostsService.BackfillPolularPosts(cmd.Context())

		// Run the server
		router.Run(":8080")

		return nil
	},
}
