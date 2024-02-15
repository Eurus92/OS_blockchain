#!/bin/bash
export SSHPASS='osgroup20!@#$'

sshpass -e ssh osgroup20@122.200.68.26 -p 8011 "rm -rf ./p3binary
          mkdir p3binary
          "

cd src||exit
GOOS=linux GOARCH=amd64 go build -o blockchain main.go broadcastBlock.go chain.go clientTransaction.go listen.go pow.go transactionPool.go
sshpass -e scp -P 8013 ./blockchain osgroup20@122.200.68.26:/home/osgroup20/project3

echo "finish binary"

exit
