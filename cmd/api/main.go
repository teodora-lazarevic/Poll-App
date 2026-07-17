package main

import (
	"log"
	"net/http"

	"github.com/teodora-lazarevic/Poll-App/internal/db"
	"github.com/teodora-lazarevic/Poll-App/internal/handlers"
	"github.com/teodora-lazarevic/Poll-App/internal/router"
	"github.com/teodora-lazarevic/Poll-App/internal/services"
)

func main() {

	dbURL := "host=localhost port=5432 user=postgres dbname=pollingapp password=secret sslmode=disable"

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
