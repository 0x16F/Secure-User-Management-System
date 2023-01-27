package users

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"test-project/src/controller/repository"
	"test-project/src/internal/user"
	mock_user "test-project/src/internal/user/mocks"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/go-pg/pg/v10"
	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestHandler_FindOne(t *testing.T) {
	type MockCallback func(s *mock_user.MockStorager, InputParam int64, ExpectedUser *user.FindUserDTO)

	type Test struct {
		Name         string
		InputParam   int64
		ExpectedUser *user.FindUserDTO
		ExpectedCode int
		MockCallback MockCallback
	}

	testTable := []Test{
		{
			Name:       "OK",
			InputParam: 1,
			ExpectedUser: &user.FindUserDTO{
				Id:          1,
				Name:        "admin",
				Surname:     "admin",
				Login:       "admin",
				Permissions: "admin",
				Birthday:    0,
			},
			ExpectedCode: http.StatusOK,
			MockCallback: func(s *mock_user.MockStorager, InputParam int64, ExpectedUser *user.FindUserDTO) {
				s.EXPECT().FindOne(InputParam).Return(ExpectedUser, nil)
			},
		},
		{
			Name:         "Not found",
			InputParam:   1,
			ExpectedCode: http.StatusNotFound,
			MockCallback: func(s *mock_user.MockStorager, InputParam int64, ExpectedUser *user.FindUserDTO) {
				s.EXPECT().FindOne(InputParam).Return(nil, pg.ErrNoRows)
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
			testCase.MockCallback(userStorage, testCase.InputParam, testCase.ExpectedUser)

			handler := NewHandler(router, cache, &repository.Storage{Users: userStorage})

			router.GET("/users/:id", handler.FindOne)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("", fmt.Sprintf("/users/%d", testCase.InputParam), nil)

			router.ServeHTTP(w, req)

			if w.Code != testCase.ExpectedCode {
				t.Fatalf("Wrong status code, exptected: %d, got %d", testCase.ExpectedCode, w.Code)
			}

			user := user.FindUserDTO{}
			json.NewDecoder(w.Body).Decode(&user)

			if testCase.ExpectedUser != nil {
				if user != *testCase.ExpectedUser {
					t.Fatalf("%v\n%v", user, testCase.ExpectedUser)
				}
			}
		})
	}
}

func TestHandler_Update(t *testing.T) {
	type MockCallback func(s *mock_user.MockStorager, dto *user.UpdateUserDTO, userId int64)

	type Test struct {
		Name           string
		InputUserId    string
		InputBody      string
		ExpectedDTO    *user.UpdateUserDTO
		ExpectedUserId int64
		ExpectedCode   int
		MockCallback   MockCallback
	}

	banned := "banned"

	testTable := []Test{
		{
			Name:        "OK",
			InputUserId: "1",
			InputBody:   `{"permissions": "banned"}`,
			ExpectedDTO: &user.UpdateUserDTO{
				Id:          1,
				Permissions: &banned,
			},
			ExpectedUserId: 1,
			ExpectedCode:   http.StatusOK,
			MockCallback: func(s *mock_user.MockStorager, dto *user.UpdateUserDTO, userId int64) {
				s.EXPECT().FindOne(userId).Return(nil, nil)
				s.EXPECT().Update(dto).Return(nil)
			},
		},
		{
			Name:         "Invalid id",
			InputUserId:  "undefined",
			InputBody:    `{"permissions": "banned"}`,
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:           "User not found",
			InputUserId:    "1",
			InputBody:      `{"permissions": "banned"}`,
			ExpectedUserId: 1,
			ExpectedCode:   http.StatusNotFound,
			MockCallback: func(s *mock_user.MockStorager, dto *user.UpdateUserDTO, userId int64) {
				s.EXPECT().FindOne(userId).Return(nil, pg.ErrNoRows)
			},
		},
		{
			Name:         "Invalid permission",
			InputUserId:  "1",
			InputBody:    `{"permissions": "read"}`,
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:         "Invalid login length",
			InputUserId:  "1",
			InputBody:    `{"login": "r"}`,
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:         "Invalid login",
			InputUserId:  "1",
			InputBody:    `{"login": "Иван!"}`,
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:         "Invalid password length",
			InputUserId:  "1",
			InputBody:    `{"password": "pass"}`,
			ExpectedCode: http.StatusBadRequest,
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
				if testCase.ExpectedDTO != nil {
					testCase.MockCallback(userStorage, testCase.ExpectedDTO, testCase.ExpectedUserId)
				} else {
					dto := &user.UpdateUserDTO{}
					testCase.MockCallback(userStorage, dto, testCase.ExpectedUserId)
				}
			}

			handler := NewHandler(router, cache, &repository.Storage{Users: userStorage})

			router.PATCH("/users/:id", handler.Update)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/users/%s", testCase.InputUserId), bytes.NewBufferString(testCase.InputBody))
			req.Header.Set("Content-Type", "application/json; charset=utf8")

			router.ServeHTTP(w, req)

			if w.Code != testCase.ExpectedCode {
				t.Fatalf("Wrong status code, exptected: %d, got %d\n%s", testCase.ExpectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	type MockCallback func(s *mock_user.MockStorager, id int64)

	type Test struct {
		Name            string
		InputUserId     string
		ExptectedUserId int64
		ExpectedCode    int
		MockCallback    MockCallback
	}

	testTable := []Test{
		{
			Name:            "OK",
			InputUserId:     "1",
			ExptectedUserId: 1,
			ExpectedCode:    http.StatusOK,
			MockCallback: func(s *mock_user.MockStorager, id int64) {
				s.EXPECT().FindOne(id).Return(nil, nil)
				s.EXPECT().Delete(id).Return(nil)
			},
		},
		{
			Name:         "Invalid id",
			InputUserId:  "fsf",
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:            "User not found",
			InputUserId:     "1",
			ExptectedUserId: 1,
			ExpectedCode:    http.StatusNotFound,
			MockCallback: func(s *mock_user.MockStorager, id int64) {
				s.EXPECT().FindOne(id).Return(nil, pg.ErrNoRows)
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
				testCase.MockCallback(userStorage, testCase.ExptectedUserId)
			}

			handler := NewHandler(router, cache, &repository.Storage{Users: userStorage})

			router.DELETE("/users/:id", handler.Delete)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/users/%s", testCase.InputUserId), nil)

			router.ServeHTTP(w, req)

			if w.Code != testCase.ExpectedCode {
				t.Fatalf("Wrong status code, exptected: %d, got %d\n%s", testCase.ExpectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestHandler_Create(t *testing.T) {
	type MockCallback func(s *mock_user.MockStorager)

	type Test struct {
		Name         string
		InputBody    string
		ExpectedCode int
		MockCallback MockCallback
	}

	testTable := []Test{
		{
			Name:         "OK",
			InputBody:    `{"name": "Иван", "surname": "Иванов", "login": "ivanov123", "password": "ivanov123", "permissions": "read-only", "birthday": 696435862}`,
			ExpectedCode: http.StatusCreated,
			MockCallback: func(s *mock_user.MockStorager) {
				id := int64(1)

				s.EXPECT().FindByLogin("ivanov123").Return(nil, pg.ErrNoRows)
				s.EXPECT().Create(gomock.Any()).Return(&id, nil)
			},
		},
		{
			Name:         "user is already exists",
			InputBody:    `{"name": "Иван", "surname": "Иванов", "login": "ivanov123", "password": "ivanov123", "permissions": "read-only", "birthday": 696435862}`,
			ExpectedCode: http.StatusConflict,
			MockCallback: func(s *mock_user.MockStorager) {
				s.EXPECT().FindByLogin("ivanov123").Return(nil, nil)
			},
		},
		{
			Name:         "bad request",
			InputBody:    `{"name": "Иван", "surname": "Иванов", "password": "ivanov123", "permissions": "read-only", "birthday": 696435862}`,
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:         "invalid login length",
			InputBody:    `{"name": "Иван", "surname": "Иванов", "login": "iv", "password": "ivanov123", "permissions": "read-only", "birthday": 696435862}`,
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:         "invlaid login",
			InputBody:    `{"name": "Иван", "surname": "Иванов", "login": "Иванов", "password": "ivanov123", "permissions": "read-only", "birthday": 696435862}`,
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:         "invalid password length",
			InputBody:    `{"name": "Иван", "surname": "Иванов", "login": "ivanov123", "password": "iv", "permissions": "read-only", "birthday": 696435862}`,
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:         "invalid permission",
			InputBody:    `{"name": "Иван", "surname": "Иванов", "login": "ivanov123", "password": "ivanov123", "permissions": "read", "birthday": 696435862}`,
			ExpectedCode: http.StatusBadRequest,
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
				testCase.MockCallback(userStorage)
			}

			handler := NewHandler(router, cache, &repository.Storage{Users: userStorage})
			router.POST("/users", handler.Create)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users", bytes.NewBufferString(testCase.InputBody))
			req.Header.Set("Content-Type", "application/json; charset=utf8")

			router.ServeHTTP(w, req)

			if w.Code != testCase.ExpectedCode {
				t.Fatalf("Wrong status code, exptected: %d, got %d", testCase.ExpectedCode, w.Code)
			}
		})
	}
}

func TestHandler_FindAll(t *testing.T) {
	type MockCallback func(s *mock_user.MockStorager, limit, offset int, ExpectedUsers *[]user.FindUserDTO)

	type Test struct {
		Name          string
		InputLimit    int
		InputOffset   int
		ExpectedUsers *[]user.FindUserDTO
		ExpectedCode  int
		MockCallback  MockCallback
	}

	testTable := []Test{
		{
			Name:        "OK",
			InputLimit:  10,
			InputOffset: 0,
			ExpectedUsers: &[]user.FindUserDTO{
				{
					Id:          1,
					Name:        "admin",
					Surname:     "admin",
					Login:       "admin",
					Permissions: "admin",
					Birthday:    0,
				},
			},
			ExpectedCode: http.StatusOK,
			MockCallback: func(s *mock_user.MockStorager, limit, offset int, ExpectedUsers *[]user.FindUserDTO) {
				s.EXPECT().FindAll(limit, offset).Return(ExpectedUsers, nil)
			},
		},
		{
			Name:         "Not found",
			InputLimit:   10,
			InputOffset:  5,
			ExpectedCode: http.StatusNotFound,
			MockCallback: func(s *mock_user.MockStorager, limit, offset int, ExpectedUsers *[]user.FindUserDTO) {
				s.EXPECT().FindAll(limit, offset).Return(nil, pg.ErrNoRows)
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
			testCase.MockCallback(userStorage, testCase.InputLimit, testCase.InputOffset, testCase.ExpectedUsers)

			handler := NewHandler(router, cache, &repository.Storage{Users: userStorage})

			router.GET("/users", handler.FindAll)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("", fmt.Sprintf("/users?limit=%d&offset=%d", testCase.InputLimit, testCase.InputOffset), nil)

			router.ServeHTTP(w, req)

			if w.Code != testCase.ExpectedCode {
				t.Fatalf("Wrong status code, exptected: %d, got %d", testCase.ExpectedCode, w.Code)
			}

			users := make([]*user.FindUserDTO, 0)
			json.NewDecoder(w.Body).Decode(&users)

			if testCase.ExpectedUsers != nil {
				for index, exptectedUser := range *testCase.ExpectedUsers {
					if *users[index] != exptectedUser {
						t.Fatalf("Expected: %v\nGet: %v", *users[index], exptectedUser)
					}
				}
			}
		})
	}
}
