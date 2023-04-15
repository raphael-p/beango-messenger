package config

import "github.com/raphael-p/beango/utils/validate"

type config struct {
	Server  serverConfig  `json:"server"`
	Logger  loggerConfig  `json:"logger"`
	Session sessionConfig `json:"session"`
}

type serverConfig struct {
	Port int `json:"port"`
}

type loggerConfig struct {
	Directory    string                  `json:"directory"`
	Filename     string                  `json:"filename"`
	DefaultLevel validate.JSONField[int] `json:"defaulLevel" zeroable:"true"`
}

type sessionConfig struct {
	SecondsUntilExpiry int `json:"secondsUntilExpiry"`
}
