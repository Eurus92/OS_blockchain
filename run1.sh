#!/bin/bash
export SSHPASS='osgroup20!@#$'

sshpass -e ssh osgroup20@122.200.68.26 -p 8016 "cd p4_tree/src; ./blockchain 16 ./result/node11t.txt" 1&
sshpass -e ssh osgroup20@122.200.68.26 -p 8017 "cd p4_tree/src; ./blockchain 17 ./result/node12t.txt" 1&
sshpass -e ssh osgroup20@122.200.68.26 -p 8018 "cd p4_tree/src; ./blockchain 18 ./result/node13t.txt" 1&
sshpass -e ssh osgroup20@122.200.68.26 -p 8019 "cd p4_tree/src; ./blockchain 19 ./result/node14t.txt" 1&
sshpass -e ssh osgroup20@122.200.68.26 -p 8020 "cd p4_tree/src; ./blockchain 20 ./result/node15t.txt" 1&


#wait

#sshpass -e scp -P 8012 -r osgroup20@122.200.68.26:/home/osgroup20/project3/result .