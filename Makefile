# default target
build:
	cd src; go build -tags debug

test:
	cd src; go test ./... -v -race

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
proto:
	cd src/proto; PATH=$$PATH:C:\\code\\go\\bin ../../../protoc -I . *.proto --go_out=.

proto-copy:
	cp src/proto/*.proto ../notekeeper-electron-frontend/app/proto/

proto-js:
	cd ../notekeeper-electron-frontend/app/proto; ../../../protoc -I . rpc.proto --js_out=import_style=commonjs,binary:./

proto-all: proto proto-copy proto-js
