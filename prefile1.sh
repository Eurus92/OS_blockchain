#!/bin/bash
export SSHPASS='osgroup20!@#$'

cd src||exit
go build -o prefile transGen.go

./prefile 20

sshpass -e scp -P 8013 -r ./data osgroup20@122.200.68.26:/home/osgroup20/project3