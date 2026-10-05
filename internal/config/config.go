package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() (Config, error) {
	var cfg Config
	i, err := os.UserHomeDir()
	if err != nil {
		return cfg, err
	}
	full := filepath.Join(i, ".gatorconfig.json")
	read, err2 := os.ReadFile(full)
	if err2 != nil {
		return cfg, err2
	}
	a := json.Unmarshal(read, &cfg)
	if a != nil {
		return cfg, a
	}
	return cfg, nil
}
func (c Config) SetUser(username string) error {
	c.Current_user_name = username
	err := write(c)
	if err != nil {
		return err
	}
	return nil
}
func write(cfg Config) error {
	i, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	a, err2 := getConfigFilePath()
	if err2 != nil {
		return err2
	}
	err3 := os.WriteFile(a, i, 0644)
	if err3 != nil {
		return err3
	}
	return nil
}
func getConfigFilePath() (string, error) {
	i, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	full := filepath.Join(i, ".gatorconfig.json")
	return full, nil
}
