package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/raphael-p/beango/utils/validate"
)

var Values *config

func CreateConfig() {
	filePath := os.Getenv(envars.configFilepath)
	if filePath == "" {
		fmt.Printf("$%s not set, using default config\n", envars.configFilepath)
		_, dirname, _, ok := runtime.Caller(0)
		if !ok {
			panic("could get absolute path for config directory")
		}
		filePath = filepath.Join(filepath.Dir(dirname), "default.json")
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
