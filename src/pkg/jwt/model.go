package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

type Token struct {
	jwt.StandardClaims
	Id          int64  `json:"id"`
	Permissions string `json:"permissions"`
}

type Service struct {
	AccessSecret  string `json:"access"`
	RefreshSecret string `json:"refresh"`
}

type GenerateDTO struct {
	Id          int64  `json:"id"`
	Permissions string `json:"permissions"`
}

type Servicer interface {
	ParseAccess(token string) (*Token, error)
	ParseRefresh(token string) (*Token, error)
	GenerateAccess(dto *GenerateDTO) (string, error)
	GenerateRefresh(dto *GenerateDTO) (string, error)
}

var (
	ErrExpired = errors.New("expired")
)
