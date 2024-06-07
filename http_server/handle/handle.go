package handle

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"mybtckb-svr/dao"
	"mybtckb-svr/http_server/api_code"
	"mybtckb-svr/lib/cache"
	"mybtckb-svr/lib/contract"
	"mybtckb-svr/txtool"
	"net"
	"strings"
)

type HttpHandle struct {
	Ctx          context.Context
	DbDao        *dao.DbDao
	RC           *cache.RedisCache
	Contracts    *contract.Contracts
	TxTool       *txtool.TxTool
	ServerScript *types.Script
}

func GetClientIp(ctx *gin.Context) (string, string) {
	clientIP := fmt.Sprintf("%v", ctx.Request.Header.Get("X-Real-IP"))
	remoteAddrIP, _, _ := net.SplitHostPort(ctx.Request.RemoteAddr)
	return clientIP, remoteAddrIP
}

func doSendTransactionError(err error, apiResp *api_code.ApiResp) error {
	if strings.Contains(err.Error(), "PoolRejectedDuplicatedTransaction") ||
		strings.Contains(err.Error(), "Dead(OutPoint(") ||
		strings.Contains(err.Error(), "Unknown(OutPoint(") ||
		(strings.Contains(err.Error(), "getInputCell") && strings.Contains(err.Error(), "not live")) {

		apiResp.ApiRespErr(api_code.ApiCodeRejectedOutPoint, "SendTransaction err")
		return fmt.Errorf("SendTransaction err: %s", err.Error())
	}

	apiResp.ApiRespErr(api_code.ApiCodeError500, "send tx err")
	return fmt.Errorf("SendTransaction err: %s", err.Error())
}

func doApiError(err error, apiResp *api_code.ApiResp) {
	if strings.Contains(err.Error(), "PoolRejectedDuplicatedTransaction") ||
		strings.Contains(err.Error(), "Dead(OutPoint(") ||
		strings.Contains(err.Error(), "Unknown(OutPoint(") ||
		(strings.Contains(err.Error(), "getInputCell") && strings.Contains(err.Error(), "not live")) {

		apiResp.ApiRespErr(api_code.ApiCodeRejectedOutPoint, "send tx err")
	}
}
