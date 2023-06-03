package config

import "github.com/raphael-p/beango/utils/validate"

type config struct {
	Server   serverConfig   `json:"server"`
	Logger   loggerConfig   `json:"logger"`
	Session  sessionConfig  `json:"session"`
	Database databaseConfig `json:"database"`
}

type serverConfig struct {
	Port uint16 `json:"port"`
}

type loggerConfig struct {
	Directory    string                    `json:"directory"`
	Filename     string                    `json:"filename"`
	DefaultLevel validate.JSONField[uint8] `json:"defaulLevel" zeroable:"true"`
}

type sessionConfig struct {
	SecondsUntilExpiry uint32 `json:"secondsUntilExpiry"`
}

type databaseConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Name string `json:"name"`
}
