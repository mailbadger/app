package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Configuration struct {
	Database      string `json:"database"`
	MigrationsDir string `json:"migrations_dir"`
}

var Config Configuration

func init() {
	file, err := ioutil.ReadFile("../config.json")

	if err != nil {
		fmt.Printf("Config file error: %v\n", err)
	}

	json.Unmarshal(file, &Config)
}
