package users

import (
	"test-project/src/controller/repository"

	"github.com/allegro/bigcache/v3"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Router  *echo.Echo
	Storage *repository.Storage
	Cache   *bigcache.BigCache
}

type IHandler interface {
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Create(c echo.Context) error
	FindOne(c echo.Context) error
	FindAll(c echo.Context) error
	CheckPermissions(next echo.HandlerFunc) echo.HandlerFunc
}
