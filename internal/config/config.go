package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// This is the name of the file that will be created in the user's home directory
const fileName = ".gatorconfig.json"

// Config is a struct that holds the configuration for the application
type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Read reads the configuration from the user's home directory
func Read() (*Config, error) {
	config := &Config{}

	file, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	configFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	err = configFile.Close()
	if err != nil {
		return nil, err
	}

	return config, nil
}

// SetUser sets the current user in the configuration
func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName

	err := write(cfg)
	if err != nil {
		return err
	}

	return nil
}

// Returns the path to the configuration file
func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, fileName), nil
}

// Writes the configuration to the user's home directory
func write(cfg *Config) error {
	file, err := getConfigFilePath()
	if err != nil {
		return err
	}

	configFile, err := os.Create(file)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(configFile)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	err = configFile.Close()
	if err != nil {
		return err
	}

	return nil
}
