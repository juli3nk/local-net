package main

import (
	"fmt"

	"github.com/juli3nk/local-net/pkg/adguardhome"
	"github.com/juli3nk/local-net/pkg/dns"
	"github.com/juli3nk/local-net/pkg/nmcli"
	"github.com/rs/zerolog/log"
)

func testDns(servers []string) []string {
	var result []string

	for _, s := range servers {
		err := dns.ServerRespond("github.com", s)
		if err != nil {
			log.Error().Err(err).Send()
		}
		if err == nil {
			result = append(result, s)
		}
	}

	return result
}

func setDnsUpstreamServers(cfgDns *Dns, aghcli *adguardhome.DnsCfg, wifiTrusted bool, wifiUuid, vpnUuid string) ([]string, error) {
	var dnsServers []string

	// Get upstream dns servers
	// vpn connected
	if vpnUuid != "" {
		ds, err := nmcli.GetConnectionDhcpDns(vpnUuid)
		if err == nil {
			result := testDns(ds)
			if len(result) > 0 {
				dnsServers = result
			}
		}
	}

	// wifi trusted
	if wifiTrusted {
		ds, err := nmcli.GetConnectionDhcpDns(wifiUuid)
		if err == nil {
			result := testDns(ds)
			if len(result) > 0 {
				dnsServers = result
			}
		}
	}

	// default
	if len(dnsServers) == 0 {
		result := testDns(cfgDns.UpstreamServers.Default)
		if len(result) > 0 {
			dnsServers = result
		}
	}

	// otherwise get dns ip from dhcp
	if len(dnsServers) == 0 {
		ds, err := nmcli.GetConnectionDhcpDns(wifiUuid)
		if err == nil {
			result := testDns(ds)
			if len(result) > 0 {
				dnsServers = result
			}
		}
	}

	if len(dnsServers) == 0 {
		return nil, fmt.Errorf("no dns server available")
	}

	// Set dns ip
	dnsConfig, err := aghcli.GetDnsConfig()
	if err != nil {
		return nil, err
	}
	dnsConfig.UpstreamDns = dnsServers

	if err := aghcli.SaveDnsConfig(dnsConfig); err != nil {
		return nil, err
	}

	return dnsServers, nil
}
