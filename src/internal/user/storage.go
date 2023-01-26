package user

func (s *Storage) FindOne(id int64) (*FindUserDTO, error) {
	user := &FindUserDTO{}

	if err := s.db.Model(user).Where("id = ?", id).Select(); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) FindByLogin(login string) (*User, error) {
	user := &User{}

	if err := s.db.Model(user).Where("login = ?", login).Select(); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) FindAll(limit, offset int) (*[]FindUserDTO, error) {
	users := make([]FindUserDTO, 0)

	if err := s.db.Model(&users).Limit(limit).Offset(offset).Select(); err != nil {
		return nil, err
	}

	return &users, nil
}

func (s *Storage) Delete(id int64) error {
	panic("not implemented") // TODO: Implement
}

func (s *Storage) Create(dto *UserDTO) (*int64, error) {
	u := User{}

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

func (s *Storage) Update(dto *UserDTO) error {
	panic("not implemented") // TODO: Implement
}
