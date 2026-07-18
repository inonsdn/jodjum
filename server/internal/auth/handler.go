package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"server/internal/response"
)

type AuthHandler struct {
	service *AuthService
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// construct decoder object for decode body request to request as expected
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// verify request that required argument in LoginRequest
	//	if cannot decode to login request, return from function
	var loginRequest LoginRequest
	if err := decoder.Decode(&loginRequest); err != nil {
		response.JSON(w, http.StatusBadRequest, "Mismatch request")
		return
	}

	// do service to handle login
	result, err := h.service.Login(r.Context(), loginRequest.Email, loginRequest.Password)

	// found error when authentication
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.JSON(w, http.StatusBadRequest, "Invalid email or password")
		} else {
			response.JSON(w, http.StatusBadRequest, "Invalid")
		}
		return
	}

	// login success
	response.JSON(w, http.StatusOK, result)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userId, ok := UserIdFromContext(r.Context())
	if !ok {
		response.JSON(w, http.StatusBadRequest, "Invalid email or password")
		return
	}

	// do service to handle login
	err := h.service.Logout(r.Context(), userId)

	// found error when authentication
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.JSON(w, http.StatusBadRequest, "Invalid email or password")
		} else {
			response.JSON(w, http.StatusBadRequest, "Invalid")
		}
		return
	}

	// login success
	response.JSON(w, http.StatusOK, "")
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// construct decoder object for decode body request to request as expected
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// verify request that required argument in RegisterRequest
	//	if cannot decode to login request, return from function
	var registerRequest RegisterRequest
	if err := decoder.Decode(&registerRequest); err != nil {
		response.JSON(w, http.StatusBadRequest, "Mismatch request")
		return
	}

	// do service to handle login
	result, err := h.service.RegisterUser(r.Context(), registerRequest.Username, registerRequest.Email, registerRequest.Password)

	// found error when authentication
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.JSON(w, http.StatusBadRequest, "Invalid email or password")
		} else {
			response.JSON(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	// login success
	response.JSON(w, http.StatusOK, result)
}
