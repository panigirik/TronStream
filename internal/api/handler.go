package api

import (
	"TronStream/internal/infrastructure/auth"
	"encoding/json"
	"net/http"
)

type Handler struct {
	AuthService *auth.Service
}

func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/auth/refresh", h.Refresh)
	mux.HandleFunc("POST /v1/auth/sign-in", h.SignIn)
}

func HealthHandler(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\":\"ok\"}"))
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	if h.AuthService == nil {
		http.Error(w, "auth service is not configured", http.StatusServiceUnavailable)
		return
	}

	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	accessToken, err := h.AuthService.SignIn(r.Context(), request.Email, request.Password)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	if h.AuthService == nil || h.AuthService.TokenService == nil {
		http.Error(w, "auth service is not configured", http.StatusServiceUnavailable)
		return
	}

	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tokenPair, err := h.AuthService.TokenService.GenerateTokenPair(r.Context(), request.RefreshToken)
	if err != nil {
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tokenPair)
}
