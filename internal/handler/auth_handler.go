package handler

import (
	"encoding/json"
	"net/http"

	"dailystep-backend/internal/service"
	"dailystep-backend/internal/transport/dto"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "invalid JSON body",
		})
		return
	}

	resp, err := h.authService.Register(r.Context(), req)
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{
				Error:   "invalid_request",
				Message: "email and password are invalid",
			})
		case service.ErrEmailAlreadyUsed:
			writeJSON(w, http.StatusConflict, dto.ErrorResponse{
				Error:   "conflict",
				Message: "email already in use",
			})
		default:
			writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "something went wrong",
			})
		}
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "invalid JSON body",
		})
		return
	}

	resp, err := h.authService.Login(r.Context(), req)
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{
				Error:   "invalid_request",
				Message: "email and password are required",
			})
		case service.ErrInvalidCredentials:
			writeJSON(w, http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "unauthorized",
				Message: "invalid credentials",
			})
		default:
			writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "something went wrong",
			})
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
