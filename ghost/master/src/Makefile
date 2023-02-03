DLL_NAME := shiori.dll
GO_VERSION := 1.19 # or "latest" for bleeding-edge Go version
GOCACHE ?= $(shell go env GOCACHE)
GOPATH ?= $(shell go env GOPATH)

.PHONY:all
all:
	make build
	make deploy
	make clean

.PHONY: clean
clean:
	rm -vf libshiori.h libshiori.a shiori.h

.PHONY: deploy
deploy:
	mv -vf $(DLL_NAME) ..

.PHONY: build
build:
	@echo ":: Building $(shell basename $(PWD)) in MinGW container..."
	docker run --rm -it \
		-u "$${UID}:$${GID}" \
		-e GOCACHE=/go/.cache \
		-e GOARCH=386 \
		-v "$(GOCACHE)":/go/.cache \
		-v "$(GOPATH)/src":/go/src \
		-v "$(GOPATH)/pkg":/go/pkg \
		-v "$(PWD)":/go/work \
		-w /go/work \
		x1unix/go-mingw:$(GO_VERSION) \
		/bin/sh -c '\
		go build -x -buildmode=c-archive -o libshiori.a .;\
		i686-w64-mingw32-gcc-win32 -shared -o $(DLL_NAME) shiori.def libshiori.a -Wl,--allow-multiple-definition -static -lstdc++ -lwinmm -lntdll -lws2_32;\
		'
