package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/teodora-lazarevic/Poll-App/internal/services"
	"github.com/teodora-lazarevic/Poll-App/internal/utils"
)

// HandleError maps domain errors to standard HTTP responses.
func HandleError(w http.ResponseWriter, err error) {
	switch {

	// 404 Not Found
	case errors.Is(err, services.ErrPollNotFound),
		errors.Is(err, services.ErrPollOptionNotFound),
		errors.Is(err, services.ErrUserNotFound):
		utils.ErrorJSON(w, http.StatusNotFound, err.Error())

	// 400 Bad Request
	case errors.Is(err, services.ErrOptionNotInPoll),
		errors.Is(err, services.ErrMinOptionsRequired),
		errors.Is(err, services.ErrInvalidInput):
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())

	// 401 Unauthorized
	case errors.Is(err, services.ErrUnauthorized),
		errors.Is(err, services.ErrInvalidCreds):
		utils.ErrorJSON(w, http.StatusUnauthorized, err.Error())

	// 409 Conflict (Duplicates / Already Voted / User Exists)
	case errors.Is(err, services.ErrAlreadyVoted),
		errors.Is(err, services.ErrDuplicateOption),
		errors.Is(err, services.ErrUserExists):
		utils.ErrorJSON(w, http.StatusConflict, err.Error())

	// 500 Internal Server Error
	default:
		log.Printf("[INTERNAL ERROR]: %v", err)
		utils.ErrorJSON(w, http.StatusInternalServerError, "An unexpected internal error occurred")
	}
}
