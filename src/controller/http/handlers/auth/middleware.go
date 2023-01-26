package auth

import (
	"bytes"
	"fmt"
	"net/http"
	"test-project/src/internal/permissions"
	headerparser "test-project/src/pkg/header-parser"

	"github.com/allegro/bigcache/v3"
	"github.com/labstack/echo/v4"
)

func (h *Handler) IsAuthorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		access := headerparser.GetUserToken(c)

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

		result, err := h.Cache.Get(fmt.Sprint(token.Id))
		if err != nil {
			if err != bigcache.ErrEntryNotFound {
				h.Router.Logger.Error(err)

				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": err.Error(),
				})
			}
		}

		if token.Permissions == permissions.BannedPermission || bytes.Equal(result, []byte(permissions.BannedPermission)) {
			return c.JSON(http.StatusForbidden, echo.Map{
				"message": "you're banned",
			})
		}

		headerparser.SetUserId(c, fmt.Sprint(token.Id))
		headerparser.SetUserPermissions(c, token.Permissions)

		return next(c)
	}
}
