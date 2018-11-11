# default target
build:
	cd src; GOPATH=`realpath $$(pwd)/../vendor` go build -tags debug; cp src.exe ../../notekeeper-electron-frontend/app/resources/backend.exe

# copy content of vendor/* into vendor/src/*
# build can't happen until this is done
VENDOR_DIRS := $(wildcard vendor/*)
CLEAN_VENDOR_DIRS := $(filter-out vendor/src vendor/manifest,$(VENDOR_DIRS))
vendor-deps:
	$(foreach dir,$(CLEAN_VENDOR_DIRS),cp -r $(dir) vendor/src/;)

test:
	cd src; go test ./... -v -race

coverage-dep:
	go get -u github.com/wadey/gocovmerge

# https://github.com/golang/go/issues/6909
coverage:
	shell ./coverage.sh

# create a release build
# The -s ldflag strips symbol table & debug info
release:
	cd src; go build -ldflags "-s"

gvt:
	go get -u github.com/FiloSottile/gvt

loc:
	cd src; find . -name '*.go' -not -path './proto/*' | xargs wc -l

# install the protoc golang plugin
# export GOPATH=C:/code/go
proto-gen:
	go get -u github.com/golang/protobuf/protoc-gen-go

# rebuild the proto definitions
# export PATH=$PATH:/cygdrive/c/code/go/bin 
# export PATH=$PATH:/c/Go/bin (use this)
proto:
	cd src/proto; PATH=$$PATH:C:\\code\\go\\bin ../../../protoc -I . *.proto --go_out=.

proto-copy:
	cp src/proto/*.proto ../notekeeper-electron-frontend/app/proto/

proto-js:
	cd ../notekeeper-electron-frontend/app/proto; ../../../protoc -I . rpc.proto --js_out=import_style=commonjs,binary:./

proto-all: proto proto-copy proto-js

.PHONY: coverage