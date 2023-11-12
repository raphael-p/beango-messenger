package config

type envarType struct {
	configFilepath,
	DatabaseHost,
	DatabaseName,
	DatabaseUsername,
	DatabasePassword string
}

var Envars envarType = envarType{
	configFilepath:   "BG_CONFIG_FILEPATH",
	DatabaseHost:     "BG_DB_HOST",
	DatabaseName:     "BG_DB_NAME",
	DatabaseUsername: "BG_DB_USERNAME",
	DatabasePassword: "BG_DB_PASSWORD",
}
