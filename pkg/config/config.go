package config

import (
  "encoding/json"
  "fmt"
  "os"
)

func New(filename string) (*Config, error) {
  config := new(Config)

  if err := config.ReadFile(filename); err != nil {
    return nil, err
  }

  // Check that ip aliases all have the same network addr
  // Check if ip aliases are already set
  // If ip aliases are not yet set, make sure they don't respond to ping

  if _, ok := config.IpAddresses["dns"]; !ok {
    return nil, fmt.Errorf("no address with label dns")
  }

  // If vpn enabled, check that connection name exists

  return config, nil
}

func (c *Config) ReadFile(filename string) error {
  if _, err := os.Lstat(filename); err != nil {
    return err
  }

  data, err := os.ReadFile(filename)
  if err != nil {
    return err
  }

  if err = json.Unmarshal(data, c); err != nil {
    return err
  }

  return nil
}

/*
func (c *Config) ValidateIpAddressNetwork() error {
}
*/

func (c *Config) IsWifiTrusted(wifiName string) bool {
  for _, name := range c.Wifi {
    if name == wifiName {
      return true
    }
  }

  return false
}
