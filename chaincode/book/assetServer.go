package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//  生成防伪标签 处理 ByDb
func (c *BookChaincode) generateSecurityLabelsByDb(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(c.generateSecurityLabels(stub,args))
}


//  生成防伪标签 处理 generateSecurityLabels
func (c *BookChaincode) generateSecurityLabels(stub shim.ChaincodeStubInterface, args []string) []byte {
	// 获取结构体 Key 
	StructType, err := Cow.generateSecurityLabels(stub, args)
	if err != nil {
		return Cow.response(401,err,nil)
	}

	stringKey err:= Cow.generateSecurityLabelsUpload(stub,args,StructType)
	if err !=nil{
		return Cow.response(401,err,nil)
	}

	return Cow.response(0,nil,stringKey)
}


//  初始化防伪标签 处理 ByDb
func (c *BookChaincode) initializesTheSecurityLabelByDb(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(c.initializesTheSecurityLabel(stub,args))
}


//  生成防伪标签 处理 generateSecurityLabels
func (c *BookChaincode) initializesTheSecurityLabel(stub shim.ChaincodeStubInterface, args []string) []byte {
	// 获取结构体 Key 
	StructType, err := Cow.initializesTheSecurityLabel(stub, args)
	if err != nil {
		return Cow.response(401,err,nil)
	}

	stringKey err:= Cow.initializesTheSecurityLabelUpload(stub,args,StructType)
	if err !=nil{
		return Cow.response(401,err,nil)
	}

	return Cow.response(0,nil,stringKey)
}




//  防伪标签激活 处理 ByDb
func (c *BookChaincode) activationOfSecurityCodeByDb(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(c.activationOfSecurityCode(stub,args))
}
//  防伪标签激活 处理 activationOfSecurityCode
func (c *BookChaincode) activationOfSecurityCode(stub shim.ChaincodeStubInterface, args []string) []byte {
	// 获取结构体 Key 
	StructType, err := Cow.activationOfSecurityCode(stub, args)
	if err != nil {
		return Cow.response(401,err,nil)
	}

	stringKey err:= Cow.activationOfSecurityCodeUpload(stub,args,StructType)
	if err !=nil{
		return Cow.response(401,err,nil)
	}

	return Cow.response(0,nil,stringKey)
}







//  中转 处理 activationOfSecurityCodeByDb
func (c *BookChaincode) theCirculationTransshipmentByDb(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(c.theCirculationTransshipment(stub,args))
}
//  中转 处理 activationOfSecurityCode
func (c *BookChaincode) theCirculationTransshipment(stub shim.ChaincodeStubInterface, args []string) []byte {
	// 获取结构体 Key 
	StructType, err := Cow.theCirculationTransshipment(stub, args)

	if err != nil {
		return Cow.response(401,err,nil)
	}

	stringKey err:= Cow.theCirculationTransshipmentUpload(stub,args,StructType)
	if err !=nil{
		return Cow.response(401,err,nil)
	}

	return Cow.response(0,nil,stringKey)
}