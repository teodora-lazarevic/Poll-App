package handlers

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/teodora-lazarevic/Poll-App/internal/utils"
)

// Registers a new user
func (appCtx *AppContext) RegisterHandler(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := utils.ReadJSON(writer, request, &req); err != nil {
		utils.ErrorJSON(writer, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := appCtx.UserService.Register(request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		HandleError(writer, err)
		return
	}

	utils.WriteJSON(writer, http.StatusCreated, "User registered successfully")
}

// Logs in a user
func (appCtx *AppContext) LoginHandler(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var req struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	if err := utils.ReadJSON(writer, request, &req); err != nil {
		utils.ErrorJSON(writer, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := appCtx.UserService.Authenticate(request.Context(), req.Identifier, req.Password)
	if err != nil {
		HandleError(writer, err)
		return
	}

	token, err := GenerateJWTToken(user.ID)
	if err != nil {
		utils.ErrorJSON(writer, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.WriteJSON(writer, http.StatusOK, map[string]string{"token": token})
}

// Get all users
func (appCtx *AppContext) GetAllUsersHandler(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	users, err := appCtx.UserService.GetAllUsers(request.Context())
	if err != nil {
		utils.ErrorJSON(writer, http.StatusInternalServerError, "Failed to get users")
		return
	}
	utils.WriteJSON(writer, http.StatusOK, users)
}

// Get a user by ID
func (appCtx *AppContext) GetUserByIdHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId, err := strconv.Atoi(params.ByName("user_id"))
	if err != nil {
		utils.ErrorJSON(writer, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := appCtx.UserService.GetUserById(request.Context(), userId)
	if err != nil {
		HandleError(writer, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, user)
}
