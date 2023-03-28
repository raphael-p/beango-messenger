package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

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

var Config *config

func CreateConfig() {
	file, err := os.Open("config/default.json")
	if err != nil {
		fl.Log(fmt.Sprint("could not open config file: ", err))
	}
	defer file.Close()

	Config = &config{}
	if err = json.NewDecoder(file).Decode(Config); err != nil {
		fl.Log(fmt.Sprint("could not parse config file: ", err))
	}
	if fields := ValidateRequiredFields(Config); len(fields) != 0 {
		fl.Log(fmt.Sprint("missing required config field(s): ", fields))
	}
}
