GOOS=windows GOARCH=amd64 go build -o bin/windows
GOOS=darwin GOARCH=amd64 go build -o bin/darwin
GOOS=linux GOARCH=amd64 go build -o bin/linux
