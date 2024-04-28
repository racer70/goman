package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Domain   string `json:"domain"`
	Services map[string]struct {
		Domain string `json:"domain"`
	} `json:"services"`
}

func (c Config) GetDomain(app string) string {
	_, ok := c.Services[app]

	if ok {
		return c.Services[app].Domain
	} else {
		return c.Domain
	}
}

func GetConfig(env string) (map[string]Config, error) {
	f, err := os.ReadFile(fmt.Sprintf("%s", "config/env.json"))

	if err != nil {
		fmt.Printf("Error reading config file, %v", err)
		return nil, err
	}

	config := map[string]Config{}

	err = json.Unmarshal(f, &config)

	if err != nil {
		fmt.Printf("Error extracting config file,  %v", err)
		return nil, err
	}

	return config, nil
}
