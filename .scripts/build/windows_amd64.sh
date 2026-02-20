BIN="./.bin"

if [[ ! -d "$BIN" ]]; then
  mkdir -p "$BIN"
fi

GOOS=windows GOARCH=amd64 go build -o "$BIN/fno_windows_amd64.exe" ./cmd/cli
