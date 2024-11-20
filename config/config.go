package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SoundProfile string `yaml:"sound_profile"`
	Volume       int    `yaml:"volume"`
	Enabled      bool   `yaml:"enabled"`
}

func LoadConfig() (*Config, error) {
	// Check if the config directory exists, create it if not
	if _, err := os.Stat("config"); os.IsNotExist(err) {
		err = os.Mkdir("config", 0755)
		if err != nil {
			return nil, err
		}
	}

	// Try to read the configuration file
	data, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			// Configuration file does not exist, create default config
			cfg := &Config{
				SoundProfile: "classic",
				Volume:       80,
				Enabled:      true,
			}
			// Save the default config file
			err = SaveConfig(cfg)
			if err != nil {
				return nil, err
			}
			return cfg, nil
		}
		// Other errors
		return nil, err
	}
	// Unmarshal the YAML data into Config struct
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	// Marshal the Config struct into YAML data
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	// Write the YAML data to the configuration file
	err = ioutil.WriteFile("config/config.yaml", data, 0644)
	if err != nil {
		return err
	}
	return nil
}
