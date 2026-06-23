BIN="./.bin"

if [[ ! -d "$BIN" ]]; then
  mkdir -p "$BIN"
fi

GOOS=linux GOARCH=amd64 go build -o "$BIN/fno_linux_amd64" ./cmd/cli
