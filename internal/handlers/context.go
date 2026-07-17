package handlers

import (
	"github.com/teodora-lazarevic/Poll-App/internal/services"
)

// Holds the database connection
type AppContext struct {
	PollService *services.PollService
	UserService *services.UserService
	VoteService *services.VoteService
}
