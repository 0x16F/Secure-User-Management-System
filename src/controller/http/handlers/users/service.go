package users

import (
	"net/http"
	"strconv"
	"test-project/src/controller/repository"
	"test-project/src/internal/user"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
)

func NewHandler(router *echo.Echo, storage *repository.Storage) IHandler {
	return &Handler{
		Router:  router,
		Storage: storage,
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
			"message": "bad request",
		})
	}

	request.Id = id

	if err := h.Storage.Users.Update(&request); err != nil {
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

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	if err := h.Storage.Users.Delete(id); err != nil {
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
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	u := user.NewUser(&request)

	id, err := h.Storage.Users.Create(u)
	if err != nil {
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

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &users)
}
