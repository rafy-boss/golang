package method

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string) (string, error) {
	 hashed,err:=bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	 return string(hashed),err
}


func CompareHashAndPassword(hashedPassword, password string) (error) {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}