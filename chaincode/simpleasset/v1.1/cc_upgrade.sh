#!/bin/bash
set -x

# 1. v1.1 install
docker exec cli peer chaincode install -n simpleasset -v 1.1.1 -p github.com/simpleasset/v1.1
# 2. v1.1 upgrade
docker exec cli peer chaincode upgrade -n simpleasset -v 1.1.1 -C mychannel -c '{"Args":[]}' -P 'AND ("Org1MSP.member")'
sleep 3

# 3. v1.1 invoke*3 a=100, b=200, transfer a->b 20
docker exec cli peer chaincode invoke -n simpleasset -C mychannel -c '{"Args":["set","a","100"]}'
docker exec cli peer chaincode invoke -n simpleasset -C mychannel -c '{"Args":["set","b","100"]}'
sleep 3
docker exec cli peer chaincode invoke -n simpleasset -C mychannel -c '{"Args":["transfer","a", "b","20"]}'
sleep 3

# 4. v1.1 query a, b, history b  
docker exec cli peer chaincode query -n simpleasset -C mychannel -c '{"Args":["get","a"]}'
docker exec cli peer chaincode query -n simpleasset -C mychannel -c '{"Args":["get","b"]}'
docker exec cli peer chaincode query -n simpleasset -C mychannel -c '{"Args":["history","b"]}'

# 5. v1.1 invoke del b 
docker exec cli peer chaincode invoke -n simpleasset -C mychannel -c '{"Args":["del","b"]}'
sleep 3
docker exec cli peer chaincode query -n simpleasset -C mychannel -c '{"Args":["get","b"]}'
docker exec cli peer chaincode query -n simpleasset -C mychannel -c '{"Args":["history","b"]}'
