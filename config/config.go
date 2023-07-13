package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/raphael-p/beango/utils/path"
	"github.com/raphael-p/beango/utils/validate"
)

var Values *config

func CreateConfig() {
	filePath := os.Getenv(envars.configFilepath)
	if filePath == "" {
		fmt.Printf("$%s not set, using default config\n", envars.configFilepath)
		path, ok := path.RelativeJoin("default.json")
		if !ok {
			panic("could not get absolute path for config directory")
		}
		filePath = path
	}

	file, err := os.Open(filePath)
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
