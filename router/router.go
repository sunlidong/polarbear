package router

import (
	"github.com/gin-gonic/gin"
)

func GetAllRounters() *gin.Engine {
	router := gin.New()

	cryptoGengenerateRouter(router)     //创建
	cryptoShowtemplateRouter(router)    // 模板
	cryptoVersionOrderUrlRouter(router) //版本
	cryptogenExtendRouter(router)       // 追加
	return router
}
