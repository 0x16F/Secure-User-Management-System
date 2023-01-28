package http

import (
	"fmt"
	"test-project/src/controller/http/handlers"
	"test-project/src/controller/repository"
	"test-project/src/pkg/jwt"

	"github.com/allegro/bigcache/v3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "test-project/docs"
)

func NewServer(storage *repository.Storage, cache *bigcache.BigCache, jwt jwt.Servicer) *Server {
	router := echo.New()

	router.Logger.SetLevel(log.DEBUG)

	server := &Server{
		Router:   router,
		Handlers: handlers.NewHandlers(router, jwt, cache, storage),
	}

	server.configureRouters()

	return server
}

func (s *Server) configureRouters() {
	// logger middleware
	s.Router.Use(middleware.Logger())

	s.Router.GET("/swagger/*", echoSwagger.WrapHandler)

	auth := s.Router.Group("/auth")
	{
		auth.POST("/login", s.Handlers.Auth.Login)
		auth.POST("/refresh", s.Handlers.Auth.Refresh)
	}

	users := s.Router.Group("/users", s.Handlers.Auth.IsAuthorized, s.Handlers.Users.CheckPermissions)
	{
		users.GET("/:id", s.Handlers.Users.FindOne)
		users.GET("", s.Handlers.Users.FindAll)
		users.DELETE("/:id", s.Handlers.Users.Delete)
		users.POST("", s.Handlers.Users.Create)
		users.PATCH("/:id", s.Handlers.Users.Update)
	}
}

func (s *Server) Start(port uint16) error {
	return s.Router.Start(fmt.Sprintf(":%d", port))
}
