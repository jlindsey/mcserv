SOURCES := $(filter-out $(wildcard *_test.go), $(wildcard *.go))
XC_OS?=linux darwin
XC_ARCH?=amd64

all: bin

build:
	mkdir -p build

bin: vendor ${SOURCES} | build
	gox -os="${XC_OS}" -arch="${XC_ARCH}" -output="build/{{.OS}}/{{.Arch}}/mcserv"

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
