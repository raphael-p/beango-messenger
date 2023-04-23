package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/raphael-p/beango/utils/validate"
)

var Values *config

func CreateConfig() {
	file, err := os.Open("config/default.json")
	if err != nil {
		panic(fmt.Sprint("could not open config file: ", err))
	}
	defer file.Close()

	Values = &config{}
	if err = json.NewDecoder(file).Decode(Values); err != nil {
		panic(fmt.Sprint("could not parse config file: ", err))
	}
	fields, err := validate.StructFromJSON(Values)
	if err != nil {
		panic(err)
	}
	if len(fields) != 0 {
		panic(fmt.Sprint("missing required config field(s): ", fields))
	}
}
