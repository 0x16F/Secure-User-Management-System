package users

import (
	"net/http"
	"test-project/src/pkg/utils"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CheckPermissions(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		permissions := c.Request().Header.Get("X-User-Permissions")
		adminOnly := []string{http.MethodDelete, http.MethodPatch, http.MethodPost}

		if permissions != "admin" && utils.Contains[string](adminOnly, c.Request().Method) {
			return c.JSON(http.StatusForbidden, echo.Map{
				"message": "you don't have enough permissions",
			})
		}

		return next(c)
	}
}
