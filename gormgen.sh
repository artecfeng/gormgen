#!/bin/bash

shellExit()
{
if [ $1 -eq 1 ]; then
    printf "\nerror!!!\n\n"
    exit 1
fi
}
printf "\n build gormgen\n\n"
go build main.go
printf "\n cp gormgen to gopath/bin\n\n"
cp main $GOPATH/bin/gormgen
printf "\n create and generating file\n\n"
gormgen -addr "root:@tcp(localhost:3306)/ginadmin"  -tables "admin"
shellExit $?



