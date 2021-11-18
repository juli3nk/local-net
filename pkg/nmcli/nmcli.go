package nmcli

import (
	"os/exec"
	"regexp"
	"strings"
)

const NMCLIBIN string = "/usr/bin/nmcli"

func GetConnectedWifi() (*Wifi, error) {
	c := exec.Command(NMCLIBIN, "c", "show", "--active")

	o, err := c.Output()
	if err != nil {
		return nil, err
	}

	result := strings.Split(string(o), "\n")

	re := regexp.MustCompile(`^(.*)\s+([a-z0-9]{8}\-[a-z0-9]{4}\-[a-z0-9]{4}\-[a-z0-9]{4}\-[a-z0-9]{12})\s+wifi\s+(.*)$`)

	for _, l := range result {
		match := re.FindStringSubmatch(l)

		if len(match) > 0 {
			w := Wifi{
				Name: strings.TrimSpace(match[1]),
				Uuid: strings.TrimSpace(match[2]),
				Nic:  strings.TrimSpace(match[3]),
			}

			return &w, nil
		}
	}

	return nil, nil
}

func GetDhcpDnsIP(uuid string) (string, error) {
	c := exec.Command(NMCLIBIN, "c", "show", uuid)

	o, err := c.Output()
	if err != nil {
		return "", err
	}

	result := strings.Split(string(o), "\n")

	re := regexp.MustCompile(`^IP4\.DNS\[1\]:\s+([0-9\.]+)$`)

	for _, l := range result {
		match := re.FindStringSubmatch(l)

		if len(match) > 0 {
			ip := match[1]

			return ip, nil
		}
	}

	return "", nil
}
