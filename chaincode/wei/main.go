package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	TxId string = "txid-"
)

type FinanceChaincode struct {
}

type Finance struct {
	Id        string `json:"id"`        //ID
	Content   string `json:"content"`   //主要内容
	Accessory string `json:"accessory"` //附件
}

type FinanceList struct {
	TxId  string `json:"txid"`  //txid
	Value []byte `json:"value"` //value
}

type ResInfo struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
}

func (t *ResInfo) error(msg string) {
	t.Status = false
	t.Msg = msg
}
func (t *ResInfo) ok(msg string) {
	t.Status = true
	t.Msg = msg
}

func (t *ResInfo) response() pb.Response {
	resJson, err := json.Marshal(t)
	if err != nil {
		return shim.Error("Failed to generate json result " + err.Error())
	}
	return shim.Success(resJson)
}

type process func(shim.ChaincodeStubInterface, []string) *ResInfo

func (t *FinanceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *FinanceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	if function == "invoke" {
		return t.invoke(stub, args)
	} else if function == "delete" {
		return t.delete(stub, args)
	} else if function == "update" {
		return t.update(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	} else if function == "updateTxid" {
		return t.updateTxid(stub, args)
	} else if function == "history" {
		return t.history(stub, args)
	} else if function == "test" {
		return t.test(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"update\" \"query\"")
}

func (t *FinanceChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return t.handleProcess(stub, args, 3, func(shim.ChaincodeStubInterface, []string) *ResInfo {
		ri := &ResInfo{true, ""}
		_id := args[0]
		_content := args[1]
		_accessory := args[2]

		_finance := &Finance{_id, _content, _accessory}

		_ejson, err := json.Marshal(_finance)

		if err != nil {
			ri.error(err.Error())
		} else {
			_old, err := stub.GetState(_id)
			if err != nil {
				ri.error(err.Error())
			} else if _old != nil {
				ri.error("the finance has exists")
			} else {
				err := stub.PutState(_id, _ejson)
				if err != nil {
					ri.error(err.Error())
				} else {
					// 追加 根据 key 查询
					_tjson, err := json.Marshal(stub.GetTxID())
					if err != nil {
						ri.error(err.Error())
					} else {
						err := stub.PutState(TxId+_id, _tjson)
						if err != nil {
							ri.error(err.Error())
						} else {
							ri.ok(stub.GetTxID())
						}
					}
				}
			}
		}
		return ri
	})
}

func (t *FinanceChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return t.handleProcess(stub, args, 1, func(shim.ChaincodeStubInterface, []string) *ResInfo {
		ri := &ResInfo{true, ""}
		_id := args[0]
		_finance, err := stub.GetState(_id)
		if err != nil {
			ri.error(err.Error())
		} else {
			if _finance == nil {
				ri.ok("Warnning finance does not exists")
			} else {
				err := stub.DelState(_id)
				if err != nil {
					ri.error(err.Error())
				} else {
					ri.ok(stub.GetTxID())
				}
			}
		}
		return ri
	})
}

func (t *FinanceChaincode) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return t.handleProcess(stub, args, 3, func(shim.ChaincodeStubInterface, []string) *ResInfo {
		ri := &ResInfo{true, ""}
		_id := args[0]
		_content := args[1]
		_accessory := args[2]

		newfinance := &Finance{_id, _content, _accessory}

		_finance, err := stub.GetState(_id)

		if err != nil {
			ri.error(err.Error())
		} else {
			if _finance == nil {
				ri.error("Error the finance does not exists")
			} else {
				_ejson, err := json.Marshal(newfinance)
				if err != nil {
					ri.error(err.Error())
				} else {
					err := stub.PutState(_id, _ejson)
					if err != nil {
						ri.error(err.Error())
					} else {
						// 追加 根据 key 查询
						_tjson, err := json.Marshal(stub.GetTxID())
						if err != nil {
							ri.error(err.Error())
						} else {
							err := stub.PutState(TxId+_id, _tjson)
							if err != nil {
								ri.error(err.Error())
							} else {
								ri.ok(stub.GetTxID())
							}
						}
					}
				}
			}
		}
		return ri
	})
}

func (t *FinanceChaincode) updateTxid(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return t.handleProcess(stub, args, 2, func(shim.ChaincodeStubInterface, []string) *ResInfo {
		ri := &ResInfo{true, ""}
		_id := args[0]
		_txid := args[1]
		err := stub.PutState(TxId+_id, []byte(_txid))
		if err != nil {
			ri.error(err.Error())
		} else {
			ri.ok(stub.GetTxID())
		}
		return ri
	})
}

func (t *FinanceChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return t.handleProcess(stub, args, 1, func(shim.ChaincodeStubInterface, []string) *ResInfo {
		ri := &ResInfo{true, ""}
		queryString := args[0]

		value, err := stub.GetState(queryString)
		//rich查询，leveldb不支持
		//queryResults, err := getQueryResultForQueryString(stub, queryString)
		if err != nil {
			ri.error(err.Error())
		} else {
			ri.ok(string(value))
		}
		return ri
	})
}

func (t *FinanceChaincode) history(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return t.handleProcess(stub, args, 1, func(shim.ChaincodeStubInterface, []string) *ResInfo {
		ri := &ResInfo{true, ""}
		queryString := args[0]

		queryResults, err := getHistoryForKeyStrcts(stub, queryString)
		if err != nil {
			ri.error(err.Error())
		} else {
			ri.ok(string(queryResults))
		}
		return ri
	})
}

func (t *FinanceChaincode) handleProcess(stub shim.ChaincodeStubInterface, args []string, expectNum int, f process) pb.Response {
	res := &ResInfo{false, ""}
	err := t.checkArgs(args, expectNum)
	if err != nil {
		res.error(err.Error())
	} else {
		res = f(stub, args)
	}
	return res.response()
}

func (t *FinanceChaincode) checkArgs(args []string, expectNum int) error {
	if len(args) != expectNum {
		return fmt.Errorf("Incorrect number of arguments. Expecting  " + strconv.Itoa(expectNum))
	}
	for p := 0; p < len(args); p++ {
		if len(args[p]) <= 0 {
			return fmt.Errorf(strconv.Itoa(p+1) + "nd argument must be a non-empty string")
		}
	}
	return nil
}

func (t *FinanceChaincode) test(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return t.handleProcess(stub, args, 1, func(shim.ChaincodeStubInterface, []string) *ResInfo {
		ri := &ResInfo{true, ""}

		ri.ok("init server is ok ")

		return ri
	})
}

func getHistoryForKeyString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetHistoryForKey(queryString)
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"txid\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.TxId)

		buffer.WriteString("\"")

		buffer.WriteString(", \"time\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}
func getHistoryForKeyStrcts(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)
	var list []FinanceList
	var row FinanceList
	resultsIterator, err := stub.GetHistoryForKey(queryString)
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		row.TxId = queryResponse.TxId
		row.Value = queryResponse.Value
		list = append(list, row)
	}
	BytesList, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}
	return BytesList, nil
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func main() {
	err := shim.Start(new(FinanceChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
