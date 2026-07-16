package user

import "net/http"

func RegisterRoutes(router *http.ServeMux, handler *UserHandler) {
	router.HandleFunc("GET /api/v1/me", handler.GetMe)
}
