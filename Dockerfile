FROM alpine AS adguardhome

ARG ADGUARDHOME_VERSION

COPY get-adguardhome.sh /tmp/

RUN apk --update --no-cache add \
	bash \
	curl \
	tar

RUN /tmp/get-adguardhome.sh


FROM golang:1.17-alpine3.14 AS builder

RUN apk --update add \
		ca-certificates \
		gcc \
		git \
		musl-dev

WORKDIR /go/src/github.com/juli3nk/local-dns

ENV GO111MODULE off

COPY . .

RUN go get
RUN go build -ldflags "-linkmode external -extldflags -static -s -w" -o /tmp/local-dns


FROM alpine

RUN apk --update --no-cache add \
	networkmanager

COPY --from=adguardhome /usr/local/bin/AdGuardHome /usr/local/bin/
COPY --from=builder /tmp/local-dns /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/local-dns"]