package http

import (
	"test-project/src/controller/http/handlers"

	"github.com/labstack/echo/v4"
)

type Server struct {
	Router   *echo.Echo
	Handlers *handlers.Handlers
}
