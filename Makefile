##
## ----------------------------------------------------------------------------
##   Docker
## ----------------------------------------------------------------------------
##

docker-start-chaturbate-dvr: ## Start docker Chaturbate DVR
	docker-compose -f docker-compose.yml up -d

docker-stop-chaturbate-dvr: ## Stop docker Chaturbate DVR
	docker-compose -f docker-compose.yml down

docker-restart-chaturbate-dvr: ## Restart project Chaturbate DVR
	docker-compose -f docker-compose.yml restart

docker-start-chaturbate-dvr-web: ## Start docker Chaturbate DVR WEB
	docker-compose -f docker-compose-web.yml up -d

docker-stop-chaturbate-dvr-web: ## Stop docker Chaturbate DVR WEB
	docker-compose -f docker-compose-web.yml down

docker-restart-chaturbate-dvr-web: ## Restart project Chaturbate DVR WEB
	docker-compose -f docker-compose-web.yml restart

.PHONY: docker-start-chaturbate-dvr docker-stop-chaturbate-dvr docker-restart-chaturbate-dvr docker-start-chaturbate-dvr-web docker-stop-chaturbate-dvr-web docker-restart-chaturbate-dvr-web

##
## ----------------------------------------------------------------------------
##   Compile
## ----------------------------------------------------------------------------
##

64bit-windows-macos-linux: ## Compile all arch amd64
	GOOS=windows GOARCH=amd64 go build -o bin/windows/chatubrate-dvr.exe && \
    GOOS=darwin GOARCH=amd64 go build -o bin/darwin/chatubrate-dvr && \
    GOOS=linux GOARCH=amd64 go build -o bin/linux/chatubrate-dvr

arm64-windows-macos-linux: ## Compile all arch arm64
	GOOS=windows GOARCH=arm64 go build -o bin/arm64/windows/chatubrate-dvr.exe && \
    GOOS=darwin GOARCH=arm64 go build -o bin/arm64/darwin/chatubrate-dvr && \
    GOOS=linux GOARCH=arm64 go build -o bin/arm64/linux/chatubrate-dvr

compile-all: ## Compile all
	GOOS=windows GOARCH=amd64 go build -o bin/windows/chatubrate-dvr.exe && \
    GOOS=darwin GOARCH=amd64 go build -o bin/darwin/chatubrate-dvr && \
    GOOS=linux GOARCH=amd64 go build -o bin/linux/chatubrate-dvr && \
    GOOS=windows GOARCH=arm64 go build -o bin/arm64/windows/chatubrate-dvr.exe && \
    GOOS=darwin GOARCH=arm64 go build -o bin/arm64/darwin/chatubrate-dvr && \
    GOOS=linux GOARCH=arm64 go build -o bin/arm64/linux/chatubrate-dvr

.PHONY: 64bit-windows-macos-linux arm64-windows-macos-linux

##
## ----------------------------------------------------------------------------
##   Help
## ----------------------------------------------------------------------------
##

.DEFAULT_GOAL := help
.PHONY: help
help: ## Show this help
	@egrep -h '(^[a-zA-Z0-9_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' \
		| sed -e 's/\[32m##/[33m/'