package main

import (
	"log"

	"twitter-golang-api/config"
	"twitter-golang-api/db"
	"twitter-golang-api/handlers"
)

func main() {
	cfg := config.LoadConfig()

	client, ctx, cancel := db.Connect(cfg.MongoURI)
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	username := "jinnyboy"
	userResponse, err := handlers.FetchUserData(username, cfg.TwitterBearerToken)
	if err != nil {
		log.Fatal(err)
	}

	handlers.StoreUserData(client, userResponse)
}
