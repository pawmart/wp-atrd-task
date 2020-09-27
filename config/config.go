package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type Config struct {
	Server struct {
		Address string `yaml:"address"`
		Port    string `yaml:"port"`
		Timeout int64  `yaml:"timeout"`
	} `yaml:"server"`

	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		DB       int    `yaml:"database"`
	} `yaml:"redis"`
}

func New(path string) (*Config, error) {
	if filepath.Ext(path) != ".yaml" {
		return nil, errors.New("support only .yaml format")
	}

	config := &Config{}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
