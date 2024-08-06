package local

import (
    "golang.org/x/crypto/bcrypt"

    viper "github.com/spf13/viper"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), viper.GetInt("hash.salt"))
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}