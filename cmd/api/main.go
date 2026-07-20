package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/teodora-lazarevic/Poll-App/internal/db"
	"github.com/teodora-lazarevic/Poll-App/internal/handlers"
	"github.com/teodora-lazarevic/Poll-App/internal/router"
	"github.com/teodora-lazarevic/Poll-App/internal/services"
)

func main() {
	// Loading environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	client := db.InitDB(dbURL)
	defer client.Close()

	userService := services.NewUserService(client)
	pollService := services.NewPollService(client)
	voteService := services.NewVoteService(client)

	appCtx := &handlers.AppContext{
		UserService: userService,
		PollService: pollService,
		VoteService: voteService,
	}
	r := router.SetupRouter(appCtx)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
