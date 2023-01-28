package auth

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"test-project/src/controller/repository"
	"test-project/src/internal/user"
	mock_user "test-project/src/internal/user/mocks"
	"test-project/src/pkg/jwt"
	"test-project/src/pkg/utils"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestHandler_Login(t *testing.T) {
	type MockCallback func(s *mock_user.MockStorager, login, password string)

	type Test struct {
		Name             string
		InputBody        string
		ExpectedLogin    string
		ExpectedPassword string
		ExpectedCode     int
		MockCallback     MockCallback
	}

	testTable := []Test{
		{
			Name:             "OK",
			InputBody:        `{"login": "admin", "password": "password"}`,
			ExpectedLogin:    "admin",
			ExpectedPassword: "password",
			ExpectedCode:     http.StatusOK,
			MockCallback: func(s *mock_user.MockStorager, login, password string) {
				salt := "saltsaltsalt"
				hash, _ := utils.HashString(password, salt)

				s.EXPECT().FindByLogin(login).Return(&user.User{
					Id:          1,
					Name:        "admin",
					Surname:     "admin",
					Login:       login,
					Password:    hash,
					Salt:        salt,
					Permissions: "admin",
					Birthday:    0,
				}, nil)
			},
		},
		{
			Name:             "Bad password",
			InputBody:        `{"login": "admin", "password": "password"}`,
			ExpectedLogin:    "admin",
			ExpectedPassword: "admin12345",
			ExpectedCode:     http.StatusForbidden,
			MockCallback: func(s *mock_user.MockStorager, login, password string) {
				salt := "saltsaltsalt"
				hash, _ := utils.HashString(password, salt)

				s.EXPECT().FindByLogin(login).Return(&user.User{
					Id:          1,
					Name:        "admin",
					Surname:     "admin",
					Login:       login,
					Password:    hash,
					Salt:        salt,
					Permissions: "admin",
					Birthday:    0,
				}, nil)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			router := echo.New()
			cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))

			userStorage := mock_user.NewMockStorager(c)

			if testCase.MockCallback != nil {
				testCase.MockCallback(userStorage, testCase.ExpectedLogin, testCase.ExpectedPassword)
			}

			jwtService := jwt.NewService(&jwt.Service{
				AccessSecret:  "secret",
				RefreshSecret: "secret",
			})

			handler := NewHandler(router, jwtService, cache, &repository.Storage{Users: userStorage})
			router.POST("/auth/login", handler.Login)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(testCase.InputBody))
			req.Header.Set("Content-Type", "application/json; charset=utf8")

			router.ServeHTTP(w, req)

			if w.Code != testCase.ExpectedCode {
				t.Fatalf("Wrong status code, exptected: %d, got %d", testCase.ExpectedCode, w.Code)
			}
		})
	}
}
