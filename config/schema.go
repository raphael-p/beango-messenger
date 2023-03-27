package config

type config struct {
	Server  ServerConfig  `json:"server"`
	Logger  LoggerConfig  `json:"logger"`
	Session SessionConfig `json:"session"`
}

type ServerConfig struct {
	Port int `json:"port"`
}

type LoggerConfig struct {
	Directory    string `json:"directory"`
	Filename     string `json:"filename"`
	DefaultLevel int    `json:"defaulLevel"`
}

type SessionConfig struct {
	SecondsUntilExpiry int `json:"secondsUntilExpiry"`
}
