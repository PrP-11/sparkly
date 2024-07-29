package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"prp.com/sparkly/internal/app/posts"
	"prp.com/sparkly/internal/pkg"
)

type PostsHandler interface {
	Log(ctx *gin.Context)
	GetPopularPosts(ctx *gin.Context)
}

type postsHandler struct {
	postsService posts.Service
}

func NewPostsHandler(
	postsService posts.Service,
) PostsHandler {
	return &postsHandler{
		postsService: postsService,
	}
}

func (handler *postsHandler) Log(ctx *gin.Context) {

	var activity pkg.PostInteraction
	if err := ctx.ShouldBindJSON(&activity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	handler.postsService.PushToQueue(ctx, activity)

	// err := handler.postsService.Log(ctx, activity)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	ctx.JSON(http.StatusOK, gin.H{"status": "activity logged"})

}

func (handler *postsHandler) GetPopularPosts(ctx *gin.Context) {

	limit := 10
	if limitStr, exists := ctx.GetQuery("limit"); exists {
		if val, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			limit = int(val)
		}
	}

	key, exists := ctx.GetQuery("key")
	if exists {
		handler.getPopularPostsByDuration(ctx, key, limit)
		return
	}

	data, err := handler.postsService.GetPopularPosts(ctx, pkg.AnalyticsTimeFrames, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)

}

func (handler *postsHandler) getPopularPostsByDuration(ctx *gin.Context, key string, limit int) {

	data, err := handler.postsService.GetPopularPostsByDuration(ctx, key, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)

}
