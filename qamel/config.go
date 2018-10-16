package main

import (
	"encoding/json"
	"os"
	fp "path/filepath"
)

// config is struct containing path to qmake, moc and rcc
type config struct {
	Qmake string
	Moc   string
	Rcc   string
}

// loadConfig loads config file in ${XDG_CONFIG_HOME}/qamel/config
func loadConfig() (config, error) {
	// Open file
	configFile, err := os.Open(configPath)
	if err != nil {
		return config{}, err
	}
	defer configFile.Close()

	// Decode JSON
	cfg := config{}
	err = json.NewDecoder(configFile).Decode(&cfg)
	return cfg, err
}

// saveConfig saves the config as JSON in ${XDG_CONFIG_HOME}/qamel/config
func saveConfig(cfg config) error {
	// Make sure config dir is exists
	os.MkdirAll(fp.Dir(configPath), os.ModePerm)

	// Create file
	configFile, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	// Encode config to JSON
	return json.NewEncoder(configFile).Encode(&cfg)
}
