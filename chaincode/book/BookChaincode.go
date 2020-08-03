/*
 * @Author  : sunlidong
 * @Date    : 2020年6月29日17:44:04
 * @Describe: 资产合约
 */
package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//
type BookChaincode struct {
}

var logger = shim.NewLogger("BookChaincode")

// 初始化操作
func (c *BookChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	Cow=newBase()
	return shim.Success(nil)
}

// 定义的各种操作
func (c *BookChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()
	logger.Debug(function, "--入参：　", args)

	var resp pb.Response

	switch function {
	// -------------------新增数据的时候 -----------------------

	/*

		------------- 常规操作
			1.  生成防伪标签
			2.	初始化防伪标签 Initializes the security label
			3.  防伪码激活 Activation of security code
			4.  流转中转 The circulation transshipment
			5.  验证 Data validation
			6.  查询二维码 Query qr code



		------------- 异常操作
			1. 冻结标签 Freeze tag
			2. 启动标签 Thaw the label
			3. 作废标签 Invalid label





	*/

	// 生成防伪标签
	case TagforNew:
		resp = c.generateSecurityLabelsByDb(stub, args)

	// 初始化标签
	case TagforIni:
		resp = c.initializesTheSecurityLabelByDb(stub, args)

	//	防伪标签激活
	case TagforAct:
		resp = c.activationOfSecurityCodeByDb(stub, args)
	
	//	中转
	case TagforTra:
	
		resp = c.theCirculationTransshipmentByDb(stub, args)

	case "getAssetsDataByID":
		resp = c.getAssetsDataByID(stub, args)
	//根据主键查询txID
	case "getAssetTxID":
		resp = c.getAssetTxID(stub, args)
	case "getAssetByBlurKey":
		// 根据 key 获取资产最新交易ID key
		// （用于区块链浏览器查询中，层级对象的key 由后台提供，这里只负责提供最新的key  txID）
		resp = c.getNewestAssetsKeyAndTxIDByBlurKey(stub, args)

		// 根据主键查询查询历史记录日志
	case "getHistoryByID":
		resp, _ = c.getHistoryByID(stub, args)

		//	根据主键查询数据的历史记录
	case "getDataHistoryByID":
		resp, _ = c.getDataHistoryByID(stub, args)
	case "getDataHistoryByTypeID":
		resp, _ = c.getDataHistoryByTypeID(stub, args)

		//-------------------------------------------2019-08-17 应收账款 数据上链
	case "uploadReceivable":
		resp = c.uploadReceivable(stub, args)

	default:
		resp = shim.Error(fmt.Sprintf("方法未定义:- %s", function))
	}
	logger.Debug(function, "--响应：　", " \n status:", resp.Status, " \n  Message:", resp.Message, " \n  Payload:", string(resp.Payload))
	return resp
}

func newBase() {
	return &Base{
		Name: "shixian"
	}
}

func main() {
	err := shim.Start(new(BookChaincode))
	if err != nil {
		fmt.Printf("Error starting ProductChaincode - %s", err)
	}
}
