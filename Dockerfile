FROM alpine AS adguardhome

ARG ADGUARDHOME_VERSION

COPY get-adguardhome.sh /tmp/

RUN apk --update --no-cache add \
	bash \
	curl \
	git \
	tar

RUN /tmp/get-adguardhome.sh $ADGUARDHOME_VERSION


FROM golang:1-alpine AS builder

RUN apk --update add \
		ca-certificates \
		gcc \
		git \
		musl-dev

WORKDIR /go/src/github.com/juli3nk/local-net

ENV GO111MODULE off

COPY . .

RUN go get
RUN go build -ldflags "-linkmode external -extldflags -static -s -w" -o /tmp/local-net


FROM alpine

RUN apk --update --no-cache add \
	networkmanager-cli

COPY --from=adguardhome /usr/local/bin/AdGuardHome /usr/local/bin/
COPY --from=builder /tmp/local-net /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/local-net"]
