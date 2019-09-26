# default target
build:
	go build -mod=vendor -tags debug; cp notekeeper-electron-backend.exe ../notekeeper-electron-frontend/app/resources/backend.exe

test: clean-test-db
	go test ./... -v -race -mod=vendor

clean-test-db:
	rm -f account/*.db user/*.db

coverage-dep:
	go get -u github.com/wadey/gocovmerge

# https://github.com/golang/go/issues/6909
coverage:
	shell ./coverage.sh

# create a release build
# The -s ldflag strips symbol table & debug info
release:
	go build -ldflags "-s" -mod=vendor

loc:
	find . -name '*.go' -not -path './proto/*' -not -path './vendor/*' | xargs wc -l

# install the protoc golang plugin
# export GOPATH=C:/code/go
proto-gen:
	go get -u github.com/golang/protobuf/protoc-gen-go

# rebuild the proto definitions
# export PATH=$PATH:/cygdrive/c/code/go/bin 
# export PATH=$PATH:/c/Go/bin (use this)
proto:
	cd proto; PATH=$$PATH:C:\\code\\go\\bin ../../protoc -I . *.proto --go_out=.

proto-copy:
	cp proto/*.proto ../notekeeper-electron-frontend/app/proto/

proto-js:
	cd ../notekeeper-electron-frontend/app/proto; ../../../protoc -I . *.proto --js_out=import_style=commonjs,binary:.

proto-all: proto proto-copy proto-js

.PHONY: coverage proto