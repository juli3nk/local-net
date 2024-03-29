package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/docker/docker/api/types/events"
	"github.com/juli3nk/go-adguardhome"
	"github.com/juli3nk/go-docker"
	"github.com/juli3nk/go-network/ip"
	"github.com/juli3nk/go-network/nmcli"
	"github.com/juli3nk/local-net/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thoas/go-funk"
)

var (
	flgConfig   string
	flgDebug    bool
	flgInterval int
)

func init() {
	flag.StringVar(&flgConfig, "config", "/tmp/local-net.json", "config file path")
	flag.BoolVar(&flgDebug, "debug", false, "enable debug log")
	flag.IntVar(&flgInterval, "interval", 30, "interval between probing")

	flag.Parse()
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if flgDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	var currentWifiName string
	var agh *os.Process

	cfg, err := config.New(flgConfig)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Print(cfg)

	// Get Wifi device name
	device, err := nmcli.DeviceStatus("wifi")
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	aghUrl := fmt.Sprintf("http://%s", cfg.IpAddresses["dns"].IpAddress)
	aghcli, err := adguardhome.New(aghUrl, cfg.Dns.Credentials.Username, cfg.Dns.Credentials.Password)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	probe := func() {
		// Check if addresses are set
		for label, ipAddr := range cfg.IpAddresses {
			i, err := ip.New(device.Name, label, ipAddr.IpAddress, ipAddr.Netmask)
			if err != nil {
				log.Fatal().Err(err).Send()
			}

			if !i.IsSet() {
				log.Info().Msgf("set %s ip address", label)

				if err := i.Set(); err != nil {
					log.Error().Str("id", "ip-set").Err(err).Send()
					return
				}
			}
		}

		// Run dns service
		if agh == nil {
			agh, err = adguardhome.Run(cfg.IpAddresses["dns"].IpAddress)
			if err != nil {
				log.Error().Str("id", "dns-run").Err(err).Send()
				return
			}

			log.Info().Msg("AdGuardHome started")
			log.Debug().Msgf("AdGuardHome process ID is %d", agh.Pid)
		}

		// Network
		var vpn *nmcli.Connection

		// Check wifi connection
		wifi, err := nmcli.ConnectionShow("wifi", "")
		if err != nil {
			currentWifiName = ""

			log.Error().Str("id", "wifi-conn-check").Err(err).Send()
			return
		}

		if wifi == nil {
			currentWifiName = ""

			log.Info().Msg("wifi not connected")

			return
		}
		log.Print(wifi)

		if wifi.Name == currentWifiName {
			return
		}

		wifiTrusted := cfg.IsWifiTrusted(wifi.Name)
		log.Debug().Any("wifi-trust", wifiTrusted).Send()

		// Check vpn connection
		if cfg.Vpn.Enable {
			vpn, err = nmcli.ConnectionShow(cfg.Vpn.Type, cfg.Vpn.Name)
			if err == nil {
				log.Error().Str("id", "vpn-conn-check").Err(err).Send()
			}
		} else {
			log.Debug().Msg("vpn is not enabled")
		}
		log.Print(vpn)

		// Check if vpn needs to be disconnected
		if wifi == nil || wifiTrusted {
			if vpn != nil && vpn.Device != "" {
				if err := nmcli.ConnectionDown(cfg.Vpn.Name); err != nil {
					log.Error().Str("id", "vpn-disconn").Err(err).Send()
					return
				}
				log.Info().Msgf("disconnected vpn %s", cfg.Vpn.Name)
			}
		}

		// Connect to vpn if wifi is not trusted
		if !wifiTrusted {
			if vpn == nil {
				if err := nmcli.ConnectionUp(cfg.Vpn.Name); err != nil {
					log.Error().Str("id", "vpn-conn-wifi-no-trust").Err(err).Send()
					return
				}
				log.Info().Msgf("connected to vpn %s", cfg.Vpn.Name)
			}
		}

		// DNS
		// Update dns upstream server
		var vpnUuid string
		if vpn != nil {
			vpnUuid = vpn.Uuid
		}

		dnsServers, err := setDnsUpstreamServers(&cfg.Dns, aghcli, wifiTrusted, wifi.Uuid, vpnUuid)
		if err != nil {
			log.Error().Str("id", "dns-upstream-update").Err(err).Send()

			return
		}
		log.Info().Msgf("set dns servers to %v", dnsServers)

		currentWifiName = wifi.Name
	}

	log.Info().Msg("start probing")
	ticker := time.NewTicker(time.Duration(flgInterval) * time.Second)
	done := make(chan bool)

	probe()

	go func() {
		log.Info().Msg("Setting network")

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				probe()
			}
		}
	}()

	// Run container routine
	doneContainer := make(chan bool)

	if cfg.Dns.Container.Enable {
		go func() {
			log.Info().Msg("Listening to Docker's events")

			dockercli, err := docker.NewDockerClient()
			if err != nil {
				log.Fatal().Err(err).Send()
			}
			defer dockercli.Close()

			dockercli.Ping()

			messages, errs := dockercli.Events(events.ContainerEventType)

			for {
				select {
				case err := <-errs:
					if err != nil && err != io.EOF {
						log.Error().Str("id", "container-event").Err(err).Send()
					}

					os.Exit(1)
				case <-doneContainer:
					return
				case e := <-messages:
					if !funk.Contains([]string{"create", "destroy"}, e.Status) {
						continue
					}

					recordDomain, ok := e.Actor.Attributes[cfg.Dns.Container.LabelDomain]
					if !ok {
						continue
					}
					recordAnswer, ok := e.Actor.Attributes[cfg.Dns.Container.LabelAnswer]
					if !ok {
						continue
					}

					/*
						if err = validation.IsValidFQDN(domain); err != nil {
							log.Error(err)
							continue
						}
					*/

					record := adguardhome.Record{
						Domain: recordDomain,
						Answer: recordAnswer,
					}

					if e.Status == "create" {
						if err := aghcli.RewriteAdd(&record); err != nil {
							log.Error().Str("id", "container-event-create").Err(err).Send()
							return
						}

						log.Info().
							Str("container_id", e.ID).
							Str("dns_domain", recordDomain).
							Str("dns_answer", recordAnswer).
							Msg("Added dns record")
					}

					if e.Status == "destroy" {
						if err := aghcli.RewriteDelete(&record); err != nil {
							log.Error().Str("id", "container-event-destroy").Err(err).Send()
							return
						}

						log.Info().
							Str("container_id", e.ID).
							Str("dns_domain", recordDomain).
							Str("dns_answer", recordAnswer).
							Msg("Deleted dns record")
					}
				}
			}
		}()
	}

	// Exit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	agh.Kill()

	ticker.Stop()
	done <- true
	doneContainer <- true
}
