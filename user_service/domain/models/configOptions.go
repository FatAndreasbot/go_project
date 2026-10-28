package models

import "time"

type ConfigOptions struct {
	dbConnString         string
	hs256Secret          string
	accsessJWTExpiration time.Duration
	refreshJWTExpiration time.Duration
	port                 int
	redisAddress         string
	redisPassword        string
	redisDB              int
}

func NewOptions(
	dbConnString, hs256Secret, redisAddress, redisPassword string,
	accsessJWTExpiration, refreshJWTExpiration time.Duration,
	port, redisDB int,
) *ConfigOptions {
	return &ConfigOptions{
		dbConnString:         dbConnString,
		hs256Secret:          hs256Secret,
		accsessJWTExpiration: accsessJWTExpiration,
		refreshJWTExpiration: refreshJWTExpiration,
		port:                 port,
		redisAddress:         redisAddress,
		redisPassword:        redisPassword,
		redisDB:              redisDB,
	}
}

func (opt *ConfigOptions) GetDBConnString() string {
	return opt.dbConnString
}

func (opt *ConfigOptions) GetHS256Secret() string {
	return opt.hs256Secret
}

func (opt *ConfigOptions) GetAccessJWTExpiration() time.Duration {
	return opt.accsessJWTExpiration
}

func (opt *ConfigOptions) GetRefreshJWTExpiration() time.Duration {
	return opt.refreshJWTExpiration
}

func (opt *ConfigOptions) GetPort() int {
	return opt.port
}

func (opt *ConfigOptions) GetRedisAddress() string {
	return opt.redisAddress
}

func (opt *ConfigOptions) GetRedisPassword() string {
	return opt.redisPassword
}

func (opt *ConfigOptions) GetRedisDB() int {
	return opt.redisDB
}
