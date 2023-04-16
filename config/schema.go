package config

import "github.com/raphael-p/beango/utils/validate"

type config struct {
	Server  serverConfig  `json:"server"`
	Logger  loggerConfig  `json:"logger"`
	Session sessionConfig `json:"session"`
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
	SecondsUntilExpiry uint16 `json:"secondsUntilExpiry"`
}
