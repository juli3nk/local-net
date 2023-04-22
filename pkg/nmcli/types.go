package nmcli

type Device struct {
	Name string
	Type string
	State string
	Connection string
}

type Connection struct {
	Name string
	Uuid string
	Type string
	Device string
}
