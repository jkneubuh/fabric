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

# Just copy the release build outputs for the target arch
COPY    release/${TARGETOS}-${TARGETARCH}/bin/orderer /usr/local/bin

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
