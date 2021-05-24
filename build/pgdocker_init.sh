#!/bin/sh
PGPASSWORD=coins psql -h localhost -p 5432 -U coins -d coins -f "init.sql"