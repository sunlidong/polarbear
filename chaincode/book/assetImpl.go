package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)



// generateSecurityLabels 参数处理 --校验-- 创建二维码
func (p *Base) generateSecurityLabels(stub shim.ChaincodeStubInterface, args []string) (interface{}, error) {

	if err := checkArgsForCount(args, 3); err != nil {
		return shim.Error(err.Error())
	}

	// 获取结构体标识符 euId
	euId, err := getArgsByfirst(args)

	if err != nil {
		return nil, err
	}

	// 获取结构体标识
	assetStruct, err := getStructByType(euId)
	if err != nil {
		return nil, err
	}

	return assetStruct, nil

}
//  generateSecurityLabelsUpload --上链-- 创建二维码
func (p *Base) generateSecurityLabelsUpload(stub shim.ChaincodeStubInterface, args []string, assetStruct interface{}) (string, error) {

	var key string 
	var err error
	if len(args) != 3 {
		return key, errors.New("Incorrect number of arguments. Expecting 3")
	}
	assetJsonStr := args[1]
	updateTime := time.Now().Format("2006-01-02 15:04:05")

	logger.Debug("开始步骤")
	//１ 反序列化
	if err := json.Unmarshal([]byte(assetJsonStr), assetStruct); err != nil {
		return key, errors.New("//１ 反序列化:" + err.Error())
	}
	//2 获取上链key
	key, err = getkey(assetStruct)
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

	// 创建复合键
	comkey := createCompositeKey(stub, CompositeKeyByQr, key)

	// 5　序列化
	assetAsBytes, err := json.Marshal(assetStruct)
	if err != nil {
		return key, errors.New("// 5　上链:" + err.Error())
	}
	if err = stub.PutState(comkey, assetAsBytes); err != nil {
		return key, err
	}

	return key, nil
}




// initializesTheSecurityLabel 参数处理 --校验-- 初始二维码
func (p *Base) initializesTheSecurityLabel(stub shim.ChaincodeStubInterface, args []string) (interface{}, error) {

	if err := checkArgsForCount(args, 2); err != nil {
		return shim.Error(err.Error())
	}

	// 获取结构体标识符 euId
	euId, err := getArgsByfirst(args)

	if err != nil {
		return nil, err
	}

	// 获取结构体标识
	assetStruct, err := getStructByType(euId)
	if err != nil {
		return nil, err
	}

	return assetStruct, nil

}
//  generateSecurityLabelsUpload --校验-- 初始化二维码
func (p *Base) initializesTheSecurityLabelUpload(stub shim.ChaincodeStubInterface, args []string, assetStruct interface{}) (string, error) {


	/*

	arg[0] -- jiegouti 标识
	arg[1]  --  Uuid
	arg[2]  --  RelevKey
	Binding --  Binding	

	arg[]

	*/
	var err error
	if len(args) != 2 {
		return key, errors.New("Incorrect number of arguments. Expecting 3")
	}
	args_Uuid := args[1]
	updateTime := time.Now().Format("2006-01-02 15:04:05")

	logger.Debug("开始步骤")

	// 创建复合键
	comkey := createCompositeKey(stub, CompositeKeyByQr, args_Uuid)

	//	comkey
	assetAsBytes, err := stub.GetState(comkey)
	if err != nil {
		return "", errors.New("//１ comkey 获取失败 ")
	}

	if assetAsBytes != nil {
		return "", errors.New("//１ comkey 获取失败 ")
	}

	//１ 反序列化
	if err := json.Unmarshal([]byte(assetJsonStr), assetStruct); err != nil {
		return comkey, errors.New("//１ 反序列化:" + err.Error())
	}
	// 获取结构体标识 
	err	=getParentTypeByQr(assetStruct,arg[2])
	if err!=nil{
		return "",err
	}

	// 4 添加push时间
	if err = addUpdateTime(updateTime, assetStruct); err != nil {
		return comkey, errors.New("// 4 添加push时间:" + err.Error())
	}

	// 5　序列化
	assetAsBytes, err := json.Marshal(assetStruct)
	if err != nil {
		return comkey, errors.New("// 5　上链:" + err.Error())
	}
	if err = stub.PutState(comkey, assetAsBytes); err != nil {
		return comkey, err
	}

	return comkey, nil
}







// activationOfSecurityCode 参数处理 --校验-- 防伪标签激活
func (p *Base) activationOfSecurityCode(stub shim.ChaincodeStubInterface, args []string) (interface{}, error) {

	if err := checkArgsForCount(args, 2); err != nil {
		return shim.Error(err.Error())
	}

	// 获取结构体标识符 euId
	euId, err := getArgsByfirst(args)

	if err != nil {
		return nil, err
	}

	// 获取结构体标识
	assetStruct, err := getStructByType(euId)
	if err != nil {
		return nil, err
	}

	return assetStruct, nil

}
//  activationOfSecurityCodeUpload --校验-- 防伪标签激活
func (p *Base) activationOfSecurityCodeUpload(stub shim.ChaincodeStubInterface, args []string, assetStruct interface{}) (string, error) {


	/*

	arg[0]  --   结构体标识
	arg[1]  --  二维码 ID 
	arg[2]  --  关联物品 ID 
	anma    --  暗码


	*/
	var err error
	if len(args) != 3 {
		return "", errors.New("Incorrect number of arguments. Expecting 3")
	}
	args_uuid := args[1]
	args_relid:= args[2]
	args_priv:=	 args[3]
	updateTime := time.Now().Format("2006-01-02 15:04:05")

	logger.Debug("开始步骤")

	// 创建复合键
	comkey := createCompositeKey(stub, CompositeKeyByQr, args_Uuid)

	//	comkey
	assetAsBytes, err := stub.GetState(comkey)
	if err != nil {
		return "", errors.New("//１ comkey 获取失败 ")
	}

	if assetAsBytes != nil {
		return "", errors.New("//１ comkey 获取失败 ")
	}

	//１ 反序列化
	if err := json.Unmarshal([]byte(assetJsonStr), assetStruct); err != nil {
		return comkey, errors.New("//１ 反序列化:" + err.Error())
	}


	// 获取结构体标识   激活   ( qrid wpid state  暗码 )
	err = getParentTypeByQrActivate(assetStruct,args_uuid,args_relid,args_priv)

	if err!= nil{
		return "",err
	}

	// 4 添加push时间
	if err = addUpdateTime(updateTime, assetStruct); err != nil {
		return comkey, errors.New("// 4 添加push时间:" + err.Error())
	}

	// 5　序列化
	assetAsBytes, err := json.Marshal(assetStruct)
	if err != nil {
		return comkey, errors.New("// 5　上链:" + err.Error())
	}

	if err = stub.PutState(comkey, assetAsBytes); err != nil {
		return comkey, err
	}

	return comkey, nil
}




// activationOfSecurityCode 参数处理 --校验-- 中转
func (p *Base) theCirculationTransshipment(stub shim.ChaincodeStubInterface, args []string) (interface{}, error) {

	// 	 检验 参数个数 
	if err := checkArgsForCount(args, 4); err != nil {
		return shim.Error(err.Error())
	}

	// 获取结构体标识符 euId
	euId, err := getArgsByfirst(args)

	if err != nil {
		return nil, err
	}

	// 获取结构体标识
	assetStruct, err := getStructByType(euId)
	if err != nil {
		return nil, err
	}

	return assetStruct, nil

}
//  activationOfSecurityCodeUpload --校验-- 中转
func (p *Base) theCirculationTransshipmentUpload(stub shim.ChaincodeStubInterface, args []string, assetStruct interface{}) (string, error) {


	/*

	1. 	首先 暗码有没有效果 
	1.   物流信息 存储   + 暗码 符合键位 
	2.   暗码失效后，生成 新的明码 暗码  替换二维码 最新  然后把暗码发到 客户端 ， 明码存到链上 
	3.	 二维码更新 

	args[0]	二维码标识类型  
	args[1]	物流数据类型   
	args[2] 物流数据json  
	args[3] 暗码
	args[4] 二维码 uuid 明码校验 


	*/

	logger.info("开始步骤")
	
	updateTime := time.Now().Format("2006-01-02 15:04:05")
	
	args_tag_type := args[1]
	args_tag_json:= args[2]
	args_privateCode:=	 args[3]
	args_qrcode_uuid := args[4]


	// len(5)
	if len(args) != 5 {
		return "", errors.New("Incorrect number of arguments. Expecting 5")
	}

	
	// 创建二维码复合键位，首先根据二维码复合键位查询出 当前最新二维码结构数据 
	qr_comkey := createCompositeKey(stub, CompositeKeyByQr, []string(args_qrcode_uuid))

	//	根据 二维码 复合键位查询 结构数据 
	qr_assetAsBytes, err := stub.GetState(qr_comkey)

	if err != nil {
		return "", errors.New("// 1 根据 二维码 复合键位查询 结构数据 失败 ",err)
	}

	if qr_assetAsBytes != nil {
		return "", errors.New("// 2 qr_assetAsBytes 获取失败 ")
	}

	// 将查询到的 二维码字节数据反序列化到 二维码对应的结构
	if err := json.Unmarshal(qr_assetAsBytes, assetStruct); err != nil {
		return comkey, errors.New("//１ 反序列化:" + err.Error())
	}


	// 验证二维码 最新的明码是否和暗码匹配 (结构体指针,最新暗码)
	err = getParentTypeByQrValidation(assetStruct,args_privateCode)

	if err!= nil{
		return "",err
	}

	// 至此，二维码和暗码校验通过，开始物流数据上链处理流程

	// 获取物流结构体
	args_tag_type_struct, err := getStructByType(args_tag_type)
	if err != nil {
		return "", err
	}

	// 将 入参中 物流数据  反序列化到  物流结构体指针中 
	if err := json.Unmarshal([]byte(args_tag_json), args_tag_type_struct); err != nil {
		return key, errors.New("//１ 反序列化:" + err.Error())
	}


	args_tag_type_key, err = getkey(args_tag_type_struct)

	if err != nil {
		return args_tag_type_key, errors.New("//2 获取上链:" + err.Error())
	}

	// 4 添加push时间
	if err = addUpdateTime(updateTime, args_tag_type_struct); err != nil {
		return args_tag_type_key, errors.New("// 4 添加push时间:" + err.Error())
	}

	// 生成复合键位 物流复合键位   (  物流标识 || 组织标识 ，数据主键，暗码 )
	args_comkey := createCompositeKey(stub, CompositeKeyByWl, []string{CompositeKeyByOrg1,args_tag_type_key,args_privateCode})

	// 5　序列化 []byte{}
	args_assetAsBytes, err := json.Marshal(args_tag_type_struct)

	if err != nil {
		return args_comkey, errors.New("// 5　上链:" + err.Error())
	}

	//  物流数据存储

	if err = stub.PutState(args_comkey, args_assetAsBytes); err != nil {
		return args_comkey, err
	}

	//	返回物流数据 key  _ args_ 


	//  更新  二维码 明码  新增一对  原先 作废
	err = getParentTypeByQrAddNumberCode(assetStruct,args_privateCode)
	if err != nil {
		return args_comkey, errors.New("// 5　上链:" + err.Error())
	}

	// TODO  核心： 2个暗码一块存 ，相当于 存两遍， 一个是向上，一个是向下 
	qr_comkey := createCompositeKey(stub, CompositeKeyByQr, []string(args_qrcode_uuid,args_privateCode,args_privateCode))


	//  二维码 上链 
	//	 二维码 次数 + 1
	if err = stub.PutState(qr_comkey, assetStruct); err != nil {
		return args_comkey, err
	}



	
	// 总次数   复合键位 (numberByCodeSecurity + CompositeKeyByQr) +1 

	number_comkey := createCompositeKey(stub, numberByCodeSecurity, []string{CompositeKeyByQr})
	
		//	根据 二维码 复合键位查询 结构数据 
		number_assetAsBytes, err := stub.GetState(number_comkey)

		if err != nil {
			return "", errors.New("// 1 根据 numberByCodeSecurity 复合键位查询 结构数据 失败 ",err)
		}
	

		if qr_assetAsBytes != nil {
			return "", errors.New("// 2 number_assetAsBytes 获取失败 ")
		}	

		 string_qr_assetAsBytes = string(qr_assetAsBytes[:])

		// string  to  int64
		 int64_qr, err := strconv.ParseInt(string_qr_assetAsBytes, 10, 64)

		 if err !=nil{
			 return "", errors.New("strconv.ParseInt(string_qr_assetAsBytes, 10, 64) is err:",err)
		 }

		 int64_qr+=int64(1)

		 // int64到string
		string_qr := strconv.FormatInt(int64_qr,10)

		 if err = stub.PutState(number_comkey, []byte(string_qr)); err != nil {
			return args_comkey, err
		}



	//   ----     todo
	

	// //	 明码 + 二维码 uuid + 物品ID          二维码uuid  物品id + 明码 + 随机码 
	
	// number_comkey := createCompositeKey(stub, numberByCodeSecurity, []string{CompositeKeyByQr})



	return args_comkey, nil

}




// dataValidation 参数处理 --校验-- 验证二维码
func (p *Base) dataValidation(stub shim.ChaincodeStubInterface, args []string) (interface{}, error) {

	// 	 检验 参数个数 
	if err := checkArgsForCount(args, 4); err != nil {
		return shim.Error(err.Error())
	}

	// 获取结构体标识符 euId
	euId, err := getArgsByfirst(args)

	if err != nil {
		return nil, err
	}

	// 获取结构体标识
	assetStruct, err := getStructByType(euId)
	if err != nil {
		return nil, err
	}

	return assetStruct, nil

}
// dataValidationUpload --校验--  验证二维码
func (p *Base) dataValidationUpload(stub shim.ChaincodeStubInterface, args []string, assetStruct interface{}) (string, error) {


	/*
	1. 获取 二维码 uuid  暗码  

	2. 双重溯源：     
			二维码 + 物品ID + 明码  --- 物品ID   ( uuid+明码)
			
			

			暗码溯源：   根据 uuid + 暗码   查询   到 二维码   然后 根据 二维码 能查询 到 物流数据 ，然后根据二维码中下一个暗码，查询到下一个、

			根据 (uuid+ 暗码) 查询到的是 二维码 当前 二维码 +  
			如果要查询下一个 ，就是查询  next  
			如果要查询上一个 ， 

	*/

	logger.info("开始步骤")
	
	updateTime := time.Now().Format("2006-01-02 15:04:05")
	
	args_tag_type := args[1]
	args_tag_json:= args[2]
	args_privateCode:=	 args[3]
	args_qrcode_uuid := args[4]


	// len(5)
	if len(args) != 5 {
		return "", errors.New("Incorrect number of arguments. Expecting 5")
	}

	
	// 创建二维码复合键位，首先根据二维码复合键位查询出 当前最新二维码结构数据 
	qr_comkey := createCompositeKey(stub, CompositeKeyByQr, args_qrcode_uuid)

	//	根据 二维码 复合键位查询 结构数据 
	qr_assetAsBytes, err := stub.GetState(qr_comkey)

	if err != nil {
		return "", errors.New("// 1 根据 二维码 复合键位查询 结构数据 失败 ",err)
	}

	if qr_assetAsBytes != nil {
		return "", errors.New("// 2 qr_assetAsBytes 获取失败 ")
	}

	// 将查询到的 二维码字节数据反序列化到 二维码对应的结构
	if err := json.Unmarshal(qr_assetAsBytes, assetStruct); err != nil {
		return comkey, errors.New("//１ 反序列化:" + err.Error())
	}


	// 验证二维码 最新的明码是否和暗码匹配 (结构体指针,最新暗码)
	err = getParentTypeByQrValidation(assetStruct,args_privateCode)

	if err!= nil{
		return "",err
	}

	// 至此，二维码和暗码校验通过，开始物流数据上链处理流程

	// 获取物流结构体
	args_tag_type_struct, err := getStructByType(args_tag_type)
	if err != nil {
		return "", err
	}

	// 将 入参中 物流数据  反序列化到  物流结构体指针中 
	if err := json.Unmarshal([]byte(args_tag_json), args_tag_type_struct); err != nil {
		return key, errors.New("//１ 反序列化:" + err.Error())
	}


	args_tag_type_key, err = getkey(args_tag_type_struct)

	if err != nil {
		return args_tag_type_key, errors.New("//2 获取上链:" + err.Error())
	}

	// 4 添加push时间
	if err = addUpdateTime(updateTime, args_tag_type_struct); err != nil {
		return args_tag_type_key, errors.New("// 4 添加push时间:" + err.Error())
	}

	// 生成复合键位 物流复合键位   (  物流标识 || 组织标识 ，数据主键，暗码 )
	args_comkey := createCompositeKey(stub, CompositeKeyByWl, []string{CompositeKeyByOrg1,args_tag_type_key,args_privateCode})

	// 5　序列化 []byte{}
	args_assetAsBytes, err := json.Marshal(args_tag_type_struct)

	if err != nil {
		return args_comkey, errors.New("// 5　上链:" + err.Error())
	}

	//  物流数据存储

	if err = stub.PutState(args_comkey, args_assetAsBytes); err != nil {
		return args_comkey, err
	}

	//	返回物流数据 key  _ args_ 


	//  更新  二维码 明码  新增一对  原先 作废
	err = getParentTypeByQrAddNumberCode(assetStruct,args_privateCode)
	if err != nil {
		return args_comkey, errors.New("// 5　上链:" + err.Error())
	}


	//  二维码 上链 
	//	 二维码 次数 + 1
	if err = stub.PutState(qr_comkey, assetStruct); err != nil {
		return args_comkey, err
	}



	
	// 总次数   复合键位 (numberByCodeSecurity + CompositeKeyByQr) +1 

	number_comkey := createCompositeKey(stub, numberByCodeSecurity, []string{CompositeKeyByQr})
	
		//	根据 二维码 复合键位查询 结构数据 
		number_assetAsBytes, err := stub.GetState(number_comkey)

		if err != nil {
			return "", errors.New("// 1 根据 numberByCodeSecurity 复合键位查询 结构数据 失败 ",err)
		}
	

		if qr_assetAsBytes != nil {
			return "", errors.New("// 2 number_assetAsBytes 获取失败 ")
		}	

		 string_qr_assetAsBytes = string(qr_assetAsBytes[:])

		// string  to  int64
		 int64_qr, err := strconv.ParseInt(string_qr_assetAsBytes, 10, 64)

		 if err !=nil{
			 return "", errors.New("strconv.ParseInt(string_qr_assetAsBytes, 10, 64) is err:",err)
		 }

		 int64_qr+=int64(1)

		 // int64到string
		string_qr := strconv.FormatInt(int64_qr,10)

		 if err = stub.PutState(number_comkey, []byte(string_qr)); err != nil {
			return args_comkey, err
		}



	//   ----     todo
	

	// //	 明码 + 二维码 uuid + 物品ID          二维码uuid  物品id + 明码 + 随机码 
	
	// number_comkey := createCompositeKey(stub, numberByCodeSecurity, []string{CompositeKeyByQr})



	return args_comkey, nil

}






































//  getParentTypeByQr 设置关联物品ID 
func getParentTypeByQr(asset interface{},RelevKey string) error{

	switch t := asset.(type) {
		// 二维码 
	case *CodeSecurity:
		
		 if t.Binding == true {
			return  errors.New(" t Binding true ")
		 } 

		 t.RelevKey = RelevKey

		 t.Binding = true
			return  nil 
	default:
		log.Println("--------------default")

		return  errors.New("t is not type ")
	}
}


//  getParentTypeByQr 激活 
func getParentTypeByQrActivate(asset interface{},uuid string, relevKey string,privateCode string) error{

	switch t := asset.(type) {
		// 二维码 
	case *CodeSecurity:
		
		// 激活的前提： uuid  reid  binding = true  暗码  明码校对 

		// binding  判断是否绑定
		 if t.Binding == false {
			return  errors.New(" t Binding false ")
		 } 

		 // 判断 uuid 和  reid ID  是否正确 
		 if t.Uuid != uuid {
			return  errors.New("t.Uuid != uuid")
		 }

		 // 判断  物品ID 
		 if t.RelevKey != relevKey {
			return  errors.New("t.RelevKey != relevKey")
		 }

		 // 判断 暗码 校对  
		 if t.PrivateCode != privateCode {
			return  errors.New("t.PrivateCode != privateCode")
		 }


		 //  判断 明码 暗码校验
		 if t.PublicCode!="" && (t.PrivateCode != privateCode){
			 
			// 明码 暗码 校验 
			// PublicCode  privateCode  TODO 

			

		 }

			return  nil 
	default:
		log.Println("--------------default")

		return  errors.New("t is not type ")
	}
}

//  getParentTypeByQr 验证暗码是否匹配  
func getParentTypeByQrValidation(asset interface{},privateCode string) error{

	switch t := asset.(type) {
		// 二维码 
	case *CodeSecurity:

		// binding  判断是否绑定
		 if t.Binding == false {
			return  errors.New(" t Binding false ")
		 } 

		// 判断是否是激活之前 0  以及 流转最后结果 9999
		if t.State !=0 || t.State !=9999 {

			return errors.New("Stats is 0 betten 9999")
		}

		 // 判断 暗码 校对  
		 if t.PrivateCode != privateCode {
			return  errors.New("t.PrivateCode != privateCode")
		 }


		 //  判断 明码 暗码校验
		 if t.PublicCode!="" && (t.PrivateCode != privateCode){
			 
			// 明码 暗码 校验 
			// PublicCode  privateCode  TODO  非对称验证 
		 }

			return  nil 
	default:
		log.Println("--------------default")

		return  errors.New("t is not type ")
	}
}


//  getParentTypeByQr 二维码 +1 更新 明码， 暗码   
func getParentTypeByQrAddNumberCode(asset interface{},privateCode string) error{

	switch t := asset.(type) {
		// 二维码 
	case *CodeSecurity:

		// binding  判断是否绑定
		 if t.Binding == false {
			return  errors.New(" t Binding false ")
		 } 

		// 判断是否是激活之前 0  以及 流转最后结果 9999
		if t.State !=0 || t.State !=9999 {

			return errors.New("Stats is 0 betten 9999")
		}

		 // 判断 暗码 校对  
		 if t.PrivateCode != privateCode {
			return  errors.New("t.PrivateCode != privateCode")
		 }


		 //  判断 明码 暗码校验
		 if t.PublicCode!="" && (t.PrivateCode != privateCode){
			 
			// 明码 暗码 校验 
			// PublicCode  privateCode  TODO  非对称验证 
		 }


		 // 通过后， 1. 新增 明码，暗码， num +1 


		//	二维码 code 使用 次数 +1 
		 t.Count+=int64(1)

		// 更新明码，暗码 TODO 





			return  nil 
	default:
		log.Println("--------------default")

		return  errors.New("t is not type ")
	}
}


// response 返回结果统一处理
func (t *Base) response(errCode int, errMsg string, result interface{}) []byte {
	success := true
	if errCode != 0 {
		success = false
	} else {
		if errMsg == "" {
			errMsg = "Success"
		}
	}

	response := map[string]interface{}{
		"success":    success,
		"error_code": errCode,
		"error_msg":  errMsg,
		"result":     result,
	}
	retBytes, _ := json.Marshal(response)
	endTime := time.Now().UnixNano() / 1e6
	totalTime := endTime - startTime
	log.Printf("request token elapsed time %d response data %s", totalTime, string(retBytes))
	return retBytes

