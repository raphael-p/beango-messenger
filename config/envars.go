package config

type envarType struct {
	configFilepath string
}

var envars envarType = envarType{
	configFilepath: "BG_CONFIG_FILEPATH",
}
