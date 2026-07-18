package user

import (
	"log/slog"
	"net/http"
)

type UserHandler struct {
	service *UserService
}

func NewHandler(service *UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id")
	slog.Info("Get info", "user", userId)
}
