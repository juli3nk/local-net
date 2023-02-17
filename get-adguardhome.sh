#!/usr/bin/env bash

ADGUARDHOME_VERSION=${1:-0.107.2}

curl -sfL https://github.com/AdguardTeam/AdGuardHome/releases/download/v${ADGUARDHOME_VERSION}/AdGuardHome_linux_amd64.tar.gz | tar -xzC /tmp
mv /tmp/AdGuardHome/AdGuardHome /usr/local/bin/
rm -rf /tmp/AdGuardHome
