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
		db:  nil,
		ctx: context.Background(),
	}
}

func (s *Database) Connect(cfg *config.Database) (*Storage, error) {
	s.db = pg.Connect(&pg.Options{
		Addr:                  fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Database:              cfg.Schema,
		User:                  cfg.User,
		Password:              cfg.Password,
		MaxRetries:            10,
		RetryStatementTimeout: true,
	})

	if err := s.db.Ping(s.ctx); err != nil {
		return nil, err
	}

	return &Storage{
		db:    s.db,
		Users: user.NewStorage(s.db),
	}, nil
}

func (s *Database) Close() error {
	if s.db != nil {
		return s.db.Close()
	}

	return nil
}
