package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/raphael-p/beango/utils/validate"
)

var Values *config

func CreateConfig(fail func(string)) {
	file, err := os.Open("config/default.json")
	if err != nil {
		fail(fmt.Sprint("could not open config file: ", err))
	}
	defer file.Close()

	Values = &config{}
	if err = json.NewDecoder(file).Decode(Values); err != nil {
		fail(fmt.Sprint("could not parse config file: ", err))
	}
	if fields := validate.DeserialisedJSON(Values); len(fields) != 0 {
		fail(fmt.Sprint("missing required config field(s): ", fields))
	}
}
