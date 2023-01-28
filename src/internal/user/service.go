package user

import (
	"test-project/src/pkg/utils"

	"github.com/go-pg/pg/v10"
)

func NewUser(dto *CreateUserDTO) *UserDTO {
	salt := utils.GenerateString(SaltLength)
	hash, _ := utils.HashString(dto.Password, salt)

	user := &UserDTO{
		Name:        dto.Name,
		Surname:     dto.Surname,
		Login:       dto.Login,
		Password:    hash,
		Salt:        salt,
		Permissions: dto.Permissions,
		Birthday:    dto.Birthday,
	}

	return user
}

func NewStorage(db *pg.DB) Storager {
	return &Storage{
		db: db,
	}
}
