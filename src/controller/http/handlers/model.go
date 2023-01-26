package handlers

import (
	"test-project/src/controller/http/handlers/auth"
	"test-project/src/controller/http/handlers/users"
)

type Handlers struct {
	Users users.IHandler
	Auth  auth.IHandler
}
