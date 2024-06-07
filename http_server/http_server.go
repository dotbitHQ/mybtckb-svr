package http_server

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	"mybtckb-svr/config"
	"mybtckb-svr/http_server/handle"
	"mybtckb-svr/lib/common"
	"net/http"
	"os"
)

type HttpServer struct {
	Ctx             context.Context
	Address         string
	InternalAddress string
	H               *handle.HttpHandle
	engine          *gin.Engine
	srv             *http.Server
	internalSrv     *http.Server
}

func (h *HttpServer) Run() {
	h.engine = gin.New()

	l := log.New()
	l.SetOutput(os.Stdout)
	if config.Cfg.Server.Net == common.NetTypeMainNet {
		l.SetFormatter(&log.JSONFormatter{})
	}
	h.engine.Use(ginlogrus.Logger(l), gin.Recovery())

	h.initRouter()

	h.srv = &http.Server{
		Addr:    h.Address,
		Handler: h.engine,
	}
	go func() {
		if err := h.srv.ListenAndServe(); err != nil {
			log.Error("http_server run err:", err)
		}
	}()
}

func (h *HttpServer) Shutdown() {
	if h.srv != nil {
		log.Warn("http server Shutdown ... ")
		if err := h.srv.Shutdown(h.Ctx); err != nil {
			log.Error("http server Shutdown err:", err.Error())
		}
	}
	if h.internalSrv != nil {
		log.Warn("http server internal Shutdown ... ")
		if err := h.internalSrv.Shutdown(h.Ctx); err != nil {
			log.Error("http server internal Shutdown err:", err.Error())
		}
	}
}
