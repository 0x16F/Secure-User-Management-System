package users

import (
	"fmt"
	"net/http"
	"strconv"
	"test-project/src/controller/http/response"
	"test-project/src/controller/repository"
	"test-project/src/internal/permissions"
	"test-project/src/internal/user"
	"test-project/src/pkg/validate"
	"time"

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
// @Success 200 {object} response.SuccessResponse
// @Failure 400,403,404 {object} response.AppError
// @Failure 500 {object} response.AppError
// @Router /users/{id} [patch]
func (h *Handler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		validationErr := response.BadRequestError("Bad Request", "Incorrect id")
		validationErr.WithParams(response.Map{
			"id": "this field should be int64",
		})
		return validationErr.Send(c)
	}

	request := user.UpdateUserDTO{}
	if err := c.Bind(&request); err != nil {
		h.Router.Logger.Error(err)
		return response.BadRequestError("Failed to parse JSON", err.Error()).Send(c)
	}

	// is request empty?
	if request == (user.UpdateUserDTO{}) {
		return response.BadRequestError("Empty request", "All params in request is nil").Send(c)
	}

	// is permission valid?
	if request.Permissions != nil {
		if !validate.Permission(*request.Permissions) {
			return response.BadRequestError("Incorrect permission", "This permission is not exists").Send(c)
		}
	}

	if request.Login != nil {
		// is login valid?
		if !validate.Login(*request.Login) {
			validationErr := response.BadRequestError("Incorrect login", "You can use only latin symbols and arabian numbers")
			return validationErr.Send(c)
		}

		// is login length valid?
		if !validate.LoginLenght(*request.Login) {
			developerMessage := fmt.Sprintf("Minimum: %d, maximum: %d", validate.MinLoginLength, validate.MaxLoginLength)
			validationErr := response.BadRequestError("Incorrect login length", developerMessage)
			return validationErr.Send(c)
		}
	}

	// is password length valid?
	if request.Password != nil {
		if !validate.PasswordLength(*request.Password) {
			developerMessage := fmt.Sprintf("Minimum: %d, maximum: %d", validate.MinPasswordLength, validate.MaxPasswordLength)
			validationErr := response.BadRequestError("Incorrect password length", developerMessage)
			return validationErr.Send(c)
		}
	}

	// is user exists
	if _, err := h.Storage.Users.FindOne(id); err != nil {
		if err == pg.ErrNoRows {
			notFoundError := response.NewAppError(http.StatusNotFound, "User not found", "")
			return notFoundError.Send(c)
		}

		h.Router.Logger.Error(err)
		systemError := response.SystemError("Internal error, try again later", err.Error())
		return systemError.Send(c)
	}

	// ban user
	if request.Permissions != nil {
		if *request.Permissions == permissions.BannedPermission {
			h.Cache.Set(c.Param("id"), []byte(permissions.BannedPermission))
		}
	}

	if err := h.Storage.Users.Update(id, &request); err != nil {
		h.Router.Logger.Error(err)
		systemError := response.SystemError("Internal error, try again later", err.Error())
		return systemError.Send(c)
	}

	success := response.Success(http.StatusOK, "OK")
	return success.Send(c)
}

// @Summary delete user
// @Security ApiKeyAuth
// @Tags users
// @Description delete user
// @ID delete user
// @Produce  json
// @Param id path integer true "user id"
// @Success 200 {object} response.SuccessResponse
// @Failure 400,403,404 {object} response.AppError
// @Failure 500 {object} response.AppError
// @Router /users/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		validationErr := response.BadRequestError("Bad Request", "Incorrect id")
		validationErr.WithParams(response.Map{
			"id": "this field should be int64",
		})
		return validationErr.Send(c)
	}

	if _, err := h.Storage.Users.FindOne(id); err != nil {
		if err == pg.ErrNoRows {
			notFoundError := response.NewAppError(http.StatusNotFound, "User not found", "")
			return notFoundError.Send(c)
		}

		h.Router.Logger.Error(err)
		systemError := response.SystemError("Internal error, try again later", err.Error())
		return systemError.Send(c)
	}

	if err := h.Storage.Users.Delete(id); err != nil {
		h.Router.Logger.Error(err)
		systemError := response.SystemError("Internal error, try again later", err.Error())
		return systemError.Send(c)
	}

	h.Cache.Append(fmt.Sprint(id), []byte("deleted"))

	success := response.Success(http.StatusOK, "OK")
	return success.Send(c)
}

// @Summary create user
// @Security ApiKeyAuth
// @Tags users
// @Description create user
// @ID create user
// @Accept  json
// @Produce  json
// @Param input body user.CreateUserDTO true "create user info"
// @Success 201 {object} CreateUserResponse
// @Failure 400,403,404 {object} response.AppError
// @Failure 500 {object} response.AppError
// @Router /users [post]
func (h *Handler) Create(c echo.Context) error {
	request := user.CreateUserDTO{}
	if err := c.Bind(&request); err != nil {
		h.Router.Logger.Error(err)
		return response.BadRequestError("Failed to parse JSON", err.Error()).Send(c)
	}

	// is permission valid?
	if !validate.Permission(request.Permissions) {
		return response.BadRequestError("Incorrect permission", "This permission is not exists").Send(c)
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

	// is user already exists?
	if _, err := h.Storage.Users.FindByLogin(request.Login); err != nil {
		if err != pg.ErrNoRows {
			h.Router.Logger.Error(err)
			systemError := response.SystemError("Internal error, try again later", err.Error())
			return systemError.Send(c)
		}
	} else {
		systemError := response.NewAppError(http.StatusConflict, "This user is already exists", "")
		return systemError.Send(c)
	}

	u := user.NewUser(&request)

	id, err := h.Storage.Users.Create(u)
	if err != nil {
		h.Router.Logger.Error(err)
		systemError := response.SystemError("Internal error, try again later", err.Error())
		return systemError.Send(c)
	}

	return c.JSON(http.StatusCreated, CreateUserResponse{
		Id: *id,
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
// @Failure 400,403,404 {object} response.AppError
// @Failure 500 {object} response.AppError
// @Router /users/{id} [get]
func (h *Handler) FindOne(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		validationErr := response.BadRequestError("Bad Request", "Incorrect id")
		validationErr.WithParams(response.Map{
			"id": "this field should be int64",
		})
		return validationErr.Send(c)
	}

	user, err := h.Storage.Users.FindOne(id)
	if err != nil {
		if err == pg.ErrNoRows {
			notFoundError := response.NewAppError(http.StatusNotFound, "User not found", "")
			return notFoundError.Send(c)
		}

		h.Router.Logger.Error(err)
		systemError := response.SystemError("Internal error, try again later", err.Error())
		return systemError.Send(c)
	}

	return c.JSON(http.StatusOK, &user)
}

// @Summary get users
// @Security ApiKeyAuth
// @Tags users
// @Description get users
// @ID get users
// @Produce  json
// @Param limit query integer false "limit" default(50)
// @Param order query string false "order" Enums(asc, desc) default(asc)
// @Param name query string false "name" default()
// @Param surname query string false "surname" default()
// @Param login query string false "login" default()
// @Param permissions query string false "permissions" Enums(read-only, banned, admin) default()
// @Param birthday query string false "birthday" default() Format(date)
// @Param offset query integer false "offset" default(0)
// @Success 200 {object} FindUsersResponse
// @Failure 400,403,404 {object} response.AppError
// @Failure 500 {object} response.AppError
// @Router /users [get]
func (h *Handler) FindAll(c echo.Context) error {
	limit, offset := validate.MaxLimit, 0
	var err error

	// is limit valid?
	if limit, err = strconv.Atoi(c.QueryParam("limit")); err != nil {
		limit = validate.MaxLimit
	}

	// is limit in limit range?
	if !validate.Limit(limit) {
		limit = validate.MaxLimit
	}

	// is offset valid?
	if offset, err = strconv.Atoi(c.QueryParam("offset")); err != nil {
		offset = 0
	}

	order := c.QueryParam("order")

	// is order valid?
	if !validate.Order(order) {
		order = validate.OrderAsc
	}

	filters := &user.FindUsersFilters{
		Name:        c.QueryParam("name"),
		Surname:     c.QueryParam("surname"),
		Login:       c.QueryParam("login"),
		Permissions: c.QueryParam("permissions"),
		Birthday:    c.QueryParam("birthday"),
	}

	// is birthday valid?
	if filters.Birthday != "" {
		if _, err := time.Parse("2006-01-02", filters.Birthday); err != nil {
			validationErr := response.BadRequestError("Bad request", "Incorrect birthday format")
			validationErr.WithParams(response.Map{
				"birthday": "this field should be date, like a '2006-02-16'",
			})
			return validationErr.Send(c)
		}
	}

	users, count, err := h.Storage.Users.FindAll(limit, offset, order, filters)
	if err != nil {
		if err == pg.ErrNoRows {
			notFoundError := response.NewAppError(http.StatusNotFound, "Users not found", "")
			return notFoundError.Send(c)
		}

		h.Router.Logger.Error(err)
		systemError := response.SystemError("Internal error, try again later", err.Error())
		return systemError.Send(c)
	}

	return c.JSON(http.StatusOK, FindUsersResponse{
		Count: count,
		Users: users,
	})
}
