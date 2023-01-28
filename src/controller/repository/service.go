package repository

import (
	"context"
	"fmt"
	"test-project/src/internal/user"
	"test-project/src/pkg/config"

	"github.com/go-pg/pg/v10"
)

func NewDatabase() Databaser {
	return &Database{
		ctx: context.Background(),
	}
}

func (s *Database) Connect(cfg *config.Database) (*Storage, error) {
	db := pg.Connect(&pg.Options{
		Addr:                  fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Database:              cfg.Schema,
		User:                  cfg.User,
		Password:              cfg.Password,
		MaxRetries:            10,
		RetryStatementTimeout: true,
	})

	if err := db.Ping(s.ctx); err != nil {
		return nil, err
	}

	return &Storage{
		db:    db,
		Users: user.NewStorage(db),
	}, nil
}
