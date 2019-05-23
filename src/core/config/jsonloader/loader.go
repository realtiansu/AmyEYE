package jsonloader

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	// "errors"

	"core/common/errors"
)

func LoadConfig(file string) (*Config, error) {
	fileDir := filepath.Dir(filepath.Dir(file))
	data, _ := ioutil.ReadFile(file)
	config := Config{}
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.New("config file not json")
	}

	config.Data.Dir = filepath.Join(fileDir, config.Data.Dir)

	for index, point := range config.Points {
		addr := point.Address
		if strings.HasPrefix(addr, "list:") {
			path := strings.TrimPrefix(addr, "list:")
			path = filepath.Join(fileDir, path)

			config.Points[index].Address = "list:" + path
		}
	}

	return &config, nil
}
