# Local Net

wifi
vpn

dns

https://ip.me

## Config example

```yaml
---
addresses:
  dns:
    ip_address: '192.168.10.1'
    netmask: '255.255.255.0'
  http:
    ip_address: '192.168.10.2'
    netmask: '255.255.255.0'
trusted:
  home: 'HomeWifi'
  phone: 'PhoneWifi'
vpn:
  enable: true
  name: 'vpn1'
dns:
  enable: true
  credentials:
    username: 'admin'
    password: 'admin'
  upstream_servers:
    default:
      - '1.1.1.1'
      - '9.9.9.9'
  container:
    enable: true
    label_name: 'dns.record'
```

The first time AdGuardhome needs to be configured at the following address: `http://192.168.10.1:3000`.
