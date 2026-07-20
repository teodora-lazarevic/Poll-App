package handlers

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/teodora-lazarevic/Poll-App/internal/middleware"
	"github.com/teodora-lazarevic/Poll-App/internal/utils"
)

// Votes for an option in a poll
func (appCtx *AppContext) VoteHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId, ok := middleware.GetUserIDFromContext(request.Context())
	if !ok {
		utils.ErrorJSON(writer, http.StatusUnauthorized, "Unauthorized")
		return
	}

	pollId, err := strconv.Atoi(params.ByName("poll_id"))
	if err != nil {
		utils.ErrorJSON(writer, http.StatusBadRequest, "Invalid poll ID")
		return
	}

	var req struct {
		OptionID int `json:"option_id"`
	}

	if err := utils.ReadJSON(writer, request, &req); err != nil {
		utils.ErrorJSON(writer, http.StatusBadRequest, "Invalid request body")
		return
	} else if req.OptionID == 0 {
		utils.ErrorJSON(writer, http.StatusBadRequest, "Option ID is required")
		return
	}

	err = appCtx.VoteService.CastVote(request.Context(), userId, pollId, req.OptionID)
	if err != nil {
		HandleError(writer, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, "Vote created successfully")
}

// Gets the results of a poll
func (appCtx *AppContext) GetPollResultsHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pollId, err := strconv.Atoi(params.ByName("poll_id"))
	if err != nil {
		utils.ErrorJSON(writer, http.StatusBadRequest, "Invalid poll ID")
		return
	}

	results, err := appCtx.VoteService.GetPollResults(request.Context(), pollId)
	if err != nil {
		HandleError(writer, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, results)
}

// Gets the voters for an option
func (appCtx *AppContext) GetVotersHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId, ok := middleware.GetUserIDFromContext(request.Context())
	if !ok {
		utils.ErrorJSON(writer, http.StatusUnauthorized, "Unauthorized")
		return
	}

	pollId, err1 := strconv.Atoi(params.ByName("poll_id"))
	optionId, err2 := strconv.Atoi(params.ByName("option_id"))
	if err1 != nil || err2 != nil {
		utils.ErrorJSON(writer, http.StatusBadRequest, "Invalid poll or option ID")
		return
	}

	voters, err := appCtx.VoteService.GetVotersForOption(request.Context(), userId, pollId, optionId)
	if err != nil {
		HandleError(writer, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, voters)
}
