package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(p string) (string, error) {
	bytePwd := []byte(p)

	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePassword(hashedPwd, inputPwd string) bool {
	byteHash := []byte(hashedPwd)
	byteInput := []byte(inputPwd)

	if err := bcrypt.CompareHashAndPassword(byteHash, byteInput); err != nil {
		return false
	}

	return true
}
