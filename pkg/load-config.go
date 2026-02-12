package pkg

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

func LoadConfigFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config

	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".json":
		err = json.Unmarshal(data, &config)
	case ".yml":
		err = yaml.Unmarshal(data, &config)
	case ".yaml":
		err = yaml.Unmarshal(data, &config)
	case ".toml":
		err = toml.Unmarshal(data, &config)
	default:
		return nil, errors.New("unsupported config format")
	}

	if err != nil {
		return nil, err
	}

	return &config, nil
}
