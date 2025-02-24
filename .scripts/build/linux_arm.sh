BIN="./.bin"

if [[ ! -d "$BIN" ]]; then
  mkdir -p "$BIN"
fi

GOOS=linux GOARCH=arm go build -o "$BIN/app" ./cmd/app
