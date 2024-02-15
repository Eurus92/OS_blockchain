#!/bin/bash
export SSHPASS='osgroup20!@#$'

sshpass -e ssh osgroup20@122.200.68.26 -p 8016 "killall -9 blockchain"
sshpass -e ssh osgroup20@122.200.68.26 -p 8017 "killall -9 blockchain"
sshpass -e ssh osgroup20@122.200.68.26 -p 8018 "killall -9 blockchain"
sshpass -e ssh osgroup20@122.200.68.26 -p 8019 "killall -9 blockchain"
sshpass -e ssh osgroup20@122.200.68.26 -p 8020 "killall -9 blockchain"

exit