package handlers

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/teodora-lazarevic/Poll-App/internal/middleware"
	"github.com/teodora-lazarevic/Poll-App/internal/utils"
)

// Returns all polls
func (appCtx *AppContext) ListPollsHandler(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	polls, err := appCtx.PollService.ListPolls(request.Context())
	if err != nil {
		HandleError(writer, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, polls)
}

// Creates a new poll
func (appCtx *AppContext) CreatePollHandler(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	userId, ok := middleware.GetUserIDFromContext(request.Context())
	if !ok {
		return
	}

	var req struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Options     []string `json:"options"`
	}

	if err := utils.ReadJSON(writer, request, &req); err != nil {
		utils.ErrorJSON(writer, http.StatusBadRequest, "Invalid request body")
		return
	} else if req.Title == "" || len(req.Options) < 2 {
		utils.ErrorJSON(writer, http.StatusBadRequest, "A title and at least two options are required")
		return
	}

	_, err := appCtx.PollService.CreatePoll(request.Context(), userId, req.Title, req.Description, req.Options)
	if err != nil {
		HandleError(writer, err)
		return
	}

	utils.WriteJSON(writer, http.StatusCreated, "Poll created successfully")
}

// Returns a poll by its ID
func (appCtx *AppContext) GetPollByIdHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pollId, err := strconv.Atoi(params.ByName("poll_id"))
	if err != nil {
		utils.ErrorJSON(writer, http.StatusBadRequest, "Invalid poll ID")
		return
	}

	poll, err := appCtx.PollService.GetPollById(request.Context(), pollId)
	if err != nil {
		HandleError(writer, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, poll)
}

// Adds a new option to a poll
func (appCtx *AppContext) AddPollOptionHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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
		Text string `json:"text"`
	}

	if err := utils.ReadJSON(writer, request, &req); err != nil || req.Text == "" {
		utils.ErrorJSON(writer, http.StatusBadRequest, "Option text is required")
		return
	}

	_, err = appCtx.PollService.AddPollOption(request.Context(), userId, pollId, req.Text)
	if err != nil {
		HandleError(writer, err)
		return
	}

	utils.WriteJSON(writer, http.StatusCreated, "Option created successfully")
}

// Deletes a poll by its ID
func (appCtx *AppContext) DeletePollHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	err = appCtx.PollService.DeletePoll(request.Context(), userId, pollId)
	if err != nil {
		HandleError(writer, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, "Poll deleted successfully")
}

// Deletes an option by its ID
func (appCtx *AppContext) DeletePollOptionHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	err := appCtx.PollService.DeletePollOption(request.Context(), userId, pollId, optionId)
	if err != nil {
		HandleError(writer, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, "Option deleted successfully")

}

func (appCtx *AppContext) ClearAllDataHandler(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	err := appCtx.PollService.ClearAllData(request.Context())
	if err != nil {
		HandleError(writer, err)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, "All data cleared successfully")
}
