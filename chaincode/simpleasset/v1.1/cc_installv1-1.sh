#!/bin/bash
set -x

# 1. v1.1 install

docker exec cli peer chaincode install -n simpleasset -v 1.1 -p github.com/simpleasset/v1.1
docker exec cli peer chaincode list --installed

# 2. v1.1 upgrade

docker exec cli peer chaincode instantiate -n simpleasset -v 1.1 -C mychannel -c '{"Args":[]}' -P 'AND ("Org1MSP.member")'
sleep 3
docker exec cli peer chaincode list --instantiated -C mychannel
