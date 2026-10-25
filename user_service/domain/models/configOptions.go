package models

import "time"

type ConfigOptions struct {
	dbConnString         string
	hs256Secret          string
	accsessJWTExpiration time.Duration
	refreshJWTExpiration time.Duration
	port                 int
}

func NewOptions(
	dbConnString, hs256Secret string,
	accsessJWTExpiration, refreshJWTExpiration time.Duration,
	port int,
) *ConfigOptions {
	return &ConfigOptions{
		dbConnString:         dbConnString,
		hs256Secret:          hs256Secret,
		accsessJWTExpiration: accsessJWTExpiration,
		refreshJWTExpiration: refreshJWTExpiration,
		port:                 port,
	}
}

func (opt *ConfigOptions) GetDBConnString() string {
	return opt.dbConnString
}

func (opt *ConfigOptions) GetHS256Secret() string {
	return opt.hs256Secret
}

func (opt *ConfigOptions) GetAccsessJWTExpiration() time.Duration {
	return opt.accsessJWTExpiration
}

func (opt *ConfigOptions) GetRefreshJWTExpiration() time.Duration {
	return opt.refreshJWTExpiration
}

func (opt *ConfigOptions) GetPort() int {
	return opt.port
}
