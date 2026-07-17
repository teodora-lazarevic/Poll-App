package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/teodora-lazarevic/Poll-App/ent"
	"github.com/teodora-lazarevic/Poll-App/internal/services"
)

// Votes for an option in a poll
func (appCtx *AppContext) VoteHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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
		OptionID int `json:"option_id"`
	}

	if err := ReadJSON(writer, request, &req); err != nil {
		ErrorJSON(writer, http.StatusBadRequest, "Invalid request body")
		return
	} else if req.OptionID == 0 {
		ErrorJSON(writer, http.StatusBadRequest, "Option ID is required")
		return
	}

	err = appCtx.VoteService.CastVote(request.Context(), userId, pollId, req.OptionID)
	if err != nil {
		handleVoteError(writer, err)
		return
	}

	WriteJSON(writer, http.StatusOK, "Vote created successfully")
}

// Gets the results of a poll
func (appCtx *AppContext) GetPollResultsHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pollId, err := strconv.Atoi(params.ByName("poll_id"))
	if err != nil {
		ErrorJSON(writer, http.StatusBadRequest, "Invalid poll ID")
		return
	}

	results, err := appCtx.VoteService.GetPollResults(request.Context(), pollId)
	if err != nil {
		handleVoteError(writer, err)
		return
	}

	WriteJSON(writer, http.StatusOK, results)
}

// Gets the voters for an option
func (appCtx *AppContext) GetVotersHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	voters, err := appCtx.VoteService.GetVotersForOption(request.Context(), userId, pollId, optionId)
	if err != nil {
		handleVoteError(writer, err)
		return
	}

	WriteJSON(writer, http.StatusOK, voters)
}

// Helper to map vote service errors neatly to standard HTTP status codes
func handleVoteError(writer http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrAlreadyVoted):
		ErrorJSON(writer, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrOptionNotInPoll):
		ErrorJSON(writer, http.StatusNotFound, err.Error())
	case errors.Is(err, services.ErrUnauthorized):
		ErrorJSON(writer, http.StatusUnauthorized, "User is not the creator of the poll")
	case errors.Is(err, services.ErrPollNotFound), ent.IsNotFound(err):
		ErrorJSON(writer, http.StatusNotFound, "Resource not found")
	default:
		ErrorJSON(writer, http.StatusInternalServerError, "An unexpected error occurred")
	}
}
