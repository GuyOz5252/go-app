package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GuyOz5252/go-app/internal/core"
	"github.com/GuyOz5252/go-app/internal/services"
	"github.com/GuyOz5252/go-app/pkg"
	"github.com/go-chi/chi/v5"
)

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
		pkg.ApiError(w, r, http.StatusBadRequest, "invalid user id", err.Error())
		return
	}

	user, err := userService.GetById(userId)
	if err != nil {
		if err == core.ErrNotFound {
			pkg.ApiError(w, r, http.StatusNotFound, "user not found", err.Error())
			return
		}
		pkg.ApiError(w, r, http.StatusInternalServerError, "failed to get user", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		pkg.ApiError(w, r, http.StatusInternalServerError, "failed to encode user", err.Error())
		return
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user core.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	userId, err := userService.Create(&user)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "{userId: \"%d\"}", userId)
}
