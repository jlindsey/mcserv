SOURCES := $(filter-out $(wildcard *_test.go), $(wildcard *.go))

all: bin

build:
	mkdir -p build

bin: vendor ${SOURCES} | build
	gox -os="darwin" -arch="amd64" -output="build/{{.OS}}/mcserv"

vendor: glide.lock glide.yaml
	glide i

clean:
	rm -rf build

distclean: | clean
	rm -rf vendor

test: | vendor
	go test
	rm -f build/test

.PHONY: all bin clean distclean test
