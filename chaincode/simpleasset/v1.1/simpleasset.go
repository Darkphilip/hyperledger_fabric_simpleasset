//  1. package 정의
package main

//  2. 외부 모듈 포함
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// 3. simpleasset 체인코드 클래스 구조체 정의
type SimpleAsset struct {
}

type Asset struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//  (TODO) Asset JSON 직렬화-marshal, 객체화-unmarshal -> WS 저장, 조회

//  4. Init 함수 구현 -> 체인코드 배포/ 업그레이드 할 때 한번 수행되는 코드
//  peer chaincode instantiate / upgrade -n simpleasset
// -v 1.0  -C mychannel -P 'AND("Org1MSP.member")' -c '{"Args:"["a", "100"]}
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// 내용삭제
	return shim.Success(nil)
}

// 5. Invoke 함수 구현
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	if fn == "set" {
		return t.Set(stub, args)
	} else if fn == "get" {
		return t.Get(stub, args)
	} else if fn == "del" {
		return t.Del(stub, args)
	} else if fn == "transfer" {
		return t.Transfer(stub, args)
	} else if fn == "history" {
		return t.History(stub, args)
	} else {
		return shim.Error("Not supported function name")
	}
}

// 6. Set 함수 구현
func (t *SimpleAsset) Set(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}

	//  JSON 구조체 생성
	var asset = Asset{Key: args[0], Value: args[1]}
	//  구조체에서 JSON 마샬
	assetAsBytes, _ := json.Marshal(asset)
	//  마샬된 결과를 Putstate
	err := stub.PutState(args[0], assetAsBytes)
	if err != nil {
		return shim.Error("Failed to set asset: " + args[0])
	}
	return shim.Success(assetAsBytes)
}

// 7.  get 함수 구현
func (t *SimpleAsset) Get(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting a key")
	}
	value, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get asset: " + args[0] + " with error:" + err.Error())
	}
	if value == nil {
		return shim.Error("Asset not found: " + args[0])
	}
	return shim.Success(value)
}

func (t *SimpleAsset) Del(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting a key")
	}
	value, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get asset: " + args[0] + " with error:" + err.Error())
	}
	if value == nil {
		return shim.Error("Asset not found: " + args[0])
	}
	err = stub.DelState(args[0])

	return shim.Success([]byte(args[0]))

}

func (t *SimpleAsset) Transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 { // from_key, to_key, amount
		return shim.Error("Incorrect arguments. Expecting a from_key, to_key, amount")
	}

	from_asset, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get asset: " + args[0] + " with error:" + err.Error())
	}
	if from_asset == nil {
		return shim.Error("Asset not found: " + args[0])
	}
	to_asset, err := stub.GetState(args[1])
	if err != nil {
		return shim.Error("Failed to get asset: " + args[1] + " with error:" + err.Error())
	}
	if to_asset == nil {
		return shim.Error("Asset not found: " + args[1])
	}
	// from과 to asset Json -> 구조체로 unmarshal
	from := Asset{}
	to := Asset{}
	json.Unmarshal(from_asset, &from)
	json.Unmarshal(to_asset, &to)

	// value string -> 계산가능한 int로 변환
	from_amount, _ := strconv.Atoi(from.Value)
	to_amount, _ := strconv.Atoi(to.Value)
	amount, _ := strconv.Atoi(args[2])

	// 검증
	if (from_amount - amount) < amount {
		shim.Error("Not enough asset value:" + args[0])
	}
	// 계산결과를 다시 구조체에 string형식으로 할당하는 부분
	from.Value = strconv.Itoa(from_amount - amount)
	to.Value = strconv.Itoa(to_amount + amount)

	// 구조체를 JSON으로 marshal
	from_asset, _ = json.Marshal(from)
	to_asset, _ = json.Marshal(to)

	// Putstate
	stub.PutState(args[0], from_asset)
	stub.PutState(args[1], to_asset)

	return shim.Success([]byte(args[0] + "-" + args[1] + ":" + args[2]))
}

func (t *SimpleAsset) History(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	assetName := args[0]

	fmt.Printf("- start getHistoryForAsset: %s\n", assetName)

	resultsIterator, err := stub.GetHistoryForKey(assetName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForAsset returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//  8. main 함수 구현
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
