BIN="./.bin"

if [[ ! -d "$BIN" ]]; then
  mkdir -p "$BIN"
fi

GOOS=linux GOARCH=arm64 go build -o "$BIN/fno_linux_arm64" ./cmd/cli
