VERSION = $(shell cat ../VERSION)

build:
	@packr2 clean
	@packr2 build -i -ldflags "-X main.Version=${VERSION}" -o neo-local .
	@packr2 clean

test:
	@go test ./... -cover