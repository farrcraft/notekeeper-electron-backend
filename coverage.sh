#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SRC_PATH=$(realpath $DIR/src)
VENDOR_PATH=$(realpath $DIR/vendor/)
BASE_PATH=$(realpath $DIR/)

mkdir -p $BASE_PATH/coverage

COVER_MODULES=$(cd src; go list -f 'src/{{.Name}}' ./... | grep -v main | grep -v notekeeper)
for dir in $COVER_MODULES; do
    pushd .
    cd $dir
    COVERAGE_FILE=$(echo "coverage_$dir" | sed -e 's/src\///g')
    go test -v -coverprofile=$BASE_PATH/coverage/$COVERAGE_FILE.out
    popd
done

PATH=$PATH:C:\\code\\go\\bin gocovmerge coverage/*.out | sed -e 's/_\\C_\\code\\NoteKeeper\.io\\notekeeper-electron-backend\\src\\//g' > coverage.out
GOPATH=C:\\code\\NoteKeeper.io\\notekeeper-electron-backend go tool cover -html=coverage.out -o coverage.html
rm coverage.out
rm coverage/*
mv coverage.html coverage/

exit 0
