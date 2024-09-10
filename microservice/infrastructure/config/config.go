package config

import (
	"encoding/json"
	"os"
)

type ServerConfig struct {
	Addr string `json:"addr"`
}

type Config struct {
	Servers  []ServerConfig `json:"servers"`
	LogLevel string         `json:"logLevel"`
	Timeout  int            `json:"timeout"`
}

func LoadConfig(path string) (*Config, error) {

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
