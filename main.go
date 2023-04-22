package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/juli3nk/local-net/pkg/ip"
	"github.com/juli3nk/local-net/pkg/nmcli"
	log "github.com/sirupsen/logrus"
)

var (
	flgConfig   string
	flgDebug    bool
	flgInterval int
)

func init() {
	flag.StringVar(&flgConfig, "config", "/tmp/local-net.yml", "config file path")
	flag.BoolVar(&flgDebug, "debug", false, "enable debug log")
	flag.IntVar(&flgInterval, "interval", 30, "interval between probing wifi informations")

	flag.Parse()
}

func main() {
	if flgDebug {
		log.SetLevel(log.DebugLevel)
	}

	var currentWifiName string
	var agh *os.Process

	cfg, err := NewConfig(flgConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug(cfg)

	ticker := time.NewTicker(time.Duration(flgInterval) * time.Second)
	done := make(chan bool)

	probe := func() {
		wifi, err := nmcli.GetConnectedWifi()
		if err != nil {
			log.Error(err)

			return
		}

		if wifi == nil {
			log.Info("wifi not connected")

			return
		}
		log.Debug(wifi)

		// Check if dns interface is set
		i, err := ip.New(wifi.Nic, cfg.NIC.Label, cfg.NIC.IpAddress, cfg.NIC.Netmask)
		if err != nil {
			log.Fatal(err)
		}

		if !i.IsSet() {
			log.Info("set dns interface ip address")

			if err := i.Set(); err != nil {
				log.Error(err)
				return
			}
		}

		if agh == nil {
			agh, err = runAdGuardHome(cfg.NIC.IpAddress)
			if err != nil {
				log.Fatal(err)
			}
			log.Info("AdGuardHome started")
			log.Debugf("AdGuardHome process ID is %d", agh.Pid)
		}

		// Update Adguardhome
		if wifi.Name == currentWifiName {
			return
		}

		dnsIP, err := setDns(cfg, wifi)
		if err != nil {
			log.Error(err)

			return
		}
		log.Infof("set dns IP to %s", *dnsIP)

		currentWifiName = wifi.Name
	}

	log.Info("start probing")
	probe()

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				probe()
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	agh.Kill()

	ticker.Stop()
	done <- true
}
