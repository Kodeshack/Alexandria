package main

import (
	"alexandria.app/models"
	"alexandria.app/server"
)

func main() {
	config := models.NewConfig()

	userStorage, err := models.LoadUserStorage(config.UserStoragePath)
	if err != nil {
		panic(err)
	}

	sessionStorage := models.NewSessionStorage()

	server.Start(userStorage, sessionStorage, config)
}
