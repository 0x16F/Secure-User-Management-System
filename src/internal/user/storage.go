package user

import (
	"fmt"
	"test-project/src/pkg/utils"

	"github.com/go-pg/pg/v10"
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
	users := make([]FindUserDTO, 0)
	query := s.db.Model(&users).Limit(limit).Offset(offset)

	if order != "" {
		query.Order(fmt.Sprintf("id %s", order))
	}

	if filters != nil {
		if *filters != (FindUsersFilters{}) {
			if filters.Name != "" {
				query.Where("lower(name) LIKE lower(?)", filters.Name+"%")
			}

			if filters.Surname != "" {
				query.Where("lower(surname) LIKE lower(?)", filters.Surname+"%")
			}

			if filters.Login != "" {
				query.Where("lower(login) LIKE lower(?)", filters.Login+"%")
			}

			if filters.Permissions != "" {
				query.Where("lower(permissions) LIKE lower(?)", filters.Permissions+"%")
			}

			if filters.Birthday != "" {
				query.Where("birthday = ?", filters.Birthday)
			}
		}
	}

	count, err := query.SelectAndCount()
	return &users, count, err
}

func (s *Storage) Delete(id int64) error {
	user := User{}

	result, err := s.db.Model(&user).Where("id = ?", id).Delete()
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pg.ErrNoRows
	}

	return nil
}

func (s *Storage) Create(dto *UserDTO) (*int64, error) {
	u, err := utils.TypeConverter[User](&dto)
	if err != nil {
		return nil, err
	}

	if _, err := s.db.Model(u).Insert(); err != nil {
		return nil, err
	}

	return &u.Id, nil
}

func (s *Storage) Update(id int64, dto *UpdateUserDTO) error {
	u, err := utils.TypeConverter[User](&dto)
	if err != nil {
		return err
	}

	if u.Password != "" {
		u.Salt = utils.GenerateString(SaltLength)
		u.Password, _ = utils.HashString(u.Password, u.Salt)
	}

	result, err := s.db.Model(u).Where("id = ?", id).UpdateNotZero()
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pg.ErrNoRows
	}

	return nil
}
