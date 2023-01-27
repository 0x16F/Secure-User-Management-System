package users

import (
	"net/http"
	"strconv"
	"test-project/src/controller/repository"
	"test-project/src/internal/permissions"
	"test-project/src/internal/user"
	"test-project/src/pkg/validate"

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

	request := user.UpdateUserDTO{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	// is permission valid?
	if request.Permissions != nil {
		if !validate.Permission(*request.Permissions) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "invalid permission",
			})
		}
	}

	if request.Login != nil {
		// is login valid?
		if !validate.Login(*request.Login) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "incorrect login, use only latin symbols and arabian numbers",
			})
		}

		// is login length valid?
		if !validate.LoginLenght(*request.Login) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "incorrect login length",
			})
		}
	}

	// is password length valid?
	if request.Password != nil {
		if !validate.PasswordLength(*request.Password) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "incorrect password length",
			})
		}
	}

	request.Id = id

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
		switch err {
		case pg.ErrNoRows:
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "user not found",
			})
		default:
			h.Router.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": err.Error(),
			})
		}
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
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "bad request",
		})
	}

	// is permission valid?
	if !validate.Permission(request.Permissions) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid permission",
		})
	}

	// is login valid?
	if !validate.Login(request.Login) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "incorrect login, use only latin symbols and arabian numbers",
		})
	}

	// is login length valid?
	if !validate.LoginLenght(request.Login) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "incorrect login length",
		})
	}

	// is password length valid?
	if !validate.PasswordLength(request.Password) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "incorrect password length",
		})
	}

	// is user already exists?
	if _, err := h.Storage.Users.FindByLogin(request.Login); err != nil {
		if err != pg.ErrNoRows {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": err.Error(),
			})
		}
	} else {
		return c.JSON(http.StatusConflict, echo.Map{
			"message": "user is already exists",
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

	return c.JSON(http.StatusCreated, echo.Map{
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
		switch err {
		case pg.ErrNoRows:
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "user not found",
			})
		default:
			h.Router.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, &user)
}

func (h *Handler) FindAll(c echo.Context) error {
	var limit, offset int
	var err error

	if limit, err = strconv.Atoi(c.QueryParam("limit")); err != nil {
		limit = 10
	}

	if offset, err = strconv.Atoi(c.QueryParam("offset")); err != nil {
		limit = 10
	}

	users, err := h.Storage.Users.FindAll(limit, offset)
	if err != nil {
		switch err {
		case pg.ErrNoRows:
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "users not found",
			})
		default:
			h.Router.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, &users)
}
