package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func NewService(secrets *Service) Servicer {
	return &Service{
		AccessSecret:  secrets.AccessSecret,
		RefreshSecret: secrets.RefreshSecret,
	}
}

func (s *Service) GenerateAccess(dto *GenerateDTO) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, Token{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Minute * 5).Unix(),
		},
		Id:          dto.Id,
		Permissions: dto.Permissions,
	})

	str, err := token.SignedString([]byte(s.AccessSecret))

	return str, err
}

func (s *Service) GenerateRefresh(dto *GenerateDTO) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, Token{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 30).Unix(),
		},
		Id:          dto.Id,
		Permissions: dto.Permissions,
	})

	str, err := token.SignedString([]byte(s.RefreshSecret))

	return str, err
}

func (s *Service) ParseAccess(token string) (*Token, error) {
	t, err := jwt.ParseWithClaims(token, &Token{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.AccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(*Token); ok && t.Valid {
		if time.Since(time.Now()).Milliseconds() > claims.ExpiresAt {
			return nil, ErrExpired
		}

		return claims, nil
	}

	return nil, err
}

func (s *Service) ParseRefresh(token string) (*Token, error) {
	t, err := jwt.ParseWithClaims(token, &Token{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.RefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(*Token); ok && t.Valid {
		if time.Since(time.Now()).Milliseconds() > claims.ExpiresAt {
			return nil, ErrExpired
		}

		return claims, nil
	}

	return nil, err
}
