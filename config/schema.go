package config

type config struct {
	Server  serverConfig  `json:"server"`
	Logger  loggerConfig  `json:"logger"`
	Session sessionConfig `json:"session"`
}

type serverConfig struct {
	Port int `json:"port"`
}

type loggerConfig struct {
	Directory    string `json:"directory"`
	Filename     string `json:"filename"`
	DefaultLevel int    `json:"defaulLevel"`
}

type sessionConfig struct {
	SecondsUntilExpiry int `json:"secondsUntilExpiry"`
}
