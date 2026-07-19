package things

import (
	"net/http"
	"server/internal/types"
)

func RegisterRoutes(router *http.ServeMux, middleware types.Middleware, handler *ThingsHandler) {
	router.Handle("GET /api/v1/myThings", middleware(http.HandlerFunc(handler.GetMyThings)))
	router.Handle("GET /api/v1/things/{id}", middleware(http.HandlerFunc(handler.GetThings)))
	router.Handle("POST /api/v1/things", middleware(http.HandlerFunc(handler.CreateThings)))
	router.Handle("PUT /api/v1/things/{id}", middleware(http.HandlerFunc(handler.UpdateThings)))
	router.Handle("DELETE /api/v1/things/{id}", middleware(http.HandlerFunc(handler.DeleteThings)))

}
