package nmcli

import (
	"os/exec"
	"regexp"
	"strings"
)

const NMCLIBIN string = "/usr/bin/nmcli"

func GetDevice(dtype string) (*Device, error) {
	cmd := exec.Command(NMCLIBIN, "device", "status")

	o, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	result := strings.Split(string(o), "\n")

	re := regexp.MustCompile(`^(.*)\s+([a-z\-]+)\s+(connected|disconnected|unmanaged)\s+(.*)$`)

	for _, l := range result {
		match := re.FindStringSubmatch(l)

		if len(match) > 0 && match[2] == dtype {
			device := Device{
				Name: strings.TrimSpace(match[1]),
				Type: strings.TrimSpace(match[2]),
				State:  strings.TrimSpace(match[3]),
				Connection: strings.TrimSpace(match[4]),
			}

			return &device, nil
		}
	}

	return nil, nil
}

func GetConnection(ctype, name string) (*Connection, error) {
	var connection Connection

	cmd := exec.Command(NMCLIBIN, "connection", "show", "--active")

	o, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	result := strings.Split(string(o), "\n")

	re := regexp.MustCompile(`^(.*)\s+([a-z0-9]{8}\-[a-z0-9]{4}\-[a-z0-9]{4}\-[a-z0-9]{4}\-[a-z0-9]{12})\s+([a-z\-]+)\s+(.*)$`)

	for _, l := range result {
		match := re.FindStringSubmatch(l)

		if len(match) > 0 && match[3] == ctype {
			connection = Connection{
				Name: strings.TrimSpace(match[1]),
				Uuid: strings.TrimSpace(match[2]),
				Type: strings.TrimSpace(match[3]),
				Device: strings.TrimSpace(match[4]),
			}

			if match[1] == name {
				break
			}
		}
	}

	return &connection, nil
}

func UpConnection(name string) error {
	//cmd := exec.Command(NMCLIBIN, "connection", "up", name)

	return nil
}

func DownConnection(name string) error {
	//cmd := exec.Command(NMCLIBIN, "connection", "down", name)

	return nil
}

func GetConnectionDhcpDns(uuid string) ([]string, error) {
	var dns []string

	cmd := exec.Command(NMCLIBIN, "c", "show", uuid)

	o, err := cmd.Output()
	if err != nil {
		return dns, err
	}

	result := strings.Split(string(o), "\n")

	re := regexp.MustCompile(`^IP4\.DNS\[1\]:\s+([0-9\.]+)$`)

	for _, l := range result {
		match := re.FindStringSubmatch(l)

		if len(match) > 0 {
			dns = append(dns, match[1])
		}
	}

	return dns, nil
}
