package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	NIC         NIC         `yaml:"nic"`
	Home        Home        `yaml:"home"`
	DnsProvider DnsProvider `yaml:"dns_provider"`
	DnsServers  []string    `yaml:"dns_servers"`
}

type NIC struct {
	Label     string `yaml:"label"`
	IpAddress string `yaml:"ip_address"`
	Netmask   string `yaml:"netmask"`
}

type Home struct {
	WifiName  string `yaml:"wifi_name"`
	DnsServer string `yaml:"dns_server"`
}

type DnsProvider struct {
	Url string      `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func NewConfig(filename string) (*Config, error) {
	if _, err := os.Lstat(filename); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := new(Config)

	if err = yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

