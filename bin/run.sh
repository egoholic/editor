#!/bin/bash
DBNAME=stoa_blogging_$1 PORT=3000 DOMAIN=wearestoa.com go run ./targets/web/main.go
