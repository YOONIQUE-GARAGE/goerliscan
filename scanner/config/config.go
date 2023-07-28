package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type DatabaseConfig struct {
	Host string `toml:"host"`
	Name string `toml:"name"`																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																															
}

type NetworkConfig struct {
	URL string `toml:"url"`
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
	Netowrk NetworkConfig `toml:"network"`
	Log LogConfig `toml:"log"`
}	

func LoadCofig() (*Config, error) {
	var config Config
	_, err := toml.DecodeFile("config/config.toml", &config)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &config, nil
}
