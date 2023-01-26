package repository

import (
	"context"
	"test-project/src/internal/user"
	"test-project/src/pkg/config"

	"github.com/go-pg/pg/v10"
)

type Database struct {
	ctx context.Context
}

type Storage struct {
	db    *pg.DB
	Users user.Storager
}

type Databaser interface {
	Connect(cfg *config.Database) (*Storage, error)
}
