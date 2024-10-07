package config

import (
	"Burst/pkg/models"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl"
)

func LoadConfig(filename string) (*models.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config models.Config
	err = hcl.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadConfigs(path string) ([]*models.Config, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var configs []*models.Config

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".burst" {
			continue
		}
		config, err := LoadConfig(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func LoadRootConfig() (*models.Config, error) {
	return LoadConfig("Burstfile")
}
