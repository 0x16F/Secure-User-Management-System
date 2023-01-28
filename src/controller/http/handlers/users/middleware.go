package users

import (
	"net/http"
	"test-project/src/controller/http/response"
	"test-project/src/internal/permissions"
	headerparser "test-project/src/pkg/header-parser"
	"test-project/src/pkg/utils"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CheckPermissions(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		adminOnly := []string{http.MethodDelete, http.MethodPatch, http.MethodPost}
		userPermissions := headerparser.GetUserPermissions(c)

		if userPermissions != permissions.AdminPermission && utils.Contains(adminOnly, c.Request().Method) {
			return response.NewAppError(http.StatusForbidden, "You don't have enough permissions", "").Send(c)
		}

		return next(c)
	}
}
