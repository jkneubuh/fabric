# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

ARG GO_VER
ARG ALPINE_VER
FROM alpine:${ALPINE_VER} as peer-base
RUN apk add --no-cache tzdata
# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/golang/go/blob/go1.9.1/src/net/conf.go#L194-L275
# - docker run --rm debian:stretch grep '^hosts:' /etc/nsswitch.conf
RUN echo 'hosts: files dns' > /etc/nsswitch.conf

FROM golang:${GO_VER}-alpine${ALPINE_VER} as golang
RUN apk add --no-cache \
	bash \
	binutils-gold \
	gcc \
	git \
	make \
	musl-dev
ADD . $GOPATH/src/github.com/hyperledger/fabric
WORKDIR $GOPATH/src/github.com/hyperledger/fabric

FROM golang as peer
ARG TARGETARCH
ARG TARGETOS
ARG GO_TAGS
ARG FABRIC_VER
ARG CCBUILDER_VER
ENV FABRIC_VER ${FABRIC_VER}
RUN make peer GO_TAGS=${GO_TAGS}
RUN make ccaasbuilder

RUN mkdir -p release/${TARGETOS}-${TARGETARCH}/builders/k8s/bin
RUN wget -qO - \
    https://github.com/hyperledger-labs/fabric-builder-k8s/releases/download/v${CCBUILDER_VER}/fabric-builder-k8s-v${CCBUILDER_VER}-${TARGETOS}-${TARGETARCH}.tgz \
    | tar zxvf - -C release/${TARGETOS}-${TARGETARCH}/builders/k8s/bin

FROM peer-base
ARG TARGETARCH
ARG TARGETOS
ENV FABRIC_CFG_PATH /etc/hyperledger/fabric
VOLUME /etc/hyperledger/fabric
VOLUME /var/hyperledger
COPY --from=peer /go/src/github.com/hyperledger/fabric/build/bin /usr/local/bin
COPY --from=peer /go/src/github.com/hyperledger/fabric/sampleconfig/msp ${FABRIC_CFG_PATH}/msp
COPY --from=peer /go/src/github.com/hyperledger/fabric/sampleconfig/core.yaml ${FABRIC_CFG_PATH}/core.yaml
COPY --from=peer /go/src/github.com/hyperledger/fabric/release/${TARGETOS}-${TARGETARCH}/builders/ccaas/bin/ /opt/hyperledger/ccaas_builder/bin/
COPY --from=peer /go/src/github.com/hyperledger/fabric/release/${TARGETOS}-${TARGETARCH}/builders/k8s/bin/ /opt/hyperledger/k8s_builder/bin/
EXPOSE 7051
CMD ["peer","node","start"]
