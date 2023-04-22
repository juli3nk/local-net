package adguardhome

type DnsCfg struct {
	url      string
	username string
	password string
}

func New(url, username, password string) (*DnsCfg, error) {
	cfg := DnsCfg{
		url:      url,
		username: username,
		password: password,
	}

	return &cfg, nil
}
