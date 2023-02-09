GOOS=windows GOARCH=amd64 go build -o bin/windows
GOOS=darwin GOARCH=amd64 go build -o bin/darwin
GOOS=linux GOARCH=amd64 go build -o bin/linux

GOOS=windows GOARCH=arm64 go build -o bin/arm64/windows/chatubrate-dvr &&
GOOS=darwin GOARCH=arm64 go build -o bin/arm64/darwin/chatubrate-dvr &&
GOOS=linux GOARCH=arm64 go build -o bin/arm64/linux/chatubrate-dvr
