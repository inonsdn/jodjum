package user

import (
	"net/http"
	"server/internal/types"
)

func RegisterRoutes(router *http.ServeMux, middleware types.Middleware, handler *UserHandler) {
	router.Handle("GET /api/v1/me", middleware(http.HandlerFunc(handler.GetMe)))
}
