
REPODIR := "/go/src/github.com/juli3nk/local-net"
IMAGE_NAME := juli3nk/local-net

.PHONY: dev
dev:
	docker container run \
		-ti \
		--rm \
		--mount type=bind,src=$$PWD,dst=${REPODIR} \
		--mount type=bind,src=/var/run/dbus,dst=/var/run/dbus \
		--mount type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
		--net host \
		--dns 1.1.1.1 \
		--security-opt seccomp=unconfined \
		--security-opt apparmor=unconfined \
		--cap-add=NET_ADMIN \
		--env DOCKER_API_VERSION=1.41 \
		--workdir ${REPODIR} \
		--name local-net_dev \
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
		--mount type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
		--mount type=bind,src=$$HOME/.config/local/net.json,dst=/tmp/local-net.json \
		--mount type=bind,src=$$HOME/Data/adguardhome/work,dst=/opt/adguardhome/work \
		--mount type=bind,src=$$HOME/Data/adguardhome/conf,dst=/opt/adguardhome/conf \
		--net host \
		--dns 1.1.1.1 \
		--security-opt seccomp=unconfined \
		--security-opt apparmor=unconfined \
		--cap-add=NET_ADMIN \
		--env DOCKER_API_VERSION=1.41 \
		--name local-net \
		juli3nk/local-net \
			-debug

.PHONY: logs
logs:
	@docker container logs local-net -f

.PHONY: ip
ip:
	@ip a

.PHONY: check
check:
	@dig $$(awk '/nameserver/ { print $$2 }' /etc/resolv.conf) github.com
