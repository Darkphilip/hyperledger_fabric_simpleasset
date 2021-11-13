# !/bin/bash

# 1. 설치
docker exec cli peer chaincode install -n simpleasset -v 1.0 -p github.com/simpleasset/v1.0
docker exec cli peer chaincode list --installed
# 2. 배포 a = 100
docker exec cli peer chaincode instantiate -n simpleasset -v 1.0 -C mychannel -c '{"Args":["a","100"]}'  -P 'AND ("Org1MSP.member")'
sleep 3
docker exec cli peer chaincode list --instantiated -C mychannel
# 3. query get, a
docker exec cli peer chaincode query -n simpleasset -C mychannel -c '{"Args":["get", "a"]}'
# 4. invoke set, b = 100
docker exec cli peer chaincode invoke -n simpleasset -C mychannel -c '{"Args":["set", "b", "100"]}'
sleep 3
# 5. query get, b
docker exec cli peer chaincode query -n simpleasset -C mychannel -c '{"Args":["get", "b"]}'
