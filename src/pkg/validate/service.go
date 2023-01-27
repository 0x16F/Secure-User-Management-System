package validate

import (
	"regexp"
	"test-project/src/internal/permissions"
	"test-project/src/pkg/utils"
)

func LoginLenght(login string) bool {
	loginLength := len(login)
	return loginLength >= minLoginLength && loginLength <= maxLoginLength
}

func Login(login string) bool {
	return regexp.MustCompile(`(?m)^[A-Za-z0-9]+$`).Match([]byte(login))
}

func PasswordLength(password string) bool {
	passwordLength := len(password)
	return passwordLength >= minPasswordLength && passwordLength <= maxPasswordLength
}

func Permission(permission string) bool {
	return utils.Contains(permissions.ArrayOfPermissions, permission)
}