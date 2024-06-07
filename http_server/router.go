package http_server

import (
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "mybtckb-svr/cmd/docs"
	"mybtckb-svr/config"
	"mybtckb-svr/http_server/middleware"
)

func (h *HttpServer) initRouter() {
	log.Info("initRouter:", len(config.Cfg.Origins))
	if len(config.Cfg.Origins) > 0 {
		middleware.AllowOriginList = append(middleware.AllowOriginList, config.Cfg.Origins...)
	}
	h.engine.Use(middleware.Cors())
	v1 := h.engine.Group("v1")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		v1.POST("/transfer/xudt/record", h.H.TransferXudtRecord)

		v1.POST("/user/xudt/list", h.H.UserXudtList)
		v1.POST("/user/cluster/list", h.H.UserClusterList)
		v1.POST("/cluster/spore/list", h.H.ClusterSporeList)
		v1.POST("/user/spore/list", h.H.UserSporeList)
		v1.POST("/transfer/xudt", h.H.TransferXudt)
		v1.POST("/transfer/spore", h.H.TransferSpore)

		v1.POST("/transaction/send", h.H.TransactionSend)

	}
}
