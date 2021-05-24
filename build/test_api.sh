#!/bin/sh
#run autotest
MY_FN=`readlink -e $0`
ROOT_DIR=`dirname $MY_FN`/..

export PGSQL_HOST=127.0.0.1
export PGSQL_NAME=coins
export PGSQL_USER=coins
export PGSQL_PASS=coins
export PGSQL_PORT=5432
export CacheExpTime=10

cd $ROOT_DIR/build

docker-compose up -d

go clean -testcache
go test $ROOT_DIR/cmd/...

docker-compose  down