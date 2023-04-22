# Local Net

Config example

```yaml
---
nic:
  label: 'dns'
  ip_address: '192.168.10.1'
  netmask: '255.255.255.0'
home:
  wifi_name: 'HomeWifi'
  dns_server: '192.168.1.1'
dns_provider:
  url: 'http://192.168.10.1'
  username: 'admin'
  password: 'admin'
dns_servers:
  - '1.1.1.1'
  - '9.9.9.9'
```

http://192.168.10.1:3000
