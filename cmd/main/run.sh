#!/bin/bash

exec go get -v "github.com/go-sql-driver/mysql" &
wait
exec go get -v "github.com/gorilla/websocket" &
wait
exec go build main.go storage.go webSocket.go &
wait
exec ./main
