package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/teodora-lazarevic/Poll-App/internal/handlers"
)

func SetupRouter(appCtx *handlers.AppContext) *httprouter.Router {
	router := httprouter.New()

	// health check endpoint
	router.GET("/health", handlers.HealthCheckHandler)

	// poll endpoints
	router.POST("/polls", appCtx.CreatePollHandler)
	router.POST("/polls/:poll_id/options", appCtx.AddPollOptionHandler)

	router.GET("/polls", appCtx.ListPollsHandler)
	router.GET("/polls/:poll_id", appCtx.GetPollByIdHandler)

	router.DELETE("/polls/:poll_id", appCtx.DeletePollHandler)
	router.DELETE("/polls/:poll_id/options/:option_id", appCtx.DeletePollOptionHandler)

	// auth endpoints
	router.POST("/users/register", appCtx.RegisterHandler)
	router.POST("/users/login", appCtx.LoginHandler)

	router.GET("/users", appCtx.GetAllUsersHandler)
	router.GET("/users/:user_id", appCtx.GetUserByIdHandler)

	// vote endpoints
	router.POST("/polls/:poll_id/vote", appCtx.VoteHandler)

	router.GET("/polls/:poll_id/results", appCtx.GetPollResultsHandler)
	router.GET("/polls/:poll_id/options/:option_id/voters", appCtx.GetVotersHandler)

	return router
}
