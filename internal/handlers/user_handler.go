package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/teodora-lazarevic/Poll-App/ent"
	"github.com/teodora-lazarevic/Poll-App/internal/services"
)

// Registers a new user
func (appCtx *AppContext) RegisterHandler(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ReadJSON(writer, request, &req); err != nil {
		ErrorJSON(writer, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := appCtx.UserService.Register(request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		handleUserError(writer, err)
		return
	}

	WriteJSON(writer, http.StatusCreated, "User registered successfully")
}

// Logs in a user
func (appCtx *AppContext) LoginHandler(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var req struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	if err := ReadJSON(writer, request, &req); err != nil {
		ErrorJSON(writer, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := appCtx.UserService.Authenticate(request.Context(), req.Identifier, req.Password)
	if err != nil {
		handleUserError(writer, err)
		return
	}

	token, err := GenerateJWTToken(user.ID)
	if err != nil {
		ErrorJSON(writer, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	WriteJSON(writer, http.StatusOK, map[string]string{"token": token})
}

// Get all users
func (appCtx *AppContext) GetAllUsersHandler(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	users, err := appCtx.UserService.GetAllUsers(request.Context())
	if err != nil {
		ErrorJSON(writer, http.StatusInternalServerError, "Failed to get users")
		return
	}
	WriteJSON(writer, http.StatusOK, users)
}

// Get a user by ID
func (appCtx *AppContext) GetUserByIdHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId, err := strconv.Atoi(params.ByName("user_id"))
	if err != nil {
		ErrorJSON(writer, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := appCtx.UserService.GetUserById(request.Context(), userId)
	if err != nil {
		handleUserError(writer, err)
		return
	}

	WriteJSON(writer, http.StatusOK, user)
}

// Helper to map service errors neatly to standard HTTP status codes
func handleUserError(writer http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidInput):
		ErrorJSON(writer, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrUserExists):
		ErrorJSON(writer, http.StatusConflict, err.Error())
	case errors.Is(err, services.ErrInvalidCreds):
		ErrorJSON(writer, http.StatusUnauthorized, err.Error())
	case errors.Is(err, services.ErrUserNotFound), ent.IsNotFound(err):
		ErrorJSON(writer, http.StatusNotFound, "User not found")
	default:
		ErrorJSON(writer, http.StatusInternalServerError, "An unexpected error occurred")
	}
}
