package user

import "github.com/go-pg/pg/v10"

//go:generate mockgen -source=model.go -destination=mocks/mock.go

type User struct {
	tableName   struct{} `pg:"users"`
	Id          int64    `json:"id"`
	Name        string   `json:"name" default:"Иван"`
	Surname     string   `json:"surname" default:"Иванов"`
	Login       string   `json:"login" default:"Ivanov"`
	Password    string   `json:"password" default:"password"`
	Salt        string   `json:"salt"`
	Permissions string   `json:"permissions" enums:"admin,read-only,banned"`
	Birthday    int64    `json:"birthday"`
}

type CreateUserDTO struct {
	Name        string `json:"name" default:"Иван"`
	Surname     string `json:"surname" default:"Иванов"`
	Login       string `json:"login" default:"Ivanov"`
	Password    string `json:"password" default:"password"`
	Permissions string `json:"permissions" enums:"admin,read-only,banned"`
	Birthday    int64  `json:"birthday"`
}

type FindUserDTO struct {
	tableName   struct{} `pg:"users"`
	Id          int64    `json:"id"`
	Name        string   `json:"name" default:"Иван"`
	Surname     string   `json:"surname" default:"Иванов"`
	Login       string   `json:"login" default:"Ivanov"`
	Permissions string   `json:"permissions" enums:"admin,read-only,banned"`
	Birthday    int64    `json:"birthday"`
}

type UserDTO struct {
	Name        string `json:"name" default:"Иван"`
	Surname     string `json:"surname" default:"Иванов"`
	Login       string `json:"login" default:"Ivanov"`
	Password    string `json:"password" default:"password"`
	Salt        string `json:"salt"`
	Permissions string `json:"permissions" enums:"admin,read-only,banned"`
	Birthday    int64  `json:"birthday"`
}

type UpdateUserDTO struct {
	tableName   struct{} `pg:"users"`
	Name        *string  `json:"name,omitempty" default:"Иван"`
	Surname     *string  `json:"surname,omitempty" default:"Иванов"`
	Login       *string  `json:"login,omitempty" default:"Ivanov"`
	Password    *string  `json:"password,omitempty" default:"password"`
	Permissions *string  `json:"permissions,omitempty" enums:"admin,read-only,banned"`
	Birthday    *int64   `json:"birthday,omitempty"`
}

type FindUsersFilters struct {
	tableName   struct{} `pg:"users"`
	Name        *string  `json:"name,omitempty" default:"Иван"`
	Surname     *string  `json:"surname,omitempty" default:"Иванов"`
	Login       *string  `json:"login,omitempty" default:"Ivanov"`
	Permissions *string  `json:"permissions,omitempty" enums:"admin,read-only,banned"`
	Birthday    *int64   `json:"birthday,omitempty"`
}

type Storage struct {
	db *pg.DB
}

type Storager interface {
	FindOne(id int64) (*FindUserDTO, error)
	FindByLogin(login string) (*User, error)
	FindAll(limit, offset int, order string) (*[]FindUserDTO, int, error)
	Delete(id int64) error
	Create(dto *UserDTO) (*int64, error)
	Update(id int64, dto *UpdateUserDTO) error
}
