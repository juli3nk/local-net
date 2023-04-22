package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	IpAddresses map[string]IpAddress `yaml:"ip_addresses"`
	Trusted   map[string]string  `yaml:"trusted"`
	Vpn       Vpn                `yaml:"vpn"`
	Dns       Dns                `yaml:"dns"`
}

type IpAddress struct {
	IpAddress string `yaml:"ip_address"`
	Netmask   string `yaml:"netmask"`
}

type Vpn struct {
	Enable bool   `yaml:"enable"`
	Name   string `yaml:"name"`
}

type Dns struct {
	Enable          bool            `yaml:"enable"`
	Credentials     Credentials     `yaml:"credentials"`
	UpstreamServers UpstreamServers `yaml:"upstream_servers"`
	Container       Container       `yaml:"container"`
}

type Credentials struct {
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type UpstreamServers struct {
	Default   []string            `yaml:"default"`
	locations map[string]Location `yaml:"locations"`
}

type Location struct {
	WifiName string `yaml:"wifi_name"`
	Server   string `yaml:"dns_server"`
}

type Container struct {
	Enable    bool   `yaml:"enable"`
	LabelDomain string `yaml:"label_domain"`
	LabelAnswer string `yaml:"label_answer"`
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

	if _, ok := config.IpAddresses["dns"]; !ok {
		return nil, fmt.Errorf("no address with label dns")
	}

	return config, nil
}
