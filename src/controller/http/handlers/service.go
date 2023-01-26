package handlers

import (
	"test-project/src/controller/http/handlers/auth"
	"test-project/src/controller/http/handlers/users"
	"test-project/src/controller/repository"
	"test-project/src/pkg/jwt"

	"github.com/allegro/bigcache/v3"
	"github.com/labstack/echo/v4"
)

func NewHandlers(router *echo.Echo, jwt jwt.Servicer, cache *bigcache.BigCache, storage *repository.Storage) *Handlers {
	return &Handlers{
		Users: users.NewHandler(router, cache, storage),
		Auth:  auth.NewHandler(router, jwt, cache, storage),
	}
}
