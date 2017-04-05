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
