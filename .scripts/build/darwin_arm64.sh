BIN="./.bin"

if [[ ! -d "$BIN" ]]; then
  mkdir -p "$BIN"
fi

GOOS=darwin GOARCH=arm64 go build -o "$BIN/fno_darwin_arm64" ./cmd/cli