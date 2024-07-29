package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"prp.com/sparkly/internal/ports/rest/handlers"
)

func Register(ctx context.Context, handler handlers.Handler) *gin.Engine {

	router := gin.Default()

	basePath := router.Group("/api")

	v1 := basePath.Group("/v1")
	{
		activity := v1.Group("/activity")
		{
			activity.POST("/logins", handler.LoginsHandler().Log)
			activity.POST("/posts", handler.PostsHandler().Log)
		}

		analysis := v1.Group("/analysis")
		{
			analysis.GET("/active-users", handler.LoginsHandler().GetActiveUsers)
			analysis.GET("/popular-posts", handler.PostsHandler().GetPopularPosts)
		}

	}

	return router
}
