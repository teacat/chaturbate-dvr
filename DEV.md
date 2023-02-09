GOOS=windows GOARCH=amd64 go build -o bin/windows
GOOS=darwin GOARCH=amd64 go build -o bin/darwin
GOOS=linux GOARCH=amd64 go build -o bin/linux

GOOS=windows GOARCH=arm64 go build -o bin/arm64/windows
GOOS=darwin GOARCH=arm64 go build -o bin/arm64/darwin
GOOS=linux GOARCH=arm64 go build -o bin/arm64/linux

GOOS=windows GOARCH=arm go build -o bin/arm/windows
GOOS=darwin GOARCH=arm go build -o bin/arm/darwin
GOOS=linux GOARCH=arm go build -o bin/arm/linux
