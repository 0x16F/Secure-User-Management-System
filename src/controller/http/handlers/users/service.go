package users

import (
	"net/http"
	"regexp"
	"strconv"
	"test-project/src/controller/repository"
	"test-project/src/internal/permissions"
	"test-project/src/internal/user"
	"test-project/src/pkg/utils"

	"github.com/allegro/bigcache/v3"
	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
)

func NewHandler(router *echo.Echo, cache *bigcache.BigCache, storage *repository.Storage) IHandler {
	return &Handler{
		Router:  router,
		Storage: storage,
		Cache:   cache,
	}
}

func (h *Handler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid id",
		})
	}

	// is user exists
	if _, err := h.Storage.Users.FindOne(id); err != nil {
		if err == pg.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "user not found",
			})
		}

		h.Router.Logger.Error(err)

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	request := user.UpdateUserDTO{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "bad request",
		})
	}

	request.Id = id

	// ban user
	if *request.Permissions == permissions.BannedPermission {
		h.Cache.Set(c.Param("id"), []byte(permissions.BannedPermission))
	}

	if err := h.Storage.Users.Update(&request); err != nil {
		h.Router.Logger.Error(err)

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "ok",
	})
}

func (h *Handler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid id",
		})
	}

	if _, err := h.Storage.Users.FindOne(id); err != nil {
		if err == pg.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "user not found",
			})
		}

		h.Router.Logger.Error(err)

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	if err := h.Storage.Users.Delete(id); err != nil {
		h.Router.Logger.Error(err)

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "ok",
	})
}

func (h *Handler) Create(c echo.Context) error {
	request := user.CreateUserDTO{}

	if err := c.Bind(&request); err != nil {
		h.Router.Logger.Error(err)

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	// is permission valid?
	if !utils.Contains(permissions.ArrayOfPermissions, request.Permissions) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid permission",
		})
	}

	// is login valid?
	if !regexp.MustCompile(`(?m)^[A-Za-z0-9]+$`).Match([]byte(request.Login)) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "incorrect login, use only latin symbols and arabian numbers",
		})
	}

	if len(request.Login) < 3 || len(request.Login) > 24 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "incorrect login length. min: 3, max: 24",
		})
	}

	// is password valid?
	passwordLength := len(request.Password)

	if passwordLength < 8 || passwordLength > 64 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "incorrect password length. min: 8, max: 64",
		})
	}

	u := user.NewUser(&request)

	id, err := h.Storage.Users.Create(u)
	if err != nil {
		h.Router.Logger.Error(err)

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id": id,
	})
}

func (h *Handler) FindOne(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Router.Logger.Error(err)

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	user, err := h.Storage.Users.FindOne(id)
	if err != nil {
		if err == pg.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "user not found",
			})
		}

		h.Router.Logger.Error(err)

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &user)
}

func (h *Handler) FindAll(c echo.Context) error {
	var err error

	limit := c.QueryParam("limit")
	limitInt := 10

	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "failed to parse limit",
			})
		}
	}

	offset := c.QueryParam("offset")
	offsetInt := 0

	if offset != "" {
		limitInt, err = strconv.Atoi(offset)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "failed to parse offset",
			})
		}
	}

	users, err := h.Storage.Users.FindAll(limitInt, offsetInt)
	if err != nil {
		if err == pg.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "users not found",
			})
		}

		h.Router.Logger.Error(err)

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &users)
}
