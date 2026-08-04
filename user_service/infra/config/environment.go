package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

var hs256_SECRET string
var jwtExpiration time.Duration
var dbConnString string
var networkPort int

func GetHS256Secret() string {
	return hs256_SECRET
}

func GetJWTExpiration() time.Time {
	return time.Now().Add(jwtExpiration)
}

func GetDBConnString() string {
	return dbConnString
}

func GetNetworkPort() int {
	return networkPort
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
		return errors.New("could not parse JWT_EXPIRATION into an integer number")
	}
	jwtExpiration = time.Second * time.Duration(expSeconds)

	dbConnString, ok = os.LookupEnv("DB_CONN")
	if !ok {
		return errors.New("environment variable DB_CONN is not set")
	}

	portString, ok := os.LookupEnv("PORT")
	if !ok {
		return errors.New("environment variable PORT is not set")
	}
	networkPort, err = strconv.Atoi(portString)
	if err != nil {
		return errors.New("could not parse PORT into an integer numbe")
	}

	return nil
}
