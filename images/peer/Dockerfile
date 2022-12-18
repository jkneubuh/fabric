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
ARG K8S_BUILDER_VER=v0.9.3

# Just copy the release build outputs for the target arch
COPY release/${TARGETOS}-${TARGETARCH}/bin/peer             /usr/local/bin
COPY release/${TARGETOS}-${TARGETARCH}/builders/ccaas/bin   /opt/hyperledger/ccaas_builder

# Install the chaincode builder
RUN mkdir -p /opt/hyperledger/k8s_builder/bin
RUN curl -sL https://github.com/hyperledger-labs/fabric-builder-k8s/releases/download/${K8S_BUILDER_VER}/fabric-builder-k8s-${K8S_BUILDER_VER}-${TARGETOS}-${TARGETARCH}.tgz \
     | tar zxvf - \
        -C /opt/hyperledger/k8s_builder/bin

ENV FABRIC_CFG_PATH         /var/hyperledger/fabric/config
ENV FABRIC_VER              ${FABRIC_VER}

COPY sampleconfig/msp       ${FABRIC_CFG_PATH}/msp
COPY sampleconfig/core.yaml ${FABRIC_CFG_PATH}/core.yaml

VOLUME /etc/hyperledger/fabric
VOLUME /var/hyperledger

EXPOSE 7051

ENTRYPOINT [ "peer" ]
CMD [ "node", "start" ]
