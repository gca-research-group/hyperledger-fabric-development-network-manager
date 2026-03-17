FROM ubuntu:22.04 AS builder

ARG FABRIC_VERSION=3.1.4

RUN apt-get update && apt-get install -y wget tar ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /tmp

RUN wget -qO fabric.tar.gz https://github.com/hyperledger/fabric/releases/download/v${FABRIC_VERSION}/hyperledger-fabric-linux-amd64-${FABRIC_VERSION}.tar.gz \
    && mkdir fabric \
    && tar -xzf fabric.tar.gz -C fabric

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y ca-certificates bash \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /tmp/fabric/bin /opt/fabric/bin

COPY --from=builder /tmp/fabric/config /etc/hyperledger/fabric

ENV PATH="/opt/fabric/bin:${PATH}"
ENV FABRIC_CFG_PATH=/etc/hyperledger/fabric

WORKDIR /workspace

CMD ["sleep", "infinity"]