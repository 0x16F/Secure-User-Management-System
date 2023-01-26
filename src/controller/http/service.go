package http

import (
	"fmt"
	"os"
	"test-project/src/controller/http/handlers"
	"test-project/src/controller/repository"
	"test-project/src/pkg/jwt"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewServer(storage *repository.Storage, cache *bigcache.BigCache, jwt jwt.Servicer) *Server {
	router := echo.New()

	server := &Server{
		Router:   router,
		Handlers: handlers.NewHandlers(router, jwt, cache, storage),
	}

	server.configureRouters()

	return server
}

func (s *Server) configureRouters() {
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
	// logger middleware
	f, err := os.OpenFile(fmt.Sprintf("logs/%s.log", time.Now().Format("2006-01-02 15-04-05")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer f.Close()

	s.Router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
		Output: f,
	}))

	return s.Router.Start(fmt.Sprintf(":%d", port))
}
