Compile All at once:

```
GOOS=windows GOARCH=amd64 go build -o bin/windows_chatubrate-dvr.exe &&
GOOS=darwin GOARCH=amd64 go build -o bin/macos_chatubrate-dvr &&
GOOS=linux GOARCH=amd64 go build -o bin/linux_chatubrate-dvr &&
GOOS=windows GOARCH=arm64 go build -o bin/arm64_windows_chatubrate-dvr.exe &&
GOOS=darwin GOARCH=arm64 go build -o bin/arm64_macos_chatubrate-dvr &&
GOOS=linux GOARCH=arm64 go build -o bin/arm64_linux_chatubrate-dvr
```

or Compile for 64-bit Windows, macOS, Linux:

```
GOOS=windows GOARCH=amd64 go build -o bin/windows_chatubrate-dvr.exe &&
GOOS=darwin GOARCH=amd64 go build -o bin/macos_chatubrate-dvr &&
GOOS=linux GOARCH=amd64 go build -o bin/linux_chatubrate-dvr
```

or for arm64 Windows, macOS, Linux:

```
GOOS=windows GOARCH=arm64 go build -o bin/arm64_windows_chatubrate-dvr.exe &&
GOOS=darwin GOARCH=arm64 go build -o bin/arm64_macos_chatubrate-dvr &&
GOOS=linux GOARCH=arm64 go build -o bin/arm64_linux_chatubrate-dvr
```
