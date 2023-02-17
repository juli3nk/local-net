package main

import (
	"fmt"
	"os"
	"os/exec"
)

func runAdGuardHome() (*os.Process, error) {
	dataPath := "/opt/adguardhome"

	args := []string{
		"-c",
		fmt.Sprintf("%s/conf/AdGuardHome.yaml", dataPath),
		"-w",
		fmt.Sprintf("%s/work", dataPath),
		"-h",
		"192.168.82.1",
	}

	cmd := exec.Command("/usr/local/bin/AdGuardHome", args...)

	go func() {
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}()

	for {
		if cmd.Process != nil {
			break
		}
	}

	return cmd.Process, nil
}
