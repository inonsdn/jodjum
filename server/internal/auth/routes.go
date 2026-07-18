package auth

import (
	"net/http"
	"server/internal/types"
)

func RegisterRoutes(router *http.ServeMux, middleware types.Middleware, handler *AuthHandler) {
	router.HandleFunc("POST /api/v1/register", handler.Register)
	router.HandleFunc("POST /api/v1/login", handler.Login)
	router.Handle("POST /api/v1/logout", middleware(http.HandlerFunc(handler.Logout)))
}
