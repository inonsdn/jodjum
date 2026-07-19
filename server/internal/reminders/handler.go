package reminders

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

type RemindersHandler struct {
	service *RemindersService
}

type CreateReminderRequest struct {
	RemindTimestamp time.Time `json:"remind_timestamp"`
	ReminderType    int       `json:"reminder_type"`
	IsActive        bool      `json:"is_active"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
}

type UpdateReminderRequest struct {
	RemindTimestamp time.Time `json:"remind_timestamp"`
	ReminderType    int       `json:"reminder_type"`
	IsActive        bool      `json:"is_active"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
}

func NewHandler(service *RemindersService) *RemindersHandler {
	return &RemindersHandler{
		service: service,
	}
}

// isValidReminderType guards the reminder_type enum: only these four values are
// allowed, matching the database's reminder_type constraint.
func isValidReminderType(t int) bool {
	switch t {
	case ONE_TIME_REMINDER_TYPE, DAILY_REMINDER_TYPE, MONTHLY_REMINDER_TYPE, YEARLY_REMINDER_TYPE:
		return true
	}
	return false
}

func (h *RemindersHandler) GetMyReminders(w http.ResponseWriter, r *http.Request) {
	userId := response.GetUserIdFromContext(r.Context(), w)

	reminderList, err := h.service.GetRemindersOfUser(r.Context(), userId)
	if err != nil {
		slog.Error("Error when get reminders", "Error", err.Error())
		response.JSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, reminderList)
}

func (h *RemindersHandler) GetReminder(w http.ResponseWriter, r *http.Request) {
	userId := response.GetUserIdFromContext(r.Context(), w)

	reminderId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Invalid reminder id")
		return
	}

	reminder, err := h.service.GetReminderById(r.Context(), userId, reminderId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.JSON(w, http.StatusNotFound, "Reminder not found")
			return
		}
		slog.Error("Error when get reminder by id", "Error", err.Error())
		response.JSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, reminder)
}

func (h *RemindersHandler) CreateReminder(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userId := response.GetUserIdFromContext(r.Context(), w)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req CreateReminderRequest
	if err := decoder.Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "Mismatch request")
		return
	}

	if !isValidReminderType(req.ReminderType) {
		response.JSON(w, http.StatusBadRequest, "reminder_type must be one of: onetime, daily, monthly, yearly")
		return
	}

	reminder, err := h.service.CreateReminder(r.Context(), userId, req.RemindTimestamp, req.ReminderType, req.IsActive, req.Name, req.Description)
	if err != nil {
		slog.Error("Error when create reminder", "Error", err.Error())
		response.JSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusCreated, reminder)
}

func (h *RemindersHandler) UpdateReminder(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userId := response.GetUserIdFromContext(r.Context(), w)

	reminderId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Invalid reminder id")
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req UpdateReminderRequest
	if err := decoder.Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, "Mismatch request")
		return
	}

	if !isValidReminderType(req.ReminderType) {
		response.JSON(w, http.StatusBadRequest, "reminder_type must be one of: onetime, daily, monthly, yearly")
		return
	}

	reminder, err := h.service.UpdateReminder(r.Context(), userId, reminderId, req.RemindTimestamp, req.ReminderType, req.IsActive, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.JSON(w, http.StatusNotFound, "Reminder not found")
			return
		}
		slog.Error("Error when update reminder", "Error", err.Error())
		response.JSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, reminder)
}
