SHELL := bash

ifeq ($(OS),Windows_NT)     # is Windows_NT on XP, 2000, 7, Vista, 10...
	_OS := windows
	ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
		_ARCH := amd64
	else ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
		_ARCH := amd64
	else
		_ARCH := arm64
	endif
else
	_OS := $(shell uname | tr '[:upper:]' '[:lower:]')
	UNAME_M := $(shell uname -m)
	ifeq ($(UNAME_M),x86_64)
		_ARCH := amd64
	else
		_ARCH := arm64
	endif
endif

ifeq ($(_OS),windows)
	BINEXT := .exe
endif

.PHONY: build
build:
	$(MAKE) dist/$(_OS)/$(_ARCH)/waitsilence$(BINEXT)


PHONY: build-all
build-all: dist/linux/amd64/waitsilence dist/linux/arm64/waitsilence dist/darwin/amd64/waitsilence dist/darwin/arm64/waitsilence dist/windows/amd64/waitsilence.exe dist/windows/arm64/waitsilence.exe


dist/linux/amd64/waitsilence:
	docker run --rm -v '$(CURDIR)':/usr/src/app -w /usr/src/app \
		-e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=0 \
		--entrypoint sh \
		golang:1.17-buster \
		-c "go get -d -v && go build -a -v -o dist/linux/amd64/waitsilence"


dist/linux/arm64/waitsilence:
	docker run --rm -v '$(CURDIR)':/usr/src/app -w /usr/src/app \
		-e GOOS=linux -e GOARCH=arm64 -e CGO_ENABLED=0 \
		--entrypoint sh \
		golang:1.17-buster \
		-c "go get -d -v && go build -a -v -o dist/linux/arm64/waitsilence"


dist/darwin/amd64/waitsilence:
	docker run --rm -v '$(CURDIR)':/usr/src/app -w /usr/src/app \
		-e GOOS=darwin -e GOARCH=amd64 -e CGO_ENABLED=0 \
		--entrypoint sh \
		golang:1.17-buster \
		-c "go get -d -v && go build -a -v -o dist/darwin/amd64/waitsilence"


dist/darwin/arm64/waitsilence:
	docker run --rm -v '$(CURDIR)':/usr/src/app -w /usr/src/app \
		-e GOOS=darwin -e GOARCH=arm64 -e CGO_ENABLED=0 \
		--entrypoint sh \
		golang:1.17-buster \
		-c "go get -d -v && go build -a -v -o dist/darwin/arm64/waitsilence"


dist/windows/amd64/waitsilence.exe:
	docker run --rm -v '$(CURDIR)':/usr/src/app -w /usr/src/app \
		-e GOOS=windows -e GOARCH=amd64 -e CGO_ENABLED=0 \
		--entrypoint sh \
		golang:1.17-buster \
		-c "go get -d -v && go build -a -v -o dist/windows/amd64/waitsilence.exe"


dist/windows/arm64/waitsilence.exe:
	docker run --rm -v '$(CURDIR)':/usr/src/app -w /usr/src/app \
		-e GOOS=windows -e GOARCH=arm64 -e CGO_ENABLED=0 \
		--entrypoint sh \
		golang:1.17-buster \
		-c "go get -d -v && go build -a -v -o dist/windows/arm64/waitsilence.exe"


PHONY: clean
clean:
	rm -rf '$(CURDIR)'/dist/*
