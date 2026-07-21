package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/teodora-lazarevic/Poll-App/internal/handlers"
	"github.com/teodora-lazarevic/Poll-App/internal/middleware"
)

func SetupRouter(appCtx *handlers.AppContext) *httprouter.Router {
	router := httprouter.New()

	// health check endpoint
	router.GET("/health", handlers.HealthCheckHandler)

	// poll endpoints
	router.POST("/polls", middleware.RequireAuth(appCtx.CreatePollHandler))
	router.POST("/polls/:poll_id/options", middleware.RequireAuth(appCtx.AddPollOptionHandler))

	router.GET("/polls", appCtx.ListPollsHandler)
	router.GET("/polls/:poll_id", appCtx.GetPollByIdHandler)

	router.DELETE("/polls/:poll_id", middleware.RequireAuth(appCtx.DeletePollHandler))
	router.DELETE("/polls/:poll_id/options/:option_id", middleware.RequireAuth(appCtx.DeletePollOptionHandler))
	router.DELETE("/polls/clear", middleware.RequireAuth(appCtx.ClearAllDataHandler))

	// auth endpoints
	router.POST("/users/register", appCtx.RegisterHandler)
	router.POST("/users/login", appCtx.LoginHandler)

	router.GET("/users", appCtx.GetAllUsersHandler)
	router.GET("/users/:user_id", appCtx.GetUserByIdHandler)

	// vote endpoints
	router.POST("/polls/:poll_id/vote", middleware.RequireAuth(appCtx.VoteHandler))

	router.GET("/polls/:poll_id/results", appCtx.GetPollResultsHandler)
	router.GET("/polls/:poll_id/options/:option_id/voters", middleware.RequireAuth(appCtx.GetVotersHandler))

	return router
}
