package things

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"server/internal/response"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ThingsHandler struct {
	service *ThingsService
}

type CreateThingsRequest struct {
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	Quantity            int        `json:"quantity"`
	ExpiredAt           *time.Time `json:"expires_at"` // nil when the user leaves expiry blank
	NextRemindTimestamp *float64   `json:"next_remind_timestamp"`
}

type UpdateThingsRequest struct {
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	Quantity            int        `json:"quantity"`
	ExpiredAt           *time.Time `json:"expires_at"` // nil when the user leaves expiry blank
	NextRemindTimestamp *float64   `json:"next_remind_timestamp"`
}

func NewHandler(service *ThingsService) *ThingsHandler {
	return &ThingsHandler{
		service: service,
	}
}

func (h *ThingsHandler) GetMyThings(w http.ResponseWriter, r *http.Request) {
	userId := response.GetUserIdFromContext(r.Context(), w)

	thingsList, err := h.service.GetThingsOfUser(r.Context(), userId)

	if err != nil {
		slog.Error("Error when get things", "Error", err.Error())
		response.JSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, thingsList)
}

func (h *ThingsHandler) GetThings(w http.ResponseWriter, r *http.Request) {
	userId := response.GetUserIdFromContext(r.Context(), w)

	thingsId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Invalid things id")
		return
	}

	things, err := h.service.GetThingsById(r.Context(), userId, thingsId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.JSON(w, http.StatusNotFound, "Things not found")
			return
		}
		slog.Error("Error when get things by id", "Error", err.Error())
		response.JSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, things)
}

func (h *ThingsHandler) CreateThings(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userId := response.GetUserIdFromContext(r.Context(), w)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req CreateThingsRequest
	if err := decoder.Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "Mismatch request")
		return
	}

	things, err := h.service.CreateThings(r.Context(), userId, req.Name, req.Description, req.Quantity, req.ExpiredAt, req.NextRemindTimestamp)
	if err != nil {
		slog.Error("Error when create things", "Error", err.Error())
		response.JSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusCreated, things)
}

func (h *ThingsHandler) DeleteThings(w http.ResponseWriter, r *http.Request) {
	userId := response.GetUserIdFromContext(r.Context(), w)

	thingsId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Invalid things id")
		return
	}

	err = h.service.DeleteThings(r.Context(), userId, thingsId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.JSON(w, http.StatusNotFound, "Things not found")
			return
		}
		slog.Error("Error when delete things", "Error", err.Error())
		response.JSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, "deleted")
}

func (h *ThingsHandler) UpdateThings(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userId := response.GetUserIdFromContext(r.Context(), w)

	thingsId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Invalid things id")
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req UpdateThingsRequest
	if err := decoder.Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "Mismatch request")
		return
	}

	things, err := h.service.UpdateThings(r.Context(), userId, thingsId, req.Name, req.Description, req.Quantity, req.ExpiredAt, req.NextRemindTimestamp)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.JSON(w, http.StatusNotFound, "Things not found")
			return
		}
		slog.Error("Error when update things", "Error", err.Error())
		response.JSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, things)
}
