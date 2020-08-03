/*
 * @Author  : sunlidong
 * @Date    : 2020年6月29日17:44:04
 * @Describe: 资产合约
 */
package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// 检验参数个数是否符合要求  (len(args)== count )
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

// 根据类型获取结构体指针
func getStructByType(assetType int64) (assetStruct interface{}, err error) {

	switch assetType {

	case assetCodeByQrCode:
		assetStruct = &BaseAssets{}

	case assetCodeByLogistics:
		assetStruct = &BaseAssets{}

	case assetCodeByManagement:
		assetStruct = &BaseAssets{}

	default:
		return nil, errors.New("The assetType is not supported:" + assetType)
	}

	return assetStruct, nil
}

//

// 获取结构体类型标识
func getArgsByfirst(args []string) (noStruct int64, err error) {

	// 不用判断 len ，前面已经判断过了

	if args[0] == "" {
		log.Printf("args[0] is %s:", args[0])
		return int64(0), errors.new("args[0] is nil ")

	}
	noStruct, err = strconv.ParseInt(arg[0], 10, 64)

	if err != nil {
		log.Printf("strconv.ParseInt args[0] is err %s:", err)
		return int64(0), err
	}

	return noStruct, nil

}

// 获取上链对象的key
func getkey(asset interface{}) (string, error) {
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
	v.FieldByName("updateDate").SetString(updateTime)
	return nil
}

// compositeKey的拼接规则
func createCompositeKey(stub shim.ChaincodeStubInterface, objectType string, attributes []string) string {
	var key, _ = stub.CreateCompositeKey(objectType, attributes)
	return key
}

// 普通key
func createCommonKey(stub shim.ChaincodeStubInterface, objectType string, attributes []string) string {
	key := objectType
	for _, value := range attributes {
		key = key + ":" + value
	}

	return key
}
