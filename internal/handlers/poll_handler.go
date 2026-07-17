package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/teodora-lazarevic/Poll-App/ent"
	"github.com/teodora-lazarevic/Poll-App/internal/services"
)

// Returns all polls
func (appCtx *AppContext) ListPollsHandler(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	polls, err := appCtx.PollService.ListPolls(request.Context())
	if err != nil {
		ErrorJSON(writer, http.StatusInternalServerError, "Failed to fetch polls")
		return
	}

	WriteJSON(writer, http.StatusOK, polls)
}

// Creates a new poll
func (appCtx *AppContext) CreatePollHandler(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	userId, ok := authenticate(writer, request)
	if !ok {
		return
	}

	var req struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Options     []string `json:"options"`
	}

	if err := ReadJSON(writer, request, &req); err != nil {
		ErrorJSON(writer, http.StatusBadRequest, "Invalid request body")
		return
	} else if req.Title == "" || len(req.Options) < 2 {
		ErrorJSON(writer, http.StatusBadRequest, "A title and at least two options are required")
		return
	}

	_, err := appCtx.PollService.CreatePoll(request.Context(), userId, req.Title, req.Description, req.Options)
	if err != nil {
		handlePollError(writer, err)
		return
	}

	WriteJSON(writer, http.StatusCreated, "Poll created successfully")
}

// Returns a poll by its ID
func (appCtx *AppContext) GetPollByIdHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pollId, err := strconv.Atoi(params.ByName("poll_id"))
	if err != nil {
		ErrorJSON(writer, http.StatusBadRequest, "Invalid poll ID")
		return
	}

	poll, err := appCtx.PollService.GetPollById(request.Context(), pollId)
	if err != nil {
		handlePollError(writer, err)
		return
	}

	WriteJSON(writer, http.StatusOK, poll)
}

// Adds a new option to a poll
func (appCtx *AppContext) AddPollOptionHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId, ok := authenticate(writer, request)
	if !ok {
		return
	}

	pollId, err := strconv.Atoi(params.ByName("poll_id"))
	if err != nil {
		ErrorJSON(writer, http.StatusBadRequest, "Invalid poll ID")
		return
	}

	var req struct {
		Text string `json:"text"`
	}

	if err := ReadJSON(writer, request, &req); err != nil || req.Text == "" {
		ErrorJSON(writer, http.StatusBadRequest, "Option text is required")
		return
	}

	_, err = appCtx.PollService.AddPollOption(request.Context(), userId, pollId, req.Text)
	if err != nil {
		handlePollError(writer, err)
		return
	}

	WriteJSON(writer, http.StatusCreated, "Option created successfully")
}

// Deletes a poll by its ID
func (appCtx *AppContext) DeletePollHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId, ok := authenticate(writer, request)
	if !ok {
		return
	}

	pollId, err := strconv.Atoi(params.ByName("poll_id"))
	if err != nil {
		ErrorJSON(writer, http.StatusBadRequest, "Invalid poll ID")
		return
	}

	err = appCtx.PollService.DeletePoll(request.Context(), userId, pollId)
	if err != nil {
		handlePollError(writer, err)
		return
	}

	WriteJSON(writer, http.StatusOK, "Poll deleted successfully")
}

// Deletes an option by its ID
func (appCtx *AppContext) DeletePollOptionHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId, ok := authenticate(writer, request)
	if !ok {
		return
	}

	pollId, err1 := strconv.Atoi(params.ByName("poll_id"))
	optionId, err2 := strconv.Atoi(params.ByName("option_id"))
	if err1 != nil || err2 != nil {
		ErrorJSON(writer, http.StatusBadRequest, "Invalid poll or option ID")
		return
	}

	err := appCtx.PollService.DeletePollOption(request.Context(), userId, pollId, optionId)
	if err != nil {
		handlePollError(writer, err)
		return
	}

	WriteJSON(writer, http.StatusOK, "Option deleted successfully")

}

func handlePollError(writer http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrUnauthorized):
		ErrorJSON(writer, http.StatusUnauthorized, err.Error())
	case errors.Is(err, services.ErrDuplicateOption):
		ErrorJSON(writer, http.StatusConflict, err.Error())
	case errors.Is(err, services.ErrPollNotFound), ent.IsNotFound(err):
		ErrorJSON(writer, http.StatusNotFound, "Poll not found")
	case errors.Is(err, services.ErrPollOptionNotFound):
		ErrorJSON(writer, http.StatusNotFound, "Option not found")
	default:
		ErrorJSON(writer, http.StatusInternalServerError, "An unexpected error occurred")
	}
}
