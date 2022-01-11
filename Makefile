GOARCH := amd64
GOOS := linux

all: build
local:
	go build -o dist/main app/main/main.go
	go build -o dist/get app/get/get.go
build: 
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o dist/main app/main/main.go
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o dist/get app/get/get.go
clean: ## Remove temporary files
	rm -f dist/main 
	rm -f dist/get
	go clean