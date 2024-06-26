# build file
GO_BUILD=go build -ldflags -s -v

BIN_BINARY_NAME=mybtckb-svr-bin
mybtckb-svr:
	GO111MODULE=on $(GO_BUILD) -o $(BIN_BINARY_NAME) cmd/main.go
	@echo "Build $(BIN_BINARY_NAME) successfully. You can run ./$(BIN_BINARY_NAME) now.If you can't see it soon,wait some seconds"

update:
	go mod tidy
