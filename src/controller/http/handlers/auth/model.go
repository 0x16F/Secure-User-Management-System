package auth

import (
	"test-project/src/controller/repository"
	"test-project/src/pkg/jwt"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Router  *echo.Echo
	Storage *repository.Storage
	JWT     jwt.Servicer
}

type IHandler interface {
	Login(c echo.Context) error
	Refresh(c echo.Context) error
	IsAuthorized(next echo.HandlerFunc) echo.HandlerFunc
}

type RequestLoginDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
