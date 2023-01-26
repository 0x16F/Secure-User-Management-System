package handlers

import (
	"test-project/src/controller/http/handlers/auth"
	"test-project/src/controller/http/handlers/users"
	"test-project/src/controller/repository"
	"test-project/src/pkg/jwt"

	"github.com/labstack/echo/v4"
)

func NewHandlers(router *echo.Echo, jwt jwt.Servicer, storage *repository.Storage) *Handlers {
	return &Handlers{
		Users: users.NewHandler(router, storage),
		Auth:  auth.NewHandler(router, jwt, storage),
	}
}
