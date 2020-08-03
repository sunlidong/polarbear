package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Roots interface {
	generateSecurityLabels(s shim.ChaincodeStubInterface, args []string) (interface{}, error) //生成防伪标签
	generateSecurityLabelsUpload(stub shim.ChaincodeStubInterface, args []string, assetStruct interface{}) (string, error)
	initializesTheSecurityLabel(stub shim.ChaincodeStubInterface, args []string) (interface{}, error) //初始化防伪标签
	initializesTheSecurityLabelUpload(stub shim.ChaincodeStubInterface, args []string, assetStruct interface{}) (string, error)
	activationOfSecurityCode(s shim.ChaincodeStubInterface, args []string) (interface{}, error)                              //防伪标签激活
	activationOfSecurityCodeUpload(stub shim.ChaincodeStubInterface, args []string, assetStruct interface{}) (string, error) //防伪标签激活
	theCirculationTransshipment(s shim.ChaincodeStubInterface, args []string) pb.Response                                    //中转
	dataValidation(s shim.ChaincodeStubInterface, args []string) pb.Response                                                 //验证
	queryQrCode(s shim.ChaincodeStubInterface, args []string) pb.Response                                                    //查询二维码
	freezeTag(s shim.ChaincodeStubInterface, args []string) pb.Response                                                      //冻结标签
	thawTheLabel(s shim.ChaincodeStubInterface, args []string) pb.Response                                                   //启动标签
	invalidLabel(s shim.ChaincodeStubInterface, args []string) pb.Response                                                   //作废标签
}
