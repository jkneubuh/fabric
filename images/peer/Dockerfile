#
# SPDX-License-Identifier: Apache-2.0
#

# TODO: why is this image 50MB to large?  It should be 48MB not 90MB


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
        bin/peer \
        builders/ccaas

RUN mv bin/peer             /usr/local/bin
RUN mkdir -p                /opt/hyperledger/ccaas_builder
RUN mv builders/ccaas/bin   /opt/hyperledger/ccaas_builder

ENV FABRIC_CFG_PATH         /var/hyperledger/fabric/config
ENV FABRIC_VER              ${FABRIC_VER}

COPY sampleconfig/msp       ${FABRIC_CFG_PATH}/msp
COPY sampleconfig/core.yaml ${FABRIC_CFG_PATH}/core.yaml

VOLUME /etc/hyperledger/fabric
VOLUME /var/hyperledger

EXPOSE 7051

ENTRYPOINT [ "peer" ]
CMD [ "node", "start" ]
