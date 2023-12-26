package inits

import (
	"gopkg.in/yaml.v3"
	"nya-captcha/config"
	"os"
)

func Config() error {
	// Read config file
	configFileBytes, err := os.ReadFile("config.yml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configFileBytes, &config.Config)
	if err != nil {
		return err
	}

	// Validate config

	return nil
}
