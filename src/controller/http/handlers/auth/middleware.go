package auth

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) IsAuthorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		access := c.Request().Header.Get("Authorization")

		if access == "" {
			return c.JSON(http.StatusForbidden, echo.Map{
				"message": "you need to login",
			})
		}

		token, err := h.JWT.ParseAccess(access)
		if err != nil {
			return c.JSON(http.StatusForbidden, echo.Map{
				"message": err.Error(),
			})
		}

		c.Request().Header.Set("X-User-Id", fmt.Sprint(token.Id))
		c.Request().Header.Set("X-User-Permissions", token.Permissions)

		return next(c)
	}
}

func (h *Handler) IsBanned(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		permissions := c.Request().Header.Get("X-User-Permissions")

		if permissions == "banned" {
			return c.JSON(http.StatusForbidden, echo.Map{
				"message": "you're banned",
			})
		}

		return next(c)
	}
}
