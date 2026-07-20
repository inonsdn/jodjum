package notification

import (
	"net/http"
	"server/internal/types"
)

func RegisterRoutes(router *http.ServeMux, middleware types.Middleware, handler *NotificationHandler) {
	// push subscription management:
	router.Handle("POST /api/v1/push/subscribe", middleware(http.HandlerFunc(handler.Subscribe)))
	router.Handle("POST /api/v1/push/unsubscribe", middleware(http.HandlerFunc(handler.Unsubscribe)))

}
