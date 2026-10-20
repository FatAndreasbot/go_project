package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HS256_SECRET         []byte
	AccessJWTExpiration  time.Duration
	RefreshJWTExpiration time.Duration
	DBConnString         string
	NetworkPort          int
}

type ConfigOption int

const (
	HS256_SECRET ConfigOption = iota
	AccessJWTExpiration
	RefreshJWTExpiration
	DBConnString
	NetworkPort
)

var configKeys map[ConfigOption]string = map[ConfigOption]string{
	DBConnString:         "DB_CONN",
	HS256_SECRET:         "HS256_SECRET",
	AccessJWTExpiration:  "ACCESS_JWT_EXPIRATION",
	RefreshJWTExpiration: "REFRESH_JWT_EXPIRATION",
	NetworkPort:          "PORT",
}

func GetConfig(option ConfigOption) (string, error) {
	key, ok := configKeys[option]
	if !ok {
		return "", errors.New("unacceptable option was given")
	}

	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}

	return value, nil
}

func LoadConfig(envFileLocation ...string) error {
	if err := godotenv.Load(envFileLocation...); err != nil {
		return err
	}

	return nil
}
