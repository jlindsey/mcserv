## mcserv
## Copyright (C) 2017 Joshua Lindsey
##
## This program is free software: you can redistribute it and/or modify
## it under the terms of the GNU Lesser General Public License as published by
## the Free Software Foundation, either version 3 of the License, or
## (at your option) any later version.
##
## This program is distributed in the hope that it will be useful,
## but WITHOUT ANY WARRANTY; without even the implied warranty of
## MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
## GNU Lesser General Public License for more details.
##
## You should have received a copy of the GNU Lesser General Public License
## along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
