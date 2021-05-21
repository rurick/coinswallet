#!/bin/sh
#
# run docker container with the postgresql

sudo docker run --rm --name coins-pgdocker -e POSTGRES_PASSWORD=coins -e POSTGRES_USER=coins -e POSTGRES_DB=coins -d -p 5432:5432 -v $HOME/docker/volumes/postgres:/var/lib/postgresql/data postgres