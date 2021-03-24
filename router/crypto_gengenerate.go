package router

import (
	"github.com/gin-gonic/gin"
)

const (
	cryptoGengenerateOrderPlaceUrl        = "/cryptoGengenerate/orderPlace"
	cryptoGengenerateCancelOrderUrl       = "/cryptoGengenerate/cancelOrder"
	cryptoGengenerateGetOneOrderUrl       = "/cryptoGengenerate/getOneOrder"
	cryptoGengenerateGetUnfinishOrdersUrl = "/cryptoGengenerate/getUnfinishOrders"
	cryptoGengenerateGetOrderHistorysUrl  = "/cryptoGengenerate/getOrderHistorys"
	cryptoGengenerateGetAccountUrl        = "/cryptoGengenerate/getAccount"
	cryptoGengenerateGetTickerUrl         = "/cryptoGengenerate/getTicker"
	cryptoGengenerateGetDepthUrl          = "/cryptoGengenerate/getDepth"
	cryptoGengenerateGetKlineRecordsUrl   = "/cryptoGengenerate/getKlineRecords"
	cryptoGengenerateGetTradesUrl         = "/cryptoGengenerate/getTrades"
	cryptoGengenerateGetcryptoGengenerateDetailUrl = "/cryptoGengenerate/getcryptoGengenerateDetail"
	cryptoGengeneratePutUrl               = "/cryptoGengenerate/put"
	cryptoGengenerateListUrl              = "/cryptoGengenerate/list"
)

func cryptoGengenerateRouter(router *gin.Engine) {
	router.POST(cryptoGengenerateOrderPlaceUrl, cryptoGengenerateOrderPlace)
}

// @Summary 发送一个新的订单到某交易所进行撮合
// @Description 发送一个新的订单到某交易所进行撮合
// @Tags 交易所相关
// @Accept   json
// @Produce   json
// @Security token
// @Param group body model.OrderPlace true "每个参数均不得为空,//	OrderType: 0 表示limitBuy 1 表示limitSell 2 表示marketBuy 3 表示marketSell ; AccountType: 1 表示point 2 表示splot"
// @Success 200 {string} string "返回成功与否"
// @Router /cryptoGengenerate/orderPlace [post]
func cryptoGengenerateOrderPlace(c *gin.Context) {
	// res := result.NewResult()
	// defer c.JSON(http.StatusOK, res)
	// reqData, _ := c.GetRawData()
	// data, code, err := cryptoGengenerateClient.OrderPlace(reqData)
	// res.Code = code
	// res.Msg = errs.GetMsg(code)
	// if err != nil {
	// 	res.Msg = err.Error()
	// 	logger.Error("错误信息：", err.Error())
	// 	return
	// }
	// res.Data = data
}
