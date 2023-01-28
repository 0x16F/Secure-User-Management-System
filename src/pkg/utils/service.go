package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"math/big"

	"github.com/goccy/go-json"
)

func GenerateString(length int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		ret[i] = letters[num.Int64()]
	}
	return string(ret)
}

func HashString(str, salt string) (string, error) {
	hasher := sha512.New()

	if _, err := hasher.Write([]byte(str + salt)); err != nil {
		return "", err
	}

	hashed := hasher.Sum(nil)

	encoded := base64.StdEncoding.EncodeToString(hashed)

	return encoded, nil
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func TypeConverter[R any](data any) (*R, error) {
	var result R
	b, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}
