package auth

import (
	"net/http"
	"test-project/src/controller/repository"
	"test-project/src/pkg/jwt"
	"test-project/src/pkg/utils"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
)

func NewHandler(router *echo.Echo, jwt jwt.Servicer, storage *repository.Storage) IHandler {
	return &Handler{
		Router:  router,
		Storage: storage,
		JWT:     jwt,
	}
}

func (h *Handler) Login(c echo.Context) error {
	request := RequestLoginDTO{}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	user, err := h.Storage.Users.FindByLogin(request.Login)
	if err != nil {
		if err == pg.ErrNoRows {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "invalid login or password",
			})
		}

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	if hash, _ := utils.HashString(request.Password, user.Salt); hash != user.Password {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid login or password",
		})
	}

	refresh, err := h.JWT.GenerateRefresh(&jwt.GenerateDTO{
		Id:          user.Id,
		Permissions: user.Permissions,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	access, err := h.JWT.GenerateAccess(&jwt.GenerateDTO{
		Id:          user.Id,
		Permissions: user.Permissions,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	c.SetCookie(&http.Cookie{
		Name:     "refresh",
		Value:    refresh,
		MaxAge:   30 * 24 * 60 * 60,
		Expires:  time.Now().Add(time.Hour * 720),
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, echo.Map{
		"access": access,
	})
}

func (h *Handler) Refresh(c echo.Context) error {
	cookie, err := c.Cookie("refresh")
	if err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "refresh is not found",
		})
	}

	token, err := h.JWT.ParseRefresh(cookie.Value)
	if err != nil {
		if err == jwt.ErrExpired {
			return c.JSON(http.StatusForbidden, echo.Map{
				"message": "refresh is expired",
			})
		}

		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "refresh is not valid",
		})
	}

	refresh, err := h.JWT.GenerateRefresh(&jwt.GenerateDTO{
		Id:          token.Id,
		Permissions: token.Permissions,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	c.SetCookie(&http.Cookie{
		Name:     "refresh",
		Value:    refresh,
		MaxAge:   30 * 24 * 60 * 60,
		Expires:  time.Now().Add(time.Hour * 720),
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, echo.Map{
		"message": "ok",
	})
}
