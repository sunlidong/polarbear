package router

import (
	"github.com/gin-gonic/gin"
)

const (
	cryptoVersionOrderUrl       = "/exchange/cryptogen_extend"
)

func cryptoVersionOrderUrlRouter(router *gin.Engine) {
	router.POST(cryptoShowtemplateOrderUrl,cryptoVersionOrderUrlOrderUrlPlace)
}

func cryptoVersionOrderUrlOrderUrlPlace(c *gin.Context) {
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
