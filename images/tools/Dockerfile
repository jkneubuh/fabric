# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

ARG ALPINE_VER
FROM alpine:${ALPINE_VER}

RUN apk add --no-cache \
    bash \
    curl \
    jq \
    tzdata

ARG TARGETOS
ARG TARGETARCH
ARG FABRIC_VER

RUN REV=$(echo $FABRIC_VER | sed -e  's/^v\(.*\)/\1/') && \
    curl -sL https://github.com/hyperledger/fabric/releases/download/${FABRIC_VER}/hyperledger-fabric-${TARGETOS}-${TARGETARCH}-${REV}.tar.gz \
    | tar zxvf - \
        -C /usr/local bin

ENV     FABRIC_VER          ${FABRIC_VER}
ENV     FABRIC_CFG_PATH     /var/hyperledger/fabric/config

COPY    sampleconfig        ${FABRIC_CFG_PATH}

VOLUME  /etc/hyperledger/fabric
VOLUME  /var/hyperledger

ENTRYPOINT [ "bash" ]
