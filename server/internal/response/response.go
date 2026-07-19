package response

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"server/internal/constants"

	"github.com/google/uuid"
)

func JSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Error("write JSON response", "error", err)
	}
}

func GetUserIdFromContext(ctx context.Context, w http.ResponseWriter) uuid.UUID {
	userId, ok := ctx.Value(constants.UserIDContextKey).(uuid.UUID)
	if !ok {
		JSON(w, http.StatusUnauthorized, "Not found user id")
		return uuid.Nil
	}
	return userId
}
