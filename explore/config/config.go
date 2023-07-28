package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type DatabaseConfig struct {
	Host string `toml:"host"`
	Name string `toml:"name"`
}
type LogConfig struct {
	Level   string
	Fpath   string
	Msize   int
	Mage    int
	Mbackup int
}

type Config struct {
	Database DatabaseConfig `toml:"database"`
	Log LogConfig `toml:"log"`
}	

func LoadConfig() (*Config, error) {
	var config Config
	_, err := toml.DecodeFile("config/config.toml", &config)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &config, nil
}