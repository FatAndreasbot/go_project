package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

var hs256_SECRET string
var jwtExpiration time.Duration

func GetHS256Secret() string {
	return hs256_SECRET
}

func GetJWTExpiration() time.Time {
	return time.Now().Add(jwtExpiration)
}

func InitEnvVars() error {
	secret, ok := os.LookupEnv("HS256_SECRET")
	if !ok {
		return errors.New("environment variable HS256_SECRET is not set")
	}
	hs256_SECRET = secret

	expString, ok := os.LookupEnv("JWT_EXPIRATION")
	expSeconds, err := strconv.Atoi(expString)
	if err != nil {
		return errors.New("could not pars JWT_EXPIRATION into an integer number")
	}
	jwtExpiration = time.Second * time.Duration(expSeconds)

	return nil
}
