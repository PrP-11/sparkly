package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"prp.com/sparkly/internal/app/logins"
	"prp.com/sparkly/internal/pkg"
)

type LoginsHandler interface {
	Log(ctx *gin.Context)
	GetActiveUsers(ctx *gin.Context)
}

type loginsHandler struct {
	loginsService logins.Service
}

func NewLoginsHandler(
	loginsService logins.Service,
) LoginsHandler {
	return &loginsHandler{
		loginsService: loginsService,
	}
}

func (handler *loginsHandler) Log(ctx *gin.Context) {

	var activity pkg.LoginActivity
	if err := ctx.ShouldBindJSON(&activity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// err := handler.loginsService.Log(ctx, activity)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	handler.loginsService.PushToQueue(ctx, activity)

	ctx.JSON(http.StatusOK, gin.H{"status": "activity logged"})

}

func (handler *loginsHandler) GetActiveUsers(ctx *gin.Context) {

	key, exists := ctx.GetQuery("key")
	if exists {
		handler.getActiveUsersByDuration(ctx, key)
		return
	}

	data, err := handler.loginsService.GetActiveUsers(ctx, pkg.AnalyticsTimeFrames)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)

}

func (handler *loginsHandler) getActiveUsersByDuration(ctx *gin.Context, key string) {

	data, err := handler.loginsService.GetActiveUsersByDuration(ctx, key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)

}
