package user

func (s *Storage) FindOne(id int64) (*FindUserDTO, error) {
	user := &FindUserDTO{}
	err := s.db.Model(user).Where("id = ?", id).Select()
	return user, err
}

func (s *Storage) FindByLogin(login string) (*User, error) {
	user := &User{}
	err := s.db.Model(user).Where("login = ?", login).Select()
	return user, err
}

func (s *Storage) FindAll(limit, offset int) (*[]FindUserDTO, error) {
	users := make([]FindUserDTO, 0)
	err := s.db.Model(&users).Limit(limit).Offset(offset).Select()
	return &users, err
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

func (s *Storage) Update(dto *UpdateUserDTO) error {
	_, err := s.db.Model(dto).WherePK().UpdateNotZero()
	return err
}
