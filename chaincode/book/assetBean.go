package main

var (
	Cow Base{}
	startTime int64
)

// 增量合约　想关信息
const (

	incrementContractName string = "IncrementChaincode"
	// 增加
	incrementContractMethod string = "addIncrementInfo"
	// 查询
	incrementContractQueryMethod string = "getIncrementInfo"

	keySplitKey string = "+"

	resultSplitKey string = "=="

	// 复合键位
	CompositeKeyByQr string = "roleByQr"
	CompositeKeyByWl string = "roleByWl"


	CompositeKeyByOrg1 string = "roleforOrg1ByWl"
	CompositeKeyByOrg2 string = "roleforOrg2ByWl"


	assetCodeByQrCode     int64 = 10001		// 结构体鉴别列

	assetCodeByLogistics  int64 = 10002

	assetCodeByManagement int64 = 10003


	// 统计数量  https://1024tools.com/hash  sha256 3
	numberByCodeSecurity string  ="c097cf28a72e8da3dd46db2dabdb86fe"
	// numberByCodeSecurity + CodeSecurity



// 	1.  生成防伪标签
// 	2.	初始化防伪标签 Initializes the security label
// 	3.  防伪码激活 Activation of security code
// 	4.  流转中转 The circulation transshipment
// 	5.  验证 Data validation
// 	6.  查询二维码 Query qr code



// ------------- 异常操作
// 	1. 冻结标签 Freeze tag
// 	2. 启动标签 Thaw the label
// 	3. 作废标签 Invalid label



	//  -----------------------  start  func name 

	// fnv1a32  https://1024tools.com/hash

	TagforNew string = "generateSecurityLabel"
	TagforIni string = "initializesTheSecurityLabel"
	TagforAct string = "activationOfSecurityCode"
	TagforTra string = "circulationTransshipment"
	TagforVal string = "dataValidation"
	TagforVer string = "queryQrCode"


	TagforFre string = "freezeTagLable"
	TagforSta string = "startTagLabel"
	TagforCan string = "invalidLabel"



	//  -----------------------  end    func name 

)



type (

	//  This.comp  
	Base struct {
		Name string `json:"name"`
		Key string `json:"key"`
	}

	//  chaincode struct 
	BookChaincode struct {
	}


	// 二维码防伪标签
    CodeSecurity struct {
		No int64 `json:"no"` 	// 编号
		Uuid 	 string   `json:"uuid" pkey:""` // uuid 二维码唯一标识符 
		Uukey 	 string   `json:"Uukey"` // Uukey 二维码唯一标识符 
		Describe string   `json:"describe"` // 描述
		Binding  bool     `json:"binding"`// 是否绑定
		State  	 string   `json:"state"`// 状态   1 2 3 4 5 6 7 8 
		QrCode   string   `json:"qrCode"` // 二维码图案 展示  url 
		RelevKey string   `json:"relevKey"`// 关联溯源ID
		PublicCode string `json:"publicCode"` // 明码
		PrivateCode string `json:"privateCode"`// 暗码
		Count int64 `json:"count"`	// 更新次数
		UpdateDate string `json:"updateDate"`	// 更新时间 
	}


	// 物流节点 A 信息

	CodeLogistics struct {
		No int64 `json:"no"` 	// 编号
		Uuid string `json:"uuid" pkey:""`
		DataforOrder int64 `json:"datafor"` 	//
		ProductName string `json:"productName"`
		Number int64 `json:"number"`
		Price  int64 `json:"price"`
		RemittanceWay string `json:"remittanceWay"`
		TransportOrderNo string `json:"transportOrderNo"`
		EstimatedTimeOfArrival string `json:"estimatedtimeofarrival"`
		ActualArrivalTime string `json:"actualarrivaltime"`
		Freight string `json:"freight"`
		CreateUserName string `json:"createusername"`
		CreateUserId string `json:"createuserid"`
		CreateTime string `json:"createtime"`
		Lock bool `json:"lock"`
		State string `json:"state"`
	}
)
