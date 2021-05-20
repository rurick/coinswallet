#!/bin/sh
MY_FN=`readlink -e $0`
MY_DIR=`dirname $MY_FN`

go test $MY_DIR/../...
