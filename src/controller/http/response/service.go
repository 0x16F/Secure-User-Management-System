package response

import (
	"net/http"

	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
)

func NewAppError(code int, message, developerMessage string) *AppError {
	return &AppError{
		code:             code,
		Message:          message,
		DeveloperMessage: developerMessage,
	}
}

func BadRequestError(message, developerMessage string) *AppError {
	return NewAppError(http.StatusBadRequest, message, developerMessage)
}

func SystemError(message, developerMessage string) *AppError {
	return NewAppError(http.StatusInternalServerError, message, developerMessage)
}

func Success(code int, message string) *SuccessResponse {
	return &SuccessResponse{
		code:    code,
		Message: message,
	}
}

func (e *AppError) WithParams(params Map) {
	e.Params = params
}

func (r *SuccessResponse) WithParams(params Map) {
	r.Params = params
}

func marshal(i interface{}) []byte {
	byteData, err := json.Marshal(i)
	if err != nil {
		return nil
	}

	return byteData
}

func (e *AppError) Send(c echo.Context) error {
	return c.String(e.code, string(marshal(e)))
}

func (r *SuccessResponse) Send(c echo.Context) error {
	return c.String(r.code, string(marshal(r)))
}
