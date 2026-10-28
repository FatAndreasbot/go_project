package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
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
	RedisAddress
	RedisPassword
	RedisDB
)

var configKeys map[ConfigOption]string = map[ConfigOption]string{
	DBConnString:         "DB_CONN",
	HS256_SECRET:         "HS256_SECRET",
	AccessJWTExpiration:  "ACCESS_JWT_EXPIRATION",
	RefreshJWTExpiration: "REFRESH_JWT_EXPIRATION",
	NetworkPort:          "PORT",
	RedisAddress:         "REDIS_ADDRESS",
	RedisPassword:        "REDIS_PASSWORD",
	RedisDB:              "REDIS_DB",
}

// Depricated: insert config values at init
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

func LoadConfig(envFileLocation ...string) (*models.ConfigOptions, error) {
	if err := godotenv.Load(envFileLocation...); err != nil {
		return nil, err
	}
	dbConnString, ok := os.LookupEnv(configKeys[DBConnString])
	if !ok {
		return nil, fmt.Errorf("%s was not found", configKeys[DBConnString])
	}
	hs256Secret, ok := os.LookupEnv(configKeys[HS256_SECRET])
	if !ok {
		return nil, fmt.Errorf("%s was not found", configKeys[HS256_SECRET])
	}
	accsessJWTExpirationString, ok := os.LookupEnv(configKeys[AccessJWTExpiration])
	if !ok {
		return nil, fmt.Errorf("%s was not found", configKeys[AccessJWTExpiration])
	}
	refreshJWTExpirationString, ok := os.LookupEnv(configKeys[RefreshJWTExpiration])
	if !ok {
		return nil, fmt.Errorf("%s was not found", configKeys[RefreshJWTExpiration])
	}
	portString, ok := os.LookupEnv(configKeys[NetworkPort])
	if !ok {
		return nil, fmt.Errorf("%s was not found", configKeys[NetworkPort])
	}
	redisAddress, ok := os.LookupEnv(configKeys[RedisAddress])
	if !ok {
		return nil, fmt.Errorf("%s was not found", configKeys[RedisAddress])
	}
	redisPassword, ok := os.LookupEnv(configKeys[RedisPassword])
	if !ok {
		return nil, fmt.Errorf("%s was not found", configKeys[RedisPassword])
	}
	redisDBString, ok := os.LookupEnv(configKeys[RedisDB])
	if !ok {
		return nil, fmt.Errorf("%s was not found", configKeys[RedisDB])
	}

	accsessJWTExpiration, err := strconv.Atoi(accsessJWTExpirationString)
	if err != nil {
		return nil, err
	}
	refreshJWTExpiration, err := strconv.Atoi(refreshJWTExpirationString)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(portString)
	if err != nil {
		return nil, err
	}
	redisDB, err := strconv.Atoi(redisDBString)
	if err != nil {
		return nil, err
	}

	options := models.NewOptions(
		dbConnString, hs256Secret, redisAddress, redisPassword,
		time.Second*time.Duration(accsessJWTExpiration),
		time.Second*time.Duration(refreshJWTExpiration),
		port, redisDB,
	)

	return options, nil
}
