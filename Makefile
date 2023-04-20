
REPODIR := "/go/src/github.com/juli3nk/local-dns"
IMAGE_NAME := juli3nk/local-dns

.PHONY: dev
dev:
	docker container run \
		-ti \
		--rm \
		--mount type=bind,src=$$PWD,dst=${REPODIR} \
		--mount type=bind,src=/var/run/dbus,dst=/var/run/dbus \
		--net host \
		--dns 1.1.1.1 \
		--cap-add=NET_ADMIN \
		-w ${REPODIR} \
		--name local-dns_dev \
		juli3nk/dev:go

.PHONY: build
build:
	docker image build \
		-t ${IMAGE_NAME} \
		.

.PHONY: run
run:
	@mkdir -p $$HOME/Data/adguardhome/{work,conf}
	docker container run \
		-d \
		--rm \
		--mount type=bind,src=/var/run/dbus,dst=/var/run/dbus \
		--mount type=bind,src=$$HOME/.config/local-dns/config.yml,dst=/tmp/local-dns.yml \
		--mount type=bind,src=$$HOME/Data/adguardhome/work,dst=/opt/adguardhome/work \
		--mount type=bind,src=$$HOME/Data/adguardhome/conf,dst=/opt/adguardhome/conf \
		--net host \
		--dns 1.1.1.1 \
		--cap-add=NET_ADMIN \
		--name local-dns \
		juli3nk/local-dns \
			-debug

.PHONY: logs
logs:
	@docker container logs local-dns -f

.PHONY: ip
ip:
	@ip a

.PHONY: check
check:
	@dig @192.168.82.1 github.com
