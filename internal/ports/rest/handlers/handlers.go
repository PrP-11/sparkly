package handlers

import "prp.com/sparkly/internal/app"

type Handler interface {
	LoginsHandler() LoginsHandler
	PostsHandler() PostsHandler
}

type handler struct {
	loginsHandler LoginsHandler
	postsHandler  PostsHandler
}

func NewHandler(
	app app.Services,
) Handler {
	return &handler{
		loginsHandler: NewLoginsHandler(app.LoginsService),
		postsHandler:  NewPostsHandler(app.PostsService),
	}
}

func (h *handler) LoginsHandler() LoginsHandler {
	return h.loginsHandler
}

func (h *handler) PostsHandler() PostsHandler {
	return h.postsHandler
}
