package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home dir: %v", err)
	}
	return filepath.Join(homeDir, configFileName), nil
}

func Read() (Config, error) {
	var myConfig Config

	configPath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read file: %v", err)
	}

	jsonErr := json.Unmarshal(content, &myConfig)
	if jsonErr != nil {
		return Config{}, fmt.Errorf("error converting to json: %v", jsonErr)
	}

	return myConfig, nil
}

func write(cfg Config) error {
	content, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error converting to json: %v", err)
	}

	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	writeErr := os.WriteFile(configPath, content, 0600)
	if writeErr != nil {
		return fmt.Errorf("failed to write to file: %v", writeErr)
	}

	return nil
}

func (c Config) SetUser(username string) error {
	c.CurrentUserName = username
	err := write(c)
	if err != nil {
		return err
	}
	return nil
}
