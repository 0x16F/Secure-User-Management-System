package auth

import (
	"fmt"
	"net/http"
	"test-project/src/controller/http/response"
	"test-project/src/controller/repository"
	"test-project/src/internal/permissions"
	"test-project/src/pkg/jwt"
	"test-project/src/pkg/utils"
	"test-project/src/pkg/validate"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
)

func NewHandler(router *echo.Echo, jwt jwt.Servicer, cache *bigcache.BigCache, storage *repository.Storage) IHandler {
	return &Handler{
		Router:  router,
		Storage: storage,
		JWT:     jwt,
		Cache:   cache,
	}
}

// @Summary login
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body LoginRequest true "credentials"
// @Success 200 {object} AccessResponse
// @Failure 400,403,404 {object} response.AppError
// @Failure 500 {object} response.AppError
// @Failure default {object} response.AppError
// @Router /auth/login [post]
func (h *Handler) Login(c echo.Context) error {
	request := LoginRequest{}

	if err := c.Bind(&request); err != nil {
		h.Router.Logger.Error(err)
		return response.BadRequestError("Failed to parse JSON", err.Error()).Send(c)
	}

	// is login valid?
	if !validate.Login(request.Login) {
		validationErr := response.BadRequestError("Incorrect login", "You can use only latin symbols and arabian numbers")
		return validationErr.Send(c)
	}

	// is login length valid?
	if !validate.LoginLenght(request.Login) {
		developerMessage := fmt.Sprintf("Minimum: %d, maximum: %d", validate.MinLoginLength, validate.MaxLoginLength)
		validationErr := response.BadRequestError("Incorrect login length", developerMessage)
		return validationErr.Send(c)
	}

	// is password length valid?
	if !validate.PasswordLength(request.Password) {
		developerMessage := fmt.Sprintf("Minimum: %d, maximum: %d", validate.MinPasswordLength, validate.MaxPasswordLength)
		validationErr := response.BadRequestError("Incorrect password length", developerMessage)
		return validationErr.Send(c)
	}

	user, err := h.Storage.Users.FindByLogin(request.Login)
	if err != nil {
		if err == pg.ErrNoRows {
			notFoundError := response.NewAppError(http.StatusNotFound, "User not found", "")
			return notFoundError.Send(c)
		}

		h.Router.Logger.Error(err)
		systemError := response.SystemError("Internal error, try again later", err.Error())
		return systemError.Send(c)
	}

	if hash, _ := utils.HashString(request.Password, user.Salt); hash != user.Password {
		return response.NewAppError(http.StatusForbidden, "Incorrect login or password", "").Send(c)
	}

	refresh, err := h.JWT.GenerateRefresh(&jwt.GenerateDTO{
		Id:          user.Id,
		Permissions: user.Permissions,
	})

	if err != nil {
		h.Router.Logger.Error(err)
		systemError := response.SystemError("Internal error, try again later", err.Error())
		return systemError.Send(c)
	}

	access, err := h.JWT.GenerateAccess(&jwt.GenerateDTO{
		Id:          user.Id,
		Permissions: user.Permissions,
	})

	if err != nil {
		h.Router.Logger.Error(err)
		systemError := response.SystemError("Internal error, try again later", err.Error())
		return systemError.Send(c)
	}

	c.SetCookie(&http.Cookie{
		Name:     "refresh",
		Value:    refresh,
		MaxAge:   30 * 24 * 60 * 60,
		Expires:  time.Now().Add(time.Hour * 720),
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, AccessResponse{
		Access: access,
	})
}

// @Summary refresh
// @Tags auth
// @Description refresh jwt access token
// @ID refresh
// @Produce  json
// @Success 200 {object} AccessResponse
// @Failure 403 {object} response.AppError
// @Failure 500 {object} response.AppError
// @Failure default {object} response.AppError
// @Router /auth/refresh [post]
func (h *Handler) Refresh(c echo.Context) error {
	cookie, err := c.Cookie("refresh")
	if err != nil {
		return response.NewAppError(http.StatusForbidden, "You need to login", `"refresh" is not found in cookie`).Send(c)
	}

	token, err := h.JWT.ParseRefresh(cookie.Value)
	if err != nil {
		if err == jwt.ErrExpired {
			return response.NewAppError(http.StatusForbidden, "Token is expired", err.Error()).Send(c)
		}

		return response.NewAppError(http.StatusForbidden, "Refresh token is not valid", err.Error()).Send(c)
	}

	user, err := h.Storage.Users.FindOne(token.Id)
	if err != nil {
		if err == pg.ErrNoRows {
			notFoundError := response.NewAppError(http.StatusNotFound, "User not found", "")
			return notFoundError.Send(c)
		}

		h.Router.Logger.Error(err)
		systemError := response.SystemError("Internal error, try again later", err.Error())
		return systemError.Send(c)
	}

	// is user banned?
	if user.Permissions == permissions.BannedPermission {
		forbiddenErr := response.NewAppError(http.StatusForbidden, "You are banned", "")
		return forbiddenErr.Send(c)
	}

	refresh, err := h.JWT.GenerateRefresh(&jwt.GenerateDTO{
		Id:          user.Id,
		Permissions: user.Permissions,
	})

	if err != nil {
		h.Router.Logger.Error(err)
		systemError := response.SystemError("Internal error, try again later", err.Error())
		return systemError.Send(c)
	}

	access, err := h.JWT.GenerateAccess(&jwt.GenerateDTO{
		Id:          user.Id,
		Permissions: user.Permissions,
	})

	if err != nil {
		h.Router.Logger.Error(err)
		systemError := response.SystemError("Internal error, try again later", err.Error())
		return systemError.Send(c)
	}

	c.SetCookie(&http.Cookie{
		Name:     "refresh",
		Value:    refresh,
		MaxAge:   30 * 24 * 60 * 60,
		Expires:  time.Now().Add(time.Hour * 720),
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, AccessResponse{
		Access: access,
	})
}
