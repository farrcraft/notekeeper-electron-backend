build:
	cd src; go build

# install the protoc golang plugin
# export GOPATH=C:/code/go
proto-gen:
	go get -u github.com/golang/protobuf/protoc-gen-go

# rebuild the proto definitions
# export PATH=$PATH:/cygdrive/c/code/go/bin 
proto:
	cd src/proto; PATH=$$PATH:/cygdrive/c/code/go/bin protoc -I . backend.proto --go_out=plugins=grpc:.
