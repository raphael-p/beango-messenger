package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/raphael-p/beango/utils"
)

var Config *config

func Init() {
	file, err := os.Open("config/default.json")
	if err != nil {
		utils.StaticFatal(fmt.Sprint("could not open config file: ", err))
	}
	defer file.Close()

	Config = &config{}
	if err = json.NewDecoder(file).Decode(Config); err != nil {
		utils.StaticFatal(fmt.Sprint("could not parse config file: ", err))
	}
	if fields := utils.ValidateRequiredFields(Config); len(fields) != 0 {
		utils.StaticFatal(fmt.Sprint("missing required config field(s): ", fields))
	}
}
