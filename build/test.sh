#!/bin/sh
MY_FN=`readlink -e $0`
ROOT_DIR=`dirname $MY_FN`/..

go clean -testcache

go test $ROOT_DIR/pkg/...
go test $ROOT_DIR/internal/...
