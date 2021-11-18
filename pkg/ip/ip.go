package ip

import (
	"fmt"

	uip "github.com/juli3nk/go-utils/ip"
	"github.com/thoas/go-funk"
	"github.com/vishvananda/netlink"
)

type IPConfig struct {
	NIC     string
	Label   string
	IP      string
	Netmask string
}

func New(nic, label, ip, netmask string) (*IPConfig, error) {
	c := IPConfig{
		NIC:     nic,
		Label:   label,
		IP:      ip,
		Netmask: netmask,
	}

	return &c, nil
}

func (c *IPConfig) Ping() bool {
	// TODO
	return false
}

func (c *IPConfig) IsSet() bool {
	interfaces := uip.New()
	interfaces.Get()

	if _, ok := interfaces[c.NIC]; !ok {
		return false
	}

	intf := interfaces.GetIntf(c.NIC)

	if funk.Contains(intf.V4, c.IP) {
		return true
	}

	return false
}

func (c *IPConfig) Set() error {
	ones := ConvertNetmaskToCIDR(c.Netmask)

	ip := fmt.Sprintf("%s/%d", c.IP, ones)

	intf, err := netlink.LinkByName(c.NIC)
	if err != nil {
		return err
	}

	addr, err := netlink.ParseAddr(ip)
	if err != nil {
		return err
	}

	label := fmt.Sprintf("%s:%s", c.NIC, c.Label)
	addr.Label = label

	netlink.AddrAdd(intf, addr)

	return nil
}

func (c *IPConfig) Unset() error {
	ones := ConvertNetmaskToCIDR(c.Netmask)

	ip := fmt.Sprintf("%s/%d", c.IP, ones)

	intf, err := netlink.LinkByName(c.NIC)
	if err != nil {
		return err
	}

	addr, err := netlink.ParseAddr(ip)
	if err != nil {
		return err
	}

	netlink.AddrDel(intf, addr)

	return nil
}
