package main

import (
	"fmt"

	"github.com/juli3nk/local-dns/pkg/adguardhome"
	"github.com/juli3nk/local-dns/pkg/dns"
	"github.com/juli3nk/local-dns/pkg/nmcli"
	"github.com/thoas/go-funk"
)

func setDns(cfg *Config, wifi *nmcli.Wifi) (*string, error) {
	// Get DNS IP
	var dnsIP string

	if cfg.Home.WifiName == wifi.Name {
		dnsIP = cfg.Home.DnsServer
	}

	// check if dns server is reachable
	if len(dnsIP) == 0 {
		for _, s := range cfg.DnsServers {
			err := dns.ServerRespond("github.com", s)
			if err != nil {
				fmt.Println(err)
			}
			if err == nil {
				dnsIP = s

				break
			}
		}
	}

	// otherwise get dns ip from dhcp
	if len(dnsIP) == 0 {
		ip, err := nmcli.GetDhcpDnsIP(wifi.Uuid)
		if err == nil {
			dnsIP = ip
		}
	}

	if len(dnsIP) == 0 {
		return nil, fmt.Errorf("no dns ip available")
	}

	// set dns ip
	d, err := adguardhome.New(cfg.DnsProvider.Url, cfg.DnsProvider.Username, cfg.DnsProvider.Password)
	if err != nil {
		return nil, err
	}

	dnsConfig, err := d.GetDnsConfig()
	if err != nil {
		return nil, err
	}

	if !funk.Contains(dnsConfig.UpstreamDns, dnsIP) {
		dnsConfig.UpstreamDns = []string{dnsIP}

		if err := d.SaveDnsConfig(dnsConfig); err != nil {
			return nil, err
		}
	}

	return &dnsIP, nil
}
