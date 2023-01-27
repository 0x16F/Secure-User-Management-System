package main

import (
	"fmt"
	"log"
	"test-project/src/pkg/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

func main() {
	// init config
	if err := config.NewConfig(); err != nil {
		panic(err)
	}

	cfg := &config.Database{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetUint16("database.port"),
		Schema:   viper.GetString("database.schema"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Schema)

	m, err := migrate.New("file://migrations", url)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		log.Println(err)
	}
}
