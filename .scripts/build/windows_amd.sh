BIN="./.bin"

if [[ ! -d "$BIN" ]]; then
  mkdir -p "$BIN"
fi

GOOS=windows GOARCH=amd64 go build -o "$BIN/app.exe" ./cmd/app
