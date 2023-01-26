package user

import "github.com/go-pg/pg/v10"

type User struct {
	tableName   struct{} `pg:"users"`
	Id          int64    `json:"id"`
	Name        string   `json:"name"`
	Surname     string   `json:"surname"`
	Login       string   `json:"login"`
	Password    string   `json:"password"`
	Salt        string   `json:"salt"`
	Permissions string   `json:"permissions"`
	Birthday    int64    `json:"birthday"`
}

type CreateUserDTO struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Permissions string `json:"permissions"`
	Birthday    int64  `json:"birthday"`
}

type FindUserDTO struct {
	tableName   struct{} `pg:"users"`
	Id          int64    `json:"id"`
	Name        string   `json:"name"`
	Surname     string   `json:"surname"`
	Login       string   `json:"login"`
	Permissions string   `json:"permissions"`
	Birthday    int64    `json:"birthday"`
}

type UserDTO struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Salt        string `json:"salt"`
	Permissions string `json:"permissions"`
	Birthday    int64  `json:"birthday"`
}

type Storage struct {
	db *pg.DB
}

type Storager interface {
	FindOne(id int64) (*FindUserDTO, error)
	FindByLogin(login string) (*User, error)
	FindAll(limit, offset int) (*[]FindUserDTO, error)
	Delete(id int64) error
	Create(dto *UserDTO) (*int64, error)
	Update(dto *UserDTO) error
}
