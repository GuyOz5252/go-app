package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GuyOz5252/go-app/internal/core"
	"github.com/GuyOz5252/go-app/internal/services"
	api "github.com/GuyOz5252/go-app/pkg/api_utils"
	"github.com/go-chi/chi/v5"
)

type CreateUserResponse struct {
	UserId int `json:"userId"`
}

var userService *services.UserService

func mountUserRoutes(app *application) http.Handler {
	userService = app.UserService
	mux := chi.NewRouter()

	mux.Post("/", createUserHandler)
	mux.Get("/{id}", getUserByIdHandler)

	return mux
}

func getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	userIdString := chi.URLParam(r, "id")
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		api.ApiError(w, r, http.StatusBadRequest, "invalid user id", err.Error())
		return
	}

	user, err := userService.GetById(userId)
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

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user core.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	userId, err := userService.Create(&user)
	if err != nil {
		if err == core.ErrUsernameConflict {
			api.ApiError(w, r, http.StatusConflict, "username already exists", err.Error())
			return
		}
		if err == core.ErrEmailConflict {
			api.ApiError(w, r, http.StatusConflict, "email already exists", err.Error())
			return
		}

		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	api.ApiResponse(w, r, http.StatusOK, CreateUserResponse{
		UserId: userId,
	})
}
