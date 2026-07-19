package reminders

import (
	"net/http"
	"server/internal/types"
)

func RegisterRoutes(router *http.ServeMux, middleware types.Middleware, handler *RemindersHandler) {
	router.Handle("GET /api/v1/myReminders", middleware(http.HandlerFunc(handler.GetMyReminders)))
	router.Handle("GET /api/v1/reminders/{id}", middleware(http.HandlerFunc(handler.GetReminder)))
	router.Handle("POST /api/v1/reminders", middleware(http.HandlerFunc(handler.CreateReminder)))
	router.Handle("PUT /api/v1/reminders/{id}", middleware(http.HandlerFunc(handler.UpdateReminder)))
}
