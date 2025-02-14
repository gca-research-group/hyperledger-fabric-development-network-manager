FROM hyperledger/fabric-peer:latest

RUN apt-get update && apt-get install -y wget

ENV GO_VERSION=1.24.0
RUN wget https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz && \
    rm go$GO_VERSION.linux-amd64.tar.gz

ENV GOPATH=/opt/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

RUN go version