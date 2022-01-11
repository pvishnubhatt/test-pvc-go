GOARCH := amd64
GOOS := linux

all: build
local:
	go build -o bin/main app/main/main.go
	go build -o bin/get app/get/get.go
build: 
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/main app/main/main.go
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/get app/get/get.go
clean: ## Remove temporary files
	rm -f bin/main 
	rm -f bin/get
	go clean