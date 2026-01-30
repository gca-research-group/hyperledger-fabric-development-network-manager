FROM golang:alpine

WORKDIR /home/hfndm/app

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
COPY ./.docker/hfndm/start.sh ./start.sh

RUN go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go get -u github.com/golang-migrate/migrate/v4/database/postgres

RUN GOOS=linux GOARCH=arm go build -o "./.bin/app" ./cmd/app

ENTRYPOINT [ "sh", "/home/hfndm/app/start.sh" ]
