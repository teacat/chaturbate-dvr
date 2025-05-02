64-bit + arm64

```
GOOS=windows GOARCH=amd64 go build -o bin/x64_windows_chatubrate-dvr.exe &&
GOOS=darwin GOARCH=amd64 go build -o bin/x64_macos_chatubrate-dvr &&
GOOS=linux GOARCH=amd64 go build -o bin/x64_linux_chatubrate-dvr &&
GOOS=windows GOARCH=arm64 go build -o bin/arm64_windows_chatubrate-dvr.exe &&
GOOS=darwin GOARCH=arm64 go build -o bin/arm64_macos_chatubrate-dvr &&
GOOS=linux GOARCH=arm64 go build -o bin/arm64_linux_chatubrate-dvr
```

64-bit Windows, macOS, Linux:

```
GOOS=windows GOARCH=amd64 go build -o bin/x64_windows_chatubrate-dvr.exe &&
GOOS=darwin GOARCH=amd64 go build -o bin/x64_macos_chatubrate-dvr &&
GOOS=linux GOARCH=amd64 go build -o bin/x64_linux_chatubrate-dvr
```

arm64 Windows, macOS, Linux:

```
GOOS=windows GOARCH=arm64 go build -o bin/arm64_windows_chatubrate-dvr.exe &&
GOOS=darwin GOARCH=arm64 go build -o bin/arm64_macos_chatubrate-dvr &&
GOOS=linux GOARCH=arm64 go build -o bin/arm64_linux_chatubrate-dvr
```

Build Docker Tag:
s

```
docker build -t yamiodymel/chaturbate-dvr:2.0.0 .
docker push yamiodymel/chaturbate-dvr:2.0.0
```
