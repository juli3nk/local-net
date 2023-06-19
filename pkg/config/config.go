package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func New(filename string) (*Config, error) {
	if _, err := os.Lstat(filename); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := new(Config)

	if err = json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	if _, ok := config.IpAddresses["dns"]; !ok {
		return nil, fmt.Errorf("no address with label dns")
	}

	return config, nil
}

func (c *Config) IsWifiTrusted(wifiName string) bool {
	for _, wt := range c.Wifi {
		if wt == wifiName {
			return true
		}
	}

	return false
}
