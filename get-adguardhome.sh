#!/usr/bin/env bash

REPO_URL="https://github.com/AdguardTeam/AdGuardHome"

LATEST_VERSION=$(git ls-remote --refs --tags "$REPO_URL" | awk '{ print $2 }' | sed 's#refs/tags/##' | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+$" | sort -V | tail -n 1)
VERSION="${1:-${LATEST_VERSION}}"


echo -e "Downloading AdGuardHome version ${VERSION}"

curl -sfL "${REPO_URL}/releases/download/${VERSION}/AdGuardHome_linux_amd64.tar.gz" | tar -xzC /tmp
mv /tmp/AdGuardHome/AdGuardHome /usr/local/bin/
rm -rf /tmp/AdGuardHome
