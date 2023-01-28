package auth

import (
	"bytes"
	"fmt"
	"net/http"
	"test-project/src/controller/http/response"
	"test-project/src/internal/permissions"
	headerparser "test-project/src/pkg/header-parser"
	"test-project/src/pkg/jwt"

	"github.com/allegro/bigcache/v3"
	"github.com/labstack/echo/v4"
)

func (h *Handler) IsAuthorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		access := headerparser.GetUserToken(c)

		if access == "" {
			return response.NewAppError(http.StatusForbidden, "You need to login", "").Send(c)
		}

		token, err := h.JWT.ParseAccess(access)
		if err != nil {
			if err == jwt.ErrExpired {
				return response.NewAppError(http.StatusForbidden, "Token is expired", err.Error()).Send(c)
			}

			return response.NewAppError(http.StatusForbidden, "Access token is not valid", err.Error()).Send(c)
		}

		result, err := h.Cache.Get(fmt.Sprint(token.Id))
		if err != nil {
			if err != bigcache.ErrEntryNotFound {
				h.Router.Logger.Error(err)
				systemError := response.SystemError("Internal error, try again later", err.Error())
				return systemError.Send(c)
			}
		}

		if bytes.Equal(result, []byte("deleted")) {
			c.SetCookie(&http.Cookie{
				Name:     "refresh",
				Value:    "",
				MaxAge:   -1,
				HttpOnly: true,
			})

			return response.NewAppError(http.StatusForbidden, "Account is not exists", "").Send(c)
		}

		if token.Permissions == permissions.BannedPermission || bytes.Equal(result, []byte(permissions.BannedPermission)) {
			return response.NewAppError(http.StatusForbidden, "You are banned", "").Send(c)
		}

		headerparser.SetUserId(c, fmt.Sprint(token.Id))
		headerparser.SetUserPermissions(c, token.Permissions)

		return next(c)
	}
}
