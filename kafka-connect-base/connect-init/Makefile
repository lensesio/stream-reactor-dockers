.PHONY: build clean deploy

# Use for go build binar
SRVER?=

build:
	@export SRVER=Dev; env GOOS=linux go build -ldflags="-s -w -X main.SRVER=$${SRVER}" -o ./connect-init ./main.go

test:
	@rm -rf /connect-init
	@export SRVER=Dev; env GOOS=linux go build -ldflags="-s -w -X main.SRVER=$${SRVER}" -o ./connect-init ./main.go
	@./connect-init
