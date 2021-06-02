#!/bin/bash
set -ex
if [ "$VERSION" = "" ]; then
    VERSION=develop
fi
TD_VERSION=
URL_PREFIX=https://github.com/taosdata/TDengine/archive/refs
if [[ $VERSION =~ ^[[:alpha:]] ]]; then
  URL=$URL_PREFIX/heads/$VERSION.tar.gz
else
  URL=$URL_PREFIX/tags/ver-$VERSION.tar.gz
  TD_VERSION=$VERSION
  VERSION=ver-$VERSION
fi

wget -O src.tar.gz -c $URL
tar xf src.tar.gz
cd TDengine-${VERSION}/
if [ "$TD_VERSION" = "" ]; then
  TD_VERSION=`grep 'TD_VER_NUMBER "' \
      cmake/version.inc |sed -E 's/[^.0-9]//g'`
fi
echo "Build TDengine version" $TD_VERSION
sed 's#set -e##' -i packaging/release.sh
./packaging/release.sh -n $TD_VERSION -l lite -V stable