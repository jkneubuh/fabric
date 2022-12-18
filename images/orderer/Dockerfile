# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

ARG ALPINE_VER
FROM alpine:${ALPINE_VER}

RUN apk add --no-cache \
    bash \
    curl \
    tzdata

ARG TARGETOS
ARG TARGETARCH
ARG FABRIC_VER

# Really simple docker build: install the statically linked binary from the multi-arch distribution
# (TODO: or from the local build context if building locally)
RUN curl -sL \
    https://github.com/hyperledger/fabric/releases/download/v${FABRIC_VER}/hyperledger-fabric-${TARGETOS}-${TARGETARCH}-${FABRIC_VER}.tar.gz \
    | tar zxvf - \
        -C /usr/local bin/orderer

ENV     FABRIC_CFG_PATH            /var/hyperledger/fabric/config
ENV     FABRIC_VER                 ${FABRIC_VER}

COPY    sampleconfig/msp           ${FABRIC_CFG_PATH}/msp
COPY    sampleconfig/orderer.yaml  ${FABRIC_CFG_PATH}
COPY    sampleconfig/configtx.yaml ${FABRIC_CFG_PATH}

VOLUME  /etc/hyperledger/fabric
VOLUME  /var/hyperledger
EXPOSE  7050

ENTRYPOINT  [ "orderer" ]
CMD         [ "start" ]
