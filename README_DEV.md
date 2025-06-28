64-bit + arm64

```
GOOS=windows GOARCH=amd64 go build -o bin/x64_windows_chaturbate-dvr.exe &&
GOOS=darwin GOARCH=amd64 go build -o bin/x64_macos_chaturbate-dvr &&
GOOS=linux GOARCH=amd64 go build -o bin/x64_linux_chaturbate-dvr &&
GOOS=windows GOARCH=arm64 go build -o bin/arm64_windows_chaturbate-dvr.exe &&
GOOS=darwin GOARCH=arm64 go build -o bin/arm64_macos_chaturbate-dvr &&
GOOS=linux GOARCH=arm64 go build -o bin/arm64_linux_chaturbate-dvr
```

64-bit Windows, macOS, Linux:

```
GOOS=windows GOARCH=amd64 go build -o bin/x64_windows_chaturbate-dvr.exe &&
GOOS=darwin GOARCH=amd64 go build -o bin/x64_macos_chaturbate-dvr &&
GOOS=linux GOARCH=amd64 go build -o bin/x64_linux_chaturbate-dvr
```

arm64 Windows, macOS, Linux:

```
GOOS=windows GOARCH=arm64 go build -o bin/arm64_windows_chaturbate-dvr.exe &&
GOOS=darwin GOARCH=arm64 go build -o bin/arm64_macos_chaturbate-dvr &&
GOOS=linux GOARCH=arm64 go build -o bin/arm64_linux_chaturbate-dvr
```

Build Docker Tag:

```
docker build -t yamiodymel/chaturbate-dvr:2.0.0 .
docker push yamiodymel/chaturbate-dvr:2.0.0
docker image tag yamiodymel/chaturbate-dvr:2.0.0 yamiodymel/chaturbate-dvr:latest
docker push yamiodymel/chaturbate-dvr:latest
```
