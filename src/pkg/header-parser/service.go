package headerparser

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func SetUserId(c echo.Context, id string) {
	c.Request().Header.Set(headerUserId, id)
}

func SetUserPermissions(c echo.Context, permissions string) {
	c.Request().Header.Set(headerUserPermissions, permissions)
}

func GetUserId(c echo.Context) (int64, error) {
	return strconv.ParseInt(c.Request().Header.Get(headerUserId), 10, 64)
}

func GetUserPermissions(c echo.Context) string {
	return c.Request().Header.Get(headerUserPermissions)
}

func GetUserToken(c echo.Context) string {
	return c.Request().Header.Get(headerUserToken)
}
