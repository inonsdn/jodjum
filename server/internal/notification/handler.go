package notification

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"server/internal/response"
)

type NotificationHandler struct {
	service *NotificationService
}

// SubscribeRequest mirrors the JSON shape of a browser PushSubscription
// (PushSubscription.toJSON()): { endpoint, keys: { p256dh, auth } }.
type SubscribeRequest struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

type UnsubscribeRequest struct {
	Endpoint string `json:"endpoint"`
}

func NewHandler(service *NotificationService) *NotificationHandler {
	return &NotificationHandler{
		service: service,
	}
}

func (h *NotificationHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userId := response.GetUserIdFromContext(r.Context(), w)

	// The browser payload also carries an "expirationTime" field we don't use,
	// so we deliberately do NOT DisallowUnknownFields here.
	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "Mismatch request")
		return
	}

	if req.Endpoint == "" || req.Keys.P256dh == "" || req.Keys.Auth == "" {
		response.JSON(w, http.StatusBadRequest, "endpoint and keys (p256dh, auth) are required")
		return
	}

	sub, err := h.service.Subscribe(r.Context(), userId, req.Endpoint, req.Keys.P256dh, req.Keys.Auth)
	if err != nil {
		slog.Error("Error when subscribe push notification", "Error", err.Error())
		response.JSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusCreated, sub)
}

func (h *NotificationHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userId := response.GetUserIdFromContext(r.Context(), w)

	var req UnsubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "Mismatch request")
		return
	}

	if req.Endpoint == "" {
		response.JSON(w, http.StatusBadRequest, "endpoint is required")
		return
	}

	if err := h.service.Unsubscribe(r.Context(), userId, req.Endpoint); err != nil {
		slog.Error("Error when unsubscribe push notification", "Error", err.Error())
		response.JSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, "unsubscribed")
}
