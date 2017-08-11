#!/bin/bash

# build binary distributions for linux/amd64 and darwin/amd64
set -e 

APP=private_s3_httpd
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
echo "working dir $DIR"

rm -rf $DIR/vendor
echo "... refreshing vendor directory"
./vendor.sh

echo "... running tests"
./test.sh || exit 1

arch=$(go env GOARCH)
version=$(cat $DIR/src/cmd/$APP/version.go | grep "const VERSION" | awk '{print $NF}' | sed 's/"//g')
goversion=$(go version | awk '{print $3}')

mkdir -p dist
for os in linux darwin; do
    echo "... building v$version for $os/$arch"
    BUILD=$(mktemp -d -t sortdb)
    TARGET="$APP-$version.$os-$arch.$goversion"
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 gb build
    mkdir -p $BUILD/$TARGET
    cp bin/$APP-$os-$arch $BUILD/$TARGET/$APP
    pushd $BUILD >/dev/null
    tar czvf $TARGET.tar.gz $TARGET
    if [ -e $DIR/dist/$TARGET.tar.gz ]; then
        echo "... WARNING overwriting dist/$TARGET.tar.gz"
    fi
    mv $TARGET.tar.gz $DIR/dist
    echo "... built dist/$TARGET.tar.gz"
    popd >/dev/null
done
