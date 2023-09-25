package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config main struct
type Config struct {
	Sources `yaml:"sources"`
}

// Sources struct for config
type Sources struct {
	Github map[string]GithubSource `yaml:"github"`
}

// GithubSource struct for congig
type GithubSource struct {
	User     string `yaml:"user"`
	Repo     string `yaml:"repo"`
	Path     string `yaml:"path"`
	Branch   string `yaml:"branch"`
	Tag      string `yaml:"tag"`
	SyncPath string `yaml:"syncPath"`
}

// Parse config fuction
func Parse(filepath string) (*Config, error) {
	configFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
