package main

import (
	"test-project/src/controller/http"
	"test-project/src/controller/repository"
	"test-project/src/pkg/config"
	"test-project/src/pkg/jwt"
)

func main() {
	// init config
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// init jwt service
	jwtService := jwt.NewService(&jwt.Service{
		AccessSecret:  config.JWT.AccessSecret,
		RefreshSecret: config.JWT.RefreshSecret,
	})

	// init database
	database := repository.NewDatabase()

	// connect to database
	storage, err := database.Connect(&config.Database)
	if err != nil {
		panic(err)
	}

	// init http server
	server := http.NewServer(storage, jwtService)

	// start listening
	if err := server.Start(config.HTTP.Port); err != nil {
		panic(err)
	}
}
