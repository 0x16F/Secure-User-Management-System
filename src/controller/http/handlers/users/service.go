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

// @Summary update user
// @Security ApiKeyAuth
// @Tags users
// @Description update user
// @ID update user
// @Accept  json
// @Produce  json
// @Param id path integer true "user id"
// @Param input body user.UpdateUserDTO true "update info"
// @Success 200 {object} successResponse
// @Failure 400,403,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users/{id} [patch]
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
	if request.Permissions != nil {
		if *request.Permissions == permissions.BannedPermission {
			h.Cache.Set(c.Param("id"), []byte(permissions.BannedPermission))
		}
	}

	if err := h.Storage.Users.Update(id, &request); err != nil {
		h.Router.Logger.Error(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "ok",
	})
}

// @Summary delete user
// @Security ApiKeyAuth
// @Tags users
// @Description delete user
// @ID delete user
// @Produce  json
// @Param id path integer true "user id"
// @Success 200 {object} successResponse
// @Failure 400,403,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users/{id} [delete]
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

// @Summary create user
// @Security ApiKeyAuth
// @Tags users
// @Description create user
// @ID create user
// @Accept  json
// @Produce  json
// @Param input body user.CreateUserDTO true "create user info"
// @Success 201 {object} createResponse
// @Failure 400,403,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users [post]
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

// @Summary get user
// @Security ApiKeyAuth
// @Tags users
// @Description get user
// @ID get user
// @Produce  json
// @Param id path integer true "user id"
// @Success 200 {object} user.FindUserDTO
// @Failure 400,403,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users/{id} [get]
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

// @Summary get users
// @Security ApiKeyAuth
// @Tags users
// @Description get users
// @ID get users
// @Produce  json
// @Param limit query integer false "limit"
// @Param offset query integer false "offset"
// @Success 200 {object} []user.FindUserDTO
// @Failure 400,403,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users [get]
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
