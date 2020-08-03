/*
 * @Author  : sunlidong
 * @Date    : 2020年6月29日17:44:04
 * @Describe: 资产合约
 */
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hyperledger/fabric/common/util"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//
type BookChaincode struct {
}

var logger = shim.NewLogger("BookChaincode")

// 增量合约　想关信息
const (
	incrementContractName string = "IncrementChaincode"
	// 增加
	incrementContractMethod string = "addIncrementInfo"
	// 查询
	incrementContractQueryMethod string = "getIncrementInfo"
)

// 初始化操作
func (c *BookChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// 定义的各种操作
func (c *BookChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	logger.Debug(function, "--入参：　", args)
	var resp pb.Response
	switch function {
	// -------------------新增数据的时候 -----------------------
	case "uploadAsset":
		resp = c.uploadAsset(stub, args)
	case "updateAsset":
		resp = c.updateAsset(stub, args)
	case "getAssetByID":
		// 根据assetID,updateTime 获取当时的对象
		resp = c.getAssetByID(stub, args)
		// -------------------历史追溯的时候 查找变动的信息-----------------------
	case "getAssetByUpdateTime":
		// 根据assetID,updateTime 获取当时的对象
		resp, _ = c.getAssetByUpdateTime(stub, args)
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

const keySplitKey string = "+"
const resultSplitKey string = "=="

//  根据条件获取当时的资产 　有就返回　没有就选最优的数据（响应的时候 同时返回 当时的txID）
// 0　　AssetID 　要获取哪个资产的子集
func (c *BookChaincode) getAssetByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if err := checkArgsForCount(args, 1); err != nil {
		return shim.Error(err.Error())
	}
	currentAssetID := args[0]
	assetAsBytes, err := stub.GetState(currentAssetID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if assetAsBytes == nil {
		return shim.Error("the asset is not existed!")
	}

	return shim.Success(assetAsBytes)
}

//  上传资产附件(普通建池　返回的是交易ID)
// 0　　资产类型
// 1　　资产结构体json
// 2　　更新操作时间
func (c *BookChaincode) uploadAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 参数检查
	if err := checkArgs(args); err != nil {
		return shim.Error(err.Error())
	}
	assetType := args[0]
	// 获取上链结构体
	assetStruct, err := getAssetStructByType(assetType)
	if err != nil {
		return shim.Error(err.Error())
	}
	// 资产的上链
	if _, err := uploadAssetInternal(stub, args, assetStruct); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(stub.GetTxID()))
	//return shim.Success([]byte(stub.GetDecorations()))
}

/*
/***********************************************************************************************
*函数名 ： updateAsset
*函数功能描述 ：更新上链数据
*函数参数 ：
*函数返回值 ：上链数据的区块信息
*作者 ：孙利栋
*函数创建日期 ：2019/07/31
*函数修改日期 ：
*修改人 ：
*修改原因 ：
*版本 ：V01
*历史版本 ：V01
***********************************************************************************************/
func (c *BookChaincode) updateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 参数检查
	if err := checkArgs(args); err != nil {
		return shim.Error(err.Error())
	}
	assetType := args[0]
	// 获取上链结构体对象
	assetStruct, err := getAssetStructByType(assetType)
	if err != nil {
		return shim.Error(err.Error())
	}
	// 资产的更新
	if _, err := updateAssetInternal(stub, args, assetStruct); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(stub.GetTxID()))
}

// 操作资产前的　参数检查
func checkArgs(args []string) error {
	return checkArgsForCount(args, 3)
}

//
// 操作资产前的　参数检查 || 应收账款
func checkArgsRec(args []string) error {
	return checkArgsForCount(args, 4)
}

func checkArgsForCount(args []string, count int) error {
	if len(args) != count {
		return fmt.Errorf("Incorrect number of arguments. Expecting :  %v", count)
	}
	// 验空
	for index := 0; index < count; index++ {
		if len(args[index]) <= 0 {
			return fmt.Errorf("index :%v  argument must be a non-empty string", index)
		}
	}
	return nil
}

// 根据资产type,获取资产结构对象
func getStructByType(assetType string) (interface{}, error) {
	var assetStruct interface{}
	switch assetType {

	//	01. 资产集
	case assetType_ASSET:
		assetStruct = &Assets{}

		//	02. 附件集
	case assetType_FUJIAN:
		assetStruct = &ProAttachmentlist{}

		//	03. 尽调结果集
	case assetType_SURVEY:
		assetStruct = &BaseSurvey{}

		//	05. 尽调报告集
	case assetType_REPORT:
		assetStruct = &BaseReport{}

		//	06. 附件集
	case assetType_ATTACHMENT:
		assetStruct = &ProAttachmentlist{}

		//	07. 历史记录
	case assetType_HISTORYLOG:
		assetStruct = &HistoryLog{}

		//	08. 应收账款
	case assetType_RECEIVABLES:
		assetStruct = &AstInfo{}

	//-------------------------------------------2019-08-17 紧急调整 || 应收账款 结构
	//	09. 关联列表
	case assetType_ASSETLIST:
		assetStruct = &BaseAssetsList{}

	default:
		return nil, errors.New("The assetType is not supported:" + assetType)
	}
	return assetStruct, nil
}

// 根据资产type,获取资产结构对象 []
func getAssetStructByTypelist(assetType string) (interface{}, error) {
	var assetStruct interface{}
	switch assetType {

	//	01. 资产集
	case assetType_ASSET:
		assetStruct = &[]Assets{}

		//	02. 附件集
	case assetType_FUJIAN:
		assetStruct = &[]ProAttachmentlist{}

		//	03. 尽调结果集
	case assetType_SURVEY:
		assetStruct = &[]BaseSurvey{}

		//	05. 尽调报告集
	case assetType_REPORT:
		assetStruct = &[]BaseReport{}

		//	06. 附件集
	case assetType_ATTACHMENT:
		assetStruct = &[]ProAttachmentlist{}

		//	07. 历史记录
	case assetType_HISTORYLOG:
		assetStruct = &[]HistoryLog{}

	default:
		return nil, errors.New("The assetType is not supported:" + assetType)
	}
	return assetStruct, nil
}

//// new  chaincode
////  根据条件获取资产ID的子集资产  getChangedAssets
//// 0　　assetID  基础资产ID
//// 1　　currentAssetID 　要获取哪个资产的子集
//// 2　　updateTime　更新操作时间
func (c *BookChaincode) getAssetsDataByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验参
	if err := checkArgsForCount(args, 3); err != nil {
		return shim.Error(err.Error())
	}
	// 资产ID
	currentAssetID := args[1]
	//// 临时存储数据  初始赋值
	var resultMap = make(map[string]*SubAssetsResult)
	resultMap[currentAssetID] = &SubAssetsResult{}
	assetID := args[0]
	updateTime := args[2]
	// 获取增量信息
	chaincodeArgs := util.ToChaincodeArgs(incrementContractQueryMethod, assetID, updateTime)
	response := stub.InvokeChaincode(incrementContractName, chaincodeArgs, "")
	if response.Status != shim.OK {
		return response
	}
	// 遍历增量信息　根据类型判断　需要找的对象
	incrementKeys := strings.Split(string(response.Payload), resultSplitKey)
	for _, incrementKey := range incrementKeys {
		if len(incrementKey) == 0 {
			// 过滤空值
			continue
		}
		fields := strings.Split(incrementKey, keySplitKey)
		//// 特殊情况排除 (非本底层资产， 长度不为5的都排除)
		//if len(fields) != 5 || !strings.EqualFold(fields[3], currentAssetID)  {
		//	continue
		//}
		myAssetType := fields[2]
		assetKey := fields[3]
		switch myAssetType {

		// 添加关联附件列表 8
		case "8":
			assetTmp, err := c.getProAttachByUpdateTime(stub, []string{assetKey, updateTime})
			if err != nil {
				return shim.Error(err.Error())
			}
			// 添加关联附件
			resultMap[currentAssetID].ProAttachmentlist = append(resultMap[currentAssetID].ProAttachmentlist, assetTmp)

			// 添加关联尽调结果列表 9
		case "9":
			assetTmp, err := c.getBaseSurveyByUpdateTime(stub, []string{assetKey, updateTime})
			if err != nil {
				return shim.Error(err.Error())
			}
			// 添加关联尽调报告
			resultMap[currentAssetID].BaseSurvey = append(resultMap[currentAssetID].BaseSurvey, assetTmp)

			// 添加关联尽调报告列表 10
		case "10":
			assetTmp, err := c.getBaseReportByUpdateTime(stub, []string{assetKey, updateTime})
			if err != nil {
				return shim.Error(err.Error())
			}
			// 添加关联尽调报告
			resultMap[currentAssetID].BaseReport = append(resultMap[currentAssetID].BaseReport, assetTmp)

		default:

		}
	}
	//格式化后　返回数据
	b, err := json.Marshal(resultMap)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(b)

}

//
func (c *BookChaincode) getAssetByUpdateTime(stub shim.ChaincodeStubInterface, args []string) (pb.Response, string) {
	if err := checkArgsForCount(args, 2); err != nil {
		return shim.Error(err.Error()), ""
	}
	currentAssetID := args[0]
	updateTime := args[1]

	assetMap := make(map[string][]byte)
	assetTxIDMap := make(map[string]string)

	//1 迭代查找历史数据
	resultsIterator, err := stub.GetHistoryForKey(currentAssetID)
	if err != nil {
		return shim.Error(err.Error()), ""
	}
	defer resultsIterator.Close()
	if !resultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("the assetID : < %s > has not any history info", currentAssetID)), ""
	}
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()

		if err != nil {
			return shim.Error(err.Error()), ""
		}
		//2 查找更新时间　并存入map中
		assetUpdateTime := getUpdateTime(string(response.Value))
		if len(assetUpdateTime) > 0 {
			assetMap[assetUpdateTime] = response.Value
			assetTxIDMap[assetUpdateTime] = response.TxId
		} else {
			return shim.Error("can not find value about field of  updateDate in key:" + currentAssetID), ""
		}
	}
	keys := []string{}
	//3 根据日期返回对象
	for key, value := range assetMap {
		keys = append(keys, key)
		if strings.EqualFold(key, updateTime) {
			return shim.Success(value), assetTxIDMap[key]
		}
	}
	//4 记录首次操作时间
	sort.Strings(keys)
	startTime := keys[0]

	//5 否则　查找最近日期对象
	keys = append(keys, updateTime)
	sort.Strings(keys)
	for index, time := range keys {
		if strings.EqualFold(time, updateTime) {
			//6 判断是否为头  如果是头的话　就有问题
			if index == 0 {
				return shim.Error(fmt.Sprintf("the  assetID :%s  startTime is  %s  ,but receive : %s", currentAssetID, startTime, updateTime)), ""
			}
			return shim.Success(assetMap[keys[index-1]]), assetTxIDMap[keys[index-1]]
		}
	}
	return shim.Error(fmt.Sprintf("failed to find asset history info  by  assetID : %s , updateTime : %s ", currentAssetID, updateTime)), ""
}

// 根据资产key 、更新时间  获取资产
func (c *BookChaincode) getProAttachByUpdateTime(stub shim.ChaincodeStubInterface, args []string) (ProAttachmentlist, error) {

	assetTmp := ProAttachmentlist{}
	if err := checkArgsForCount(args, 2); err != nil {
		return assetTmp, err
	}
	assetKey := args[0]
	updateTime := args[1]
	// 根据key获取到底层资产
	res, _ := c.getAssetByUpdateTime(stub, []string{assetKey, updateTime})
	if res.Status != shim.OK {
		return assetTmp, fmt.Errorf("%s", res.Message)
	}
	// 反序列化对象
	if err := json.Unmarshal(res.Payload, &assetTmp); err != nil {
		return assetTmp, err
	}
	fmt.Println("un:", assetTmp)
	return assetTmp, nil
}

//  根据条件获取当时的资产 　有就返回　没有就选最优的数据（响应的时候 同时返回 当时的txID）
// 0　　AssetID 　要获取哪个资产的子集
func (c *BookChaincode) getAssetTxID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if err := checkArgsForCount(args, 1); err != nil {
		return shim.Error(err.Error())
	}
	currentAssetID := args[0]
	assetAsBytes, err := stub.GetState(currentAssetID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if assetAsBytes == nil {
		return shim.Error("the asset is not existed!")
	}

	return shim.Success([]byte(stub.GetTxID()))
}

// 根据key产品
//0  assetKey(产品 或者 底层)
func (c *BookChaincode) getNewestAssetsKeyAndTxIDByBlurKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验参
	if err := checkArgsForCount(args, 1); err != nil {
		return shim.Error(err.Error())
	}
	assetKey := args[0]
	//获取历史记录
	resultsIterator, err := stub.GetHistoryForKey(assetKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	// 获取最新的txid
	var currentTxID string
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		currentTxID = response.TxId
		// logger.Debug(response)
	}
	if currentTxID == "" {
		return shim.Error("The asset is not existed about key " + assetKey)
	}
	// 组装结果
	result := &SubAssetKeyAndTxIDResult{TxID: currentTxID, Key: assetKey}
	resultAsBytes, err := json.Marshal(result)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(resultAsBytes)
}

func (c *BookChaincode) getHistoryByID(stub shim.ChaincodeStubInterface, args []string) (pb.Response, string) {
	if err := checkArgsForCount(args, 1); err != nil {
		return shim.Error(err.Error()), ""
	}
	currentAssetID := args[0]
	var HisList []HistoryLog

	resultsIterator, err := stub.GetHistoryForKey(currentAssetID)
	if err != nil {
		return shim.Error(err.Error()), ""
	}
	defer resultsIterator.Close()
	if !resultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("the assetID : < %s > has not any history info", currentAssetID)), ""
	}
	for resultsIterator.HasNext() {
		var row HistoryLog
		response, err := resultsIterator.Next()

		if err != nil {
			return shim.Error(err.Error()), ""
		}

		json.Unmarshal(response.Value, &row)

		if row.HisUUID != "" {
			HisList = append(HisList, row)
		}
	}

	jsonAsBytes, err := json.Marshal(HisList)

	if err != nil {
		return shim.Error(err.Error()), ""
	}
	return shim.Success(jsonAsBytes), ""
}

//  getDataHistoryByID
func (c *BookChaincode) getDataHistoryByID(stub shim.ChaincodeStubInterface, args []string) (pb.Response, string) {
	if err := checkArgsForCount(args, 1); err != nil {
		return shim.Error(err.Error()), ""
	}
	currentAssetID := args[0]
	var HisList []HistoryLog

	resultsIterator, err := stub.GetHistoryForKey(currentAssetID)
	if err != nil {
		return shim.Error(err.Error()), ""
	}
	defer resultsIterator.Close()
	if !resultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("the assetID : < %s > has not any history info", currentAssetID)), ""
	}
	for resultsIterator.HasNext() {
		var row HistoryLog
		response, err := resultsIterator.Next()

		if err != nil {
			return shim.Error(err.Error()), ""
		}

		json.Unmarshal(response.Value, &row)

		if row.HisUUID != "" {
			HisList = append(HisList, row)
		}
	}

	jsonAsBytes, err := json.Marshal(HisList)

	if err != nil {
		return shim.Error(err.Error()), ""
	}
	return shim.Success(jsonAsBytes), ""
}

//  查询对应ID的历史信息
// 0　　资产类型
// 1　　资产结构体json
//  getDataHistoryByID
func (c *BookChaincode) getDataHistoryByTypeID(stub shim.ChaincodeStubInterface, args []string) (pb.Response, string) {
	if err := checkArgsForCount(args, 1); err != nil {
		return shim.Error(err.Error()), ""
	}
	currentAssetID := args[0]
	//currentAssetID := args[1]
	// 获取上链结构体
	var arr []string

	resultsIterator, err := stub.GetHistoryForKey(currentAssetID)
	if err != nil {
		return shim.Error(err.Error()), ""
	}
	defer resultsIterator.Close()
	if !resultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("the assetID : < %s > has not any history info", currentAssetID)), ""
	}
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()

		if err != nil {
			return shim.Error(err.Error()), ""
		}
		arr = append(arr, string(response.Value))
	}
	jsonAsBytes, err := json.Marshal(arr)

	if err != nil {
		return shim.Error(err.Error()), ""
	}
	return shim.Success(jsonAsBytes), ""
}

//
// 根据尽调结果 key 、更新时间  获取尽调结果
func (c *BookChaincode) getBaseSurveyByUpdateTime(stub shim.ChaincodeStubInterface, args []string) (BaseSurvey, error) {

	assetTmp := BaseSurvey{}
	if err := checkArgsForCount(args, 2); err != nil {
		return assetTmp, err
	}
	assetKey := args[0]
	updateTime := args[1]
	// 根据key获取到底层资产
	res, _ := c.getAssetByUpdateTime(stub, []string{assetKey, updateTime})
	if res.Status != shim.OK {
		return assetTmp, fmt.Errorf("%s", res.Message)
	}
	// 反序列化对象
	if err := json.Unmarshal(res.Payload, &assetTmp); err != nil {
		return assetTmp, err
	}
	fmt.Println("un:", assetTmp)
	return assetTmp, nil
}
 

//
//  上传应收账款资产(普通建池　返回的是交易ID)
// 0　　资产类型
// 1　　资产结构体json
// 2　　更新操作时间
// 3   	历史记录
func (c *BookChaincode) uploadReceivable(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 参数检查
	if err := checkArgsRec(args); err != nil {
		return shim.Error(err.Error())
	}
	assetType := args[0]
	// 获取上链结构体
	assetStruct, err := getAssetStructByType(assetType)
	if err != nil {
		return shim.Error(err.Error())
	}
	// 资产的上链
	if _, err := uploadAssetInternalRec(stub, args, assetStruct); err != nil {
		return shim.Error(err.Error())
	}

	historyStruct, err := getAssetStructByType(assetType_HISTORYLOG)
	if err != nil {
		return shim.Error(err.Error())
	}
	txID := stub.GetTxID()
	//历史记录上链
	if _, err := uploadAssetInternalHistory(stub, args, historyStruct, txID); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(stub.GetTxID()))
}

func main() {
	err := shim.Start(new(BookChaincode))
	if err != nil {
		fmt.Printf("Error starting ProductChaincode - %s", err)
	}
}
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// 资产更新
func updateAssetInternal(stub shim.ChaincodeStubInterface, args []string, asset interface{}) (string, error) {
	var key string
	var err error
	if len(args) != 3 {
		return key, errors.New("Incorrect number of arguments. Expecting 3")
	}
	assetJsonStr := args[1]
	updateTime := args[2]
	logger.Debug("获取到asset:", asset)
	//１ 反序列化传递进来的的对象
	if err := json.Unmarshal([]byte(assetJsonStr), asset); err != nil {
		return key, errors.New("第一次反序列话的时候 出错," + err.Error())
	}
	//2 获取上链key
	key, err = getPkey(asset)
	if err != nil {
		return key, err
	}
	//3 检查是否上链
	assetAsBytes, err := stub.GetState(key)
	if err != nil {
		return key, err
	}
	if assetAsBytes == nil {
		return key, errors.New("The asset  is　not existed! the key :" + key)
	}
	//4 反序列链上对象
	if err := json.Unmarshal(assetAsBytes, asset); err != nil {
		return key, errors.New("4 反序列链上对象 出错," + err.Error())
	}
	//5　获取原有　父级信息
	parentID, parentType := getParentIDAndParentType(asset)
	// 6 为新对象添加 父级信息
	if err := json.Unmarshal([]byte(assetJsonStr), asset); err != nil {
		return key, errors.New("6 为新对象添加 父级信息 出错," + err.Error())
	}
	if err = addParentInfo(parentID, parentType, asset); err != nil {
		return key, err
	}
	// 7 添加push时间
	if err = addUpdateTime(updateTime, asset); err != nil {
		return key, err
	}
	//// 8 　检查上链资产中的子资产是否关联了本资产
	//if err = addRelationShipForSubAsset(stub, asset, key); err != nil {
	//	return key, err
	//}
	// 9　上链
	assetAsBytes, err = json.Marshal(asset)
	if err != nil {
		return key, err
	}
	if err = stub.PutState(key, assetAsBytes); err != nil {
		return key, err
	}
	return key, nil
}

// 获取现有资产的　parentID type
func getParentIDAndParentType(asset interface{}) (string, string) {
	var parentID, parentType string
	switch t := asset.(type) {
	case *Assets:
		parentID = t.ParentID
		parentType = t.ParentType
	case *ProAttachmentlist:
		parentID = t.ParentID
		parentType = t.ParentType
	case *BaseSurvey:
		parentID = t.ParentID
		parentType = t.ParentType
	case *BaseReport:
		parentID = t.ParentID
		parentType = t.ParentType
	}
	return parentID, parentType
}

// 资产上链
func uploadAssetInternal(stub shim.ChaincodeStubInterface, args []string, assetStruct interface{}) (string, error) {
	var key string
	var err error
	if len(args) != 3 {
		return key, errors.New("Incorrect number of arguments. Expecting 3")
	}
	assetJsonStr := args[1]
	updateTime := args[2]
	logger.Debug("开始步骤")
	//１ 反序列化
	if err := json.Unmarshal([]byte(assetJsonStr), assetStruct); err != nil {
		return key, errors.New("//１ 反序列化:" + err.Error())
	}
	//2 获取上链key
	key, err = getPkey(assetStruct)
	if err != nil {
		return key, errors.New("//2 获取上链:" + err.Error())
	}
	//3 检查是否上链
	existed, err := checkIsExisted(stub, key)
	if err != nil {
		return key, errors.New("//3 检查是否上链:" + err.Error())
	}
	if existed {
		return key, fmt.Errorf("The asset  is existed! the key : %s", key)
	}
	// 4 添加push时间
	if err = addUpdateTime(updateTime, assetStruct); err != nil {
		return key, errors.New("// 4 添加push时间:" + err.Error())
	}
	// 5　上链
	assetAsBytes, err := json.Marshal(assetStruct)
	if err != nil {
		return key, errors.New("// 5　上链:" + err.Error())
	}
	if err = stub.PutState(key, assetAsBytes); err != nil {
		return key, err
	}
	// 6 添加关联关系
	/*	if err = addRelationShipForSubAsset(stub, assetStruct, key); err != nil {
			return key, errors.New("// 6 添加关联关系:"+err.Error())

		}
	*/
	return key, nil
}

// 资产上链 ||应收账款
func uploadAssetInternalRec(stub shim.ChaincodeStubInterface, args []string, assetStruct interface{}) (string, error) {
	var key string
	var err error
	if len(args) != 4 {
		return key, errors.New("Incorrect number of arguments. Expecting 3")
	}
	assetJsonStr := args[1]
	updateTime := args[2]
	logger.Debug("开始步骤")
	//１ 反序列化
	if err := json.Unmarshal([]byte(assetJsonStr), assetStruct); err != nil {
		return key, errors.New("//１ 反序列化:" + err.Error())
	}
	//2 获取上链key
	key, err = getPkey(assetStruct)
	if err != nil {
		return key, errors.New("//2 获取上链:" + err.Error())
	}
	//3 检查是否上链
	existed, err := checkIsExisted(stub, key)
	if err != nil {
		return key, errors.New("//3 检查是否上链:" + err.Error())
	}
	if existed {
		return key, fmt.Errorf("The asset  is existed! the key : %s", key)
	}
	// 4 添加push时间
	if err = addUpdateTime(updateTime, assetStruct); err != nil {
		return key, errors.New("// 4 添加push时间:" + err.Error())
	}
	// 5　上链
	assetAsBytes, err := json.Marshal(assetStruct)
	if err != nil {
		return key, errors.New("// 5　上链:" + err.Error())
	}
	if err = stub.PutState(key, assetAsBytes); err != nil {
		return key, err
	}
	return key, nil
}

//　检查是否上链了
func checkIsExisted(stub shim.ChaincodeStubInterface, key string) (bool, error) {
	// 校验key
	if key == "" {
		return false, errors.New("the pkey‘s value  is empty")
	}
	// 判断链上是否存在
	assetAsBytes, err := stub.GetState(key)
	if err != nil {
		return false, err
	}
	if assetAsBytes != nil {
		return true, nil
	}
	return false, nil
}

// 为结构体添加更新时间
func addUpdateTime(updateTime string, asset interface{}) error {
	v := reflect.ValueOf(asset).Elem()
	v.FieldByName("UpdateDate").SetString(updateTime)
	return nil
}

// 为结构体添加txID
func addUpdateTxID(updateTxID string, asset interface{}) error {
	v := reflect.ValueOf(asset).Elem()
	v.FieldByName("HisCurrentTx").SetString(updateTxID)
	return nil
}

// 为结构体添加更新时间
func addParentInfo(parentID, parentType string, asset interface{}) error {
	v := reflect.ValueOf(asset).Elem()
	v.FieldByName("ParentID").SetString(parentID)
	//v.FieldByName("ParentID")
	v.FieldByName("ParentType").SetString(parentType)
	return nil
}

// 获取上链对象的key
func getPkey(asset interface{}) (string, error) {
	t := reflect.TypeOf(asset).Elem()
	v := reflect.ValueOf(asset).Elem()
	fieldCount := t.NumField()
	for index := 0; index < fieldCount; index++ {
		// 有这个标识的　字段　对应的值
		_, hasKey := t.Field(index).Tag.Lookup("pkey")
		if hasKey {
			return v.FieldByName(t.Field(index).Name).String(), nil
		}
	}
	return "", errors.New("The asset has no field  about primary key!")
}

// 添加资产关联关系
func addRelationShipForSubAsset(stub shim.ChaincodeStubInterface, asset interface{}, assetKey string) error {
	switch v := asset.(type) {
	//基础资产
	case *Assets:
		// id := v.PreRegNo   //产品信息不能强制关联  子资产
		id := v.ProUUID // 资产唯一标识的
		//  为底层资产添加　parentID
		var i interface{}

		//根据资产类型 　　确定需要关联　保理资产　还是　底层资产
		switch v.ProType {
		case "0":
			rrAssetTmp := &ProAttachmentlist{}
			rrAssetTmp.ParentType = assetType_FUJIAN
			i = rrAssetTmp
		default:
			return errors.New("Product assetType only supports 0 and 1 ,　but recived : " + v.ProType)
		}
		if err := changeRelationByArrary(stub, id, v.ProNote, i); err != nil {
			return err
		}
		//  汇总产品信息  key  挂靠到asset下面
		products, err := stub.GetState("UUID-" + v.ProUUID)
		if err != nil {
			return err
		}
		// 如果是 value 为空  直接添加
		if products == nil {
			if err := stub.PutState("UUID-"+v.ProUUID, []byte(assetKey)); err != nil {
				return err
			}
			// 如果不包含这个key　　组装后添加进去
		} else if !strings.Contains(string(products), assetKey) {
			products = []byte(string(products) + "-" + assetKey)
			if err := stub.PutState("UUID-"+v.ProUUID, products); err != nil {
				return err
			}
		}
	}
	return nil
}

//根据key改变其　parentID
// 逻辑：　根据key 获取对象　　判读是否有parentID 　没有的话就赋值
func changeRelationByArrary(stub shim.ChaincodeStubInterface, parentID string, keys []string, asset interface{}) error {
	// 通过反射改值
	v := reflect.ValueOf(asset).Elem()
	parentTypeValue := v.FieldByName("ParentType").String()
	//遍历发票、贸易合同、其他附件集合　　更改其parentID
	for _, key := range keys {
		// 取对象
		objAsBytes, err := stub.GetState(key)
		if err != nil {
			return err
		}
		if objAsBytes == nil {
			return errors.New("The asset  is　not existed! the key :" + key)
		}
		if err := json.Unmarshal(objAsBytes, asset); err != nil {
			return err
		}
		v := reflect.ValueOf(asset).Elem()
		fieldValue := v.FieldByName("ParentID").String()
		if strings.Contains(fieldValue, "<") { //判断值是否为string 类型
			//　非string 类型
			return errors.New(fieldValue)

		} else if len(fieldValue) > 0 { //parentID已经有值了
			continue
		}
		// 　设置父级 id 以及 type
		v.FieldByName("ParentID").SetString(parentID)
		v.FieldByName("ParentType").SetString(parentTypeValue)
		// 存对象
		if objAsBytes, err = json.Marshal(asset); err != nil {
			return err
		}
		if err := stub.PutState(key, objAsBytes); err != nil {
			return err
		}
	}
	return nil
}

// 根据jsonObj 获取　更新时间
func getUpdateTime(str string) string {
	// logger.Debug("getUpdateTime 入参   ： ", str)
	resultStr := ""
	ss := strings.Split(str, ",")
	for _, value := range ss {
		if strings.Contains(value, "updateDate") {
			fieldValue := strings.Split(value, ":")[1]
			resultStr = strings.Split(fieldValue, "\"")[1]
			break
		}
	}
	// logger.Debug("getUpdateTime 结果   ： ", resultStr)
	return resultStr
}

// 资产上链 历史记录上链
func uploadAssetInternalHistory(stub shim.ChaincodeStubInterface, args []string, historyStruct interface{}, TxID string) (string, error) {
	var key string
	var err error
	if len(args) != 4 {
		return key, errors.New("Incorrect number of arguments. Expecting 4")
	}
	assetJsonStrHistory := args[3]
	updateTime := args[2]
	logger.Debug("开始步骤")
	//１ 反序列化
	if err := json.Unmarshal([]byte(assetJsonStrHistory), historyStruct); err != nil {
		return key, errors.New("//１ 反序列化:" + err.Error())
	}
	//2 获取上链key
	key, err = getPkey(historyStruct)
	if err != nil {
		return key, errors.New("//2 获取上链:" + err.Error())
	}
	//3 检查是否上链
	existed, err := checkIsExisted(stub, key)
	if err != nil {
		return key, errors.New("//3 检查是否上链:" + err.Error())
	}
	if existed {
		return key, fmt.Errorf("The asset  is existed! the key : %s", key)
	}
	// 4 添加push时间
	if err = addUpdateTime(updateTime, historyStruct); err != nil {
		return key, errors.New("// 4 添加push时间:" + err.Error())
	}

	// 4 添加tx ID
	if err = addUpdateTxID(TxID, historyStruct); err != nil {
		return key, errors.New("// 4 添加push TxID:" + err.Error())
	}
	// 5　上链
	historyBytes, err := json.Marshal(historyStruct)
	if err != nil {
		return key, errors.New("// 5　上链:" + err.Error())
	}
	if err = stub.PutState(key, historyBytes); err != nil {
		return key, err
	}

	return key, nil
}
