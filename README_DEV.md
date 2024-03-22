Compile All at once:

```
GOOS=windows GOARCH=amd64 go build -o bin/windows/chatubrate-dvr.exe &&
GOOS=darwin GOARCH=amd64 go build -o bin/darwin/chatubrate-dvr &&
GOOS=linux GOARCH=amd64 go build -o bin/linux/chatubrate-dvr &&
GOOS=windows GOARCH=arm64 go build -o bin/arm64/windows/chatubrate-dvr.exe &&
GOOS=darwin GOARCH=arm64 go build -o bin/arm64/darwin/chatubrate-dvr &&
GOOS=linux GOARCH=arm64 go build -o bin/arm64/linux/chatubrate-dvr
```

or Compile for 64-bit Windows, macOS, Linux:

```
GOOS=windows GOARCH=amd64 go build -o bin/windows/chatubrate-dvr.exe &&
GOOS=darwin GOARCH=amd64 go build -o bin/darwin/chatubrate-dvr &&
GOOS=linux GOARCH=amd64 go build -o bin/linux/chatubrate-dvr
```

or for arm64 Windows, macOS, Linux:

```
GOOS=windows GOARCH=arm64 go build -o bin/arm64/windows/chatubrate-dvr.exe &&
GOOS=darwin GOARCH=arm64 go build -o bin/arm64/darwin/chatubrate-dvr &&
GOOS=linux GOARCH=arm64 go build -o bin/arm64/linux/chatubrate-dvr
```
