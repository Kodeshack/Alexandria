package main

import (
	"alexandria.app/articledb"
	"alexandria.app/models"
	"alexandria.app/server"
)

func main() {
	config := models.NewConfig()

	userStorage, err := models.LoadUserStorage(config.UserStoragePath)
	if err != nil {
		panic(err)
	}

	articledb := articledb.New(config.ContentPath)

	sessionStorage := models.NewSessionStorage()

	server.Start(userStorage, sessionStorage, articledb, config)
}
