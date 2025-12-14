package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"time"

	"github.com/GuyOz5252/go-app/internal/core"
	"github.com/GuyOz5252/go-app/internal/services"
	api "github.com/GuyOz5252/go-app/pkg/api_utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type UserHandler struct {
	userService     *services.UserService
	tokenAuth       *jwtauth.JWTAuth
	tokenExpiration time.Duration
}

func NewUserHandler(userService *services.UserService, tokenAuth *jwtauth.JWTAuth, tokenExpiration time.Duration) *UserHandler {
	return &UserHandler{
		userService:     userService,
		tokenAuth:       tokenAuth,
		tokenExpiration: tokenExpiration,
	}
}

type RegisterResponse struct {
	UserId int `json:"userId"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userIdString := chi.URLParam(r, "id")
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		api.ApiError(w, r, http.StatusBadRequest, "invalid user id", err.Error())
		return
	}

	user, err := h.userService.GetById(r.Context(), userId)
	if err != nil {
		if err == core.ErrNotFound {
			api.ApiError(w, r, http.StatusNotFound, "user not found", err.Error())
			return
		}
		api.ApiError(w, r, http.StatusInternalServerError, "failed to get user", err.Error())
		return
	}

	api.ApiResponse(w, r, http.StatusOK, user)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ApiError(w, r, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	userId, err := h.userService.Create(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		if err == core.ErrUsernameConflict {
			api.ApiError(w, r, http.StatusConflict, "username already exists", err.Error())
			return
		}
		if err == core.ErrEmailConflict {
			api.ApiError(w, r, http.StatusConflict, "email already exists", err.Error())
			return
		}

		api.ApiError(w, r, http.StatusInternalServerError, "failed to create user", err.Error())
		return
	}

	api.ApiResponse(w, r, http.StatusCreated, RegisterResponse{
		UserId: userId,
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ApiError(w, r, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	user, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if err == core.ErrInvalidCredentials {
			api.ApiError(w, r, http.StatusUnauthorized, "invalid credentials", err.Error())
			return
		}
		api.ApiError(w, r, http.StatusInternalServerError, "failed to login", err.Error())
		return
	}

	claims := map[string]any{
		"userId": user.Id,
	}
	jwtauth.SetExpiryIn(claims, h.tokenExpiration)
	jwtauth.SetIssuedNow(claims)

	_, tokenString, err := h.tokenAuth.Encode(claims)
	if err != nil {
		api.ApiError(w, r, http.StatusInternalServerError, "failed to generate token", err.Error())
		return
	}

	api.ApiResponse(w, r, http.StatusOK, LoginResponse{
		Token: tokenString,
	})
}
