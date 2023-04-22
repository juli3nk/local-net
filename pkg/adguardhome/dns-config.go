package adguardhome

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

/*
{
	"bootstrap_dns": ["9.9.9.10","149.112.112.10","2620:fe::10","2620:fe::fe:10"],
	"upstream_mode": "",
	"resolve_clients": true,
	"local_ptr_upstreams": [],
	"upstream_dns": ["1.1.1.1"]
}
*/

type DnsConfig struct {
	BootstrapDns      []string `json:"bootstrap_dns"`
	UpstreamMode      string   `json:"upstream_mode,omitempty"`
	ResolveClients    bool     `json:"resolve_clients"`
	LocalPtrUpstreams []string `json:"local_ptr_upstreams,omitempty"`
	UpstreamDns       []string `json:"upstream_dns"`
}

func (c *DnsCfg) GetDnsConfig() (*DnsConfig, error) {
	url := c.url + "/control/dns_info"

	contentType := "application/json"

	tr := &http.Transport{
		IdleConnTimeout: 5 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Got error %s", err.Error())
	}
	req.Header.Add("Content-Type", contentType)
	req.SetBasicAuth(c.username, c.password)

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Got error %s", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Please contact an administrator if the problem persists (method: GetDnsConfig, status code: %d)", response.StatusCode)
	}

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	data := new(DnsConfig)

	if err := json.Unmarshal(respBody, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *DnsCfg) SaveDnsConfig(payload *DnsConfig) error {
	url := c.url + "/control/dns_config"

	contentType := "application/json"

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		IdleConnTimeout: 5 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	req.Header.Add("Content-Type", contentType)
	req.SetBasicAuth(c.username, c.password)

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("Please contact an administrator if the problem persists (method: SaveDnsConfig, status code: %d)", response.StatusCode)
	}

	return nil
}
