#!/bin/bash
CPPATH=./builds/blog.exe
go build targets/web/main.go
cp main $CPPATH
rm main
echo "Compiled and stored in: $CPPATH"
