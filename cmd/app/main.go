package main

import (
	"context"
	"test-project/src/controller/http"
	"test-project/src/controller/repository"
	"test-project/src/pkg/config"
	"test-project/src/pkg/jwt"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/spf13/viper"
)

// @title Test Project
// @version 1.0
// @description Just solving a test task

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// init config
	if err := config.NewConfig(); err != nil {
		panic(err)
	}

	// init jwt service
	jwtService := jwt.NewService(&jwt.Service{
		AccessSecret:  viper.GetString("jwt.access_secret"),
		RefreshSecret: viper.GetString("jwt.refresh_secret"),
	})

	// init database
	database := repository.NewDatabase()

	cfg := &config.Database{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetUint16("database.port"),
		Schema:   viper.GetString("database.schema"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
	}

	// connect to database
	storage, err := database.Connect(cfg)

	if err != nil {
		panic(err)
	}

	// init cache (we can use redis, but there is no need)
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		panic(err)
	}

	// init http server
	server := http.NewServer(storage, cache, jwtService)

	// start listening
	if err := server.Start(viper.GetUint16("http.port")); err != nil {
		panic(err)
	}
}
