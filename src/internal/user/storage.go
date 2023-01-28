package user

import (
	"fmt"
	"test-project/src/pkg/validate"
)

func (s *Storage) FindOne(id int64) (*FindUserDTO, error) {
	user := &FindUserDTO{}
	err := s.db.Model(user).Where("id = ?", id).Select()
	return user, err
}

func (s *Storage) FindByLogin(login string) (*User, error) {
	user := &User{}
	err := s.db.Model(user).Where("lower(login) = lower(?)", login).Select()
	return user, err
}

func (s *Storage) FindAll(limit, offset int, order string, filters *FindUsersFilters) (*[]FindUserDTO, int, error) {
	if order == "" {
		order = validate.OrderAsc
	}

	users := make([]FindUserDTO, 0)
	query := s.db.Model(&users).Limit(limit).Offset(offset).Order(fmt.Sprintf("id %s", order))

	// Name        *string  `json:"name,omitempty" default:"Иван"`
	// Surname     *string  `json:"surname,omitempty" default:"Иванов"`
	// Login       *string  `json:"login,omitempty" default:"Ivanov"`
	// Permissions *string  `json:"permissions,omitempty" enums:"admin,read-only,banned"`
	// Birthday    *int64   `json:"birthday,omitempty"`

	if filters != nil {
		if *filters != (FindUsersFilters{}) {
			if filters.Name != "" {
				query.Where("name LIKE ?", filters.Name+"%")
			}

			if filters.Surname != "" {
				query.Where("surname LIKE ?", filters.Surname+"%")
			}

			if filters.Login != "" {
				query.Where("login LIKE ?", filters.Login+"%")
			}

			if filters.Permissions != "" {
				query.Where("permissions LIKE ?", filters.Permissions+"%")
			}
		}
	}

	count, err := query.SelectAndCount()
	return &users, count, err
}

func (s *Storage) Delete(id int64) error {
	user := User{}
	_, err := s.db.Model(&user).Where("id = ?", id).Delete()
	return err
}

func (s *Storage) Create(dto *UserDTO) (*int64, error) {
	u := User{}

	// выглядит не очень, но мне нужно возвращать id пользователя, который был создан
	// возможно есть какой-то более элегантный способ сделать это на ORM'ке, но я его не знаю

	u.Name = dto.Name
	u.Surname = dto.Surname
	u.Login = dto.Login
	u.Password = dto.Password
	u.Salt = dto.Salt
	u.Permissions = dto.Permissions
	u.Birthday = dto.Birthday

	if _, err := s.db.Model(&u).Insert(); err != nil {
		return nil, err
	}

	return &u.Id, nil
}

func (s *Storage) Update(id int64, dto *UpdateUserDTO) error {
	_, err := s.db.Model(dto).Where("id = ?", id).UpdateNotZero()
	return err
}
