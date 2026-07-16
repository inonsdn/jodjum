package auth

import "net/http"

func RegisterRoutes(router *http.ServeMux, handler *AuthHandler) {
	router.HandleFunc("POST /api/v1/login", handler.Login)
}
