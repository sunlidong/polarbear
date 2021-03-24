package router

import (
	"github.com/gin-gonic/gin"
)

const (
	cryptogenExtendOrderUrl       = "/exchange/cryptogen_extend"
	// exchangeCancelOrderUrl       = "/exchange/cancelOrder"
	// exchangeGetOneOrderUrl       = "/exchange/getOneOrder"
	// exchangeGetUnfinishOrdersUrl = "/exchange/getUnfinishOrders"
	// exchangeGetOrderHistorysUrl  = "/exchange/getOrderHistorys"
	// exchangeGetAccountUrl        = "/exchange/getAccount"
	// exchangeGetTickerUrl         = "/exchange/getTicker"
	// exchangeGetDepthUrl          = "/exchange/getDepth"
	// exchangeGetKlineRecordsUrl   = "/exchange/getKlineRecords"
	// exchangeGetTradesUrl         = "/exchange/getTrades"
	// exchangeGetExchangeDetailUrl = "/exchange/getExchangeDetail"
	// exchangePutUrl               = "/exchange/put"
	// exchangeListUrl              = "/exchange/list"
)

func cryptogenExtendRouter(router *gin.Engine) {
	router.POST(cryptogenExtendOrderUrl, cryptogenExtendOrderUrlPlace)
}

// @Summary 发送一个新的订单到某交易所进行撮合
// @Description 发送一个新的订单到某交易所进行撮合
// @Tags 交易所相关
// @Accept   json
// @Produce   json
// @Security token
// @Param group body model.OrderPlace true "每个参数均不得为空,//	OrderType: 0 表示limitBuy 1 表示limitSell 2 表示marketBuy 3 表示marketSell ; AccountType: 1 表示point 2 表示splot"
// @Success 200 {string} string "返回成功与否"
// @Router /exchange/orderPlace [post]
func cryptogenExtendOrderUrlPlace(c *gin.Context) {
	// res := result.NewResult()
	// defer c.JSON(http.StatusOK, res)
	// reqData, _ := c.GetRawData()
	// data, code, err := exchangeClient.OrderPlace(reqData)
	// res.Code = code
	// res.Msg = errs.GetMsg(code)
	// if err != nil {
	// 	res.Msg = err.Error()
	// 	logger.Error("错误信息：", err.Error())
	// 	return
	// }
	// res.Data = data
}
