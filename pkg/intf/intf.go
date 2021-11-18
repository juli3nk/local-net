package intf

import (
	"fmt"
	"net"
)

func Get() error {
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	//iface, _ := net.InterfaceByName("enp1s0")

	for _, i := range ifaces {
		fmt.Printf("%+v\n", i)

		addrs, err := i.Addrs()
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", addrs)

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				/*
				if !v.IP.IsLoopback() {
						if v.IP.To4() != nil {//Verify if IP is IPV4
								ip = v.IP
						}
				}
				*/
				ip = v.IP
			case *net.IPAddr:
							ip = v.IP
			}

			fmt.Println(ip)
		}
	}

	return nil
}
