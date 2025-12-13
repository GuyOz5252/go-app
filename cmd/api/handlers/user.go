package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GuyOz5252/go-app/internal/core"
	"github.com/GuyOz5252/go-app/internal/services"
	api "github.com/GuyOz5252/go-app/pkg/api_utils"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type CreateUserResponse struct {
	UserId int `json:"userId"`
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
	var user core.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		api.ApiError(w, r, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	userId, err := h.userService.Create(r.Context(), &user)
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

	api.ApiResponse(w, r, http.StatusCreated, CreateUserResponse{
		UserId: userId,
	})
}
