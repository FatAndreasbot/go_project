package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HS256_SECRET string
	AccessJWTExpiration time.Duration
	RefreshJWTExpiration time.Duration
	DBConnString string
	NetworkPort int
}

var configInstance *Config

func GetConfig() *Config {
	return configInstance
}

func LoadConfig() error {
	secret, ok := os.LookupEnv("HS256_SECRET")
	if !ok {
		return errors.New("environment variable HS256_SECRET is not set")
	}
	configInstance.HS256_SECRET = secret

	expString, ok := os.LookupEnv("ACCESS_JWT_EXPIRATION")
	expSeconds, err := strconv.Atoi(expString)
	if err != nil {
		return errors.New("could not parse ACCESS_JWT_EXPIRATION into an integer number")
	}
	configInstance.AccessJWTExpiration = time.Second * time.Duration(expSeconds)

	expString, ok = os.LookupEnv("REFRESH_JWT_EXPIRATION")
	expSeconds, err = strconv.Atoi(expString)
	if err != nil {
		return errors.New("could not parse REFRESH_JWT_EXPIRATION into an integer number")
	}
	configInstance.RefreshJWTExpiration = time.Second * time.Duration(expSeconds)


	configInstance.DBConnString, ok = os.LookupEnv("DB_CONN")
	if !ok {
		return errors.New("environment variable DB_CONN is not set")
	}

	portString, ok := os.LookupEnv("PORT")
	if !ok {
		return errors.New("environment variable PORT is not set")
	}
	configInstance.NetworkPort, err = strconv.Atoi(portString)
	if err != nil {
		return errors.New("could not parse PORT into an integer numbe")
	}

	return nil
}
