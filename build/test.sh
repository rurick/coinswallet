#!/bin/sh
#run autotest
MY_FN=`readlink -e $0`
ROOT_DIR=`dirname $MY_FN`/..

export PGSQL_HOST=127.0.0.1
export PGSQL_NAME=coins
export PGSQL_USER=coins
export PGSQL_PASS=coins
export PGSQL_PORT=5433
export CacheExpTime=10

docker run --rm --name coins-pgdocker -e POSTGRES_PASSWORD=coins -e POSTGRES_USER=coins -e POSTGRES_DB=coins -d -p 5433:5432 -v $HOME/docker/volumes/postgres:/var/lib/postgresql/data postgres

go clean -testcache
go test $ROOT_DIR/pkg/...
go test $ROOT_DIR/internal/...

docker stop coins-pgdocker
