package inits

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"nyacap-mini/config"
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
	if len(config.Config.Captcha.Characters) == 0 {
		return fmt.Errorf("missing captcha characters")
	}

	// Fill with defaults
	if config.Config.Captcha.CheckTextLen.Max == 0 {
		config.Config.Captcha.CheckTextLen.Max = 5
	}
	if config.Config.Captcha.CheckTextLen.Min == 0 {
		config.Config.Captcha.CheckTextLen.Min = 3
	}

	return nil
}
