package config

type Config struct {
	IpAddresses map[string]IpAddress `json:"ip_addresses"`
	Wifi   map[string]string  `json:"wifi"`
	Vpn       Vpn                `json:"vpn"`
	Domain    string                `json:"domain"`
	Dns       Dns                `json:"dns"`
}

type IpAddress struct {
	IpAddress string `json:"ip_address"`
	Netmask   string `json:"netmask"`
}

type Vpn struct {
	Enable bool   `json:"enable"`
	Name   string `json:"name"`
}

type Dns struct {
	Enable          bool            `json:"enable"`
	Credentials     Credentials     `json:"credentials"`
	UpstreamServers UpstreamServers `json:"upstream_servers"`
	Container       Container       `json:"container"`
}

type Credentials struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpstreamServers struct {
	Default   []string            `json:"default"`
	locations map[string]Location `json:"locations"`
}

type Location struct {
	WifiName string `json:"wifi_name"`
	Server   string `json:"dns_server"`
}

type Container struct {
	Enable    bool   `json:"enable"`
	LabelDomain string `json:"label_domain"`
	LabelAnswer string `json:"label_answer"`
}
