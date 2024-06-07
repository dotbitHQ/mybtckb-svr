package handle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/http_server/api_code"
	"mybtckb-svr/lib/contract"
	"mybtckb-svr/lib/txbuilder"
	"mybtckb-svr/txtool"
	"net/http"
)

type ReqTransferCkb struct {
	Address        string `json:"address"  binding:"required"`
	Amount         uint64 `json:"amount" binding:"required"`
	ReceiptAddress string `json:"receipt_address" binding:"required"`
}

func (h *HttpHandle) TransferCkb(ctx *gin.Context) {
	var (
		funcName               = "TransferXudt"
		clientIp, remoteAddrIP = GetClientIp(ctx)
		req                    ReqTransferCkb
		apiResp                api_code.ApiResp
		err                    error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName, clientIp, remoteAddrIP, ctx)
		apiResp.ApiRespErr(api_code.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, clientIp, gconv.String(req), ctx)

	if err = h.doTransferCkb(&req, &apiResp); err != nil {
		log.Error("doTransferCkb err:", err.Error(), funcName, clientIp, ctx)
		apiResp.ApiRespErr(api_code.ApiCodeError500, "doTransferCkb err: "+err.Error())
	}
	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doTransferCkb(req *ReqTransferCkb, apiResp *api_code.ApiResp) error {
	amount := req.Amount * 1e8
	txBuilder, err := h.TxTool.TransferCkb(&txtool.TransferCkbParams{
		Amount:         amount,
		Address:        req.Address,
		ReceiptAddress: req.ReceiptAddress,
	})
	if err != nil {
		apiResp.ApiRespErr(api_code.ApiCodeError500, "Transfer err")
		return err
	}

	genSignInfoParam := &txbuilder.SignInfoCache{
		Action: contract.ActionTransferCkb,
		TransferCkb: &txbuilder.TransferCkbCache{
			Address:        req.Address,
			ReceiptAddress: req.ReceiptAddress,
			Amount:         amount,
		},
	}
	signInfo, err := txBuilder.GenSignInfo(genSignInfoParam)
	if err != nil {
		apiResp.ApiRespErr(api_code.ApiCodeError500, "GenSignInfo err")
		return err
	}

	signList, err := txBuilder.GenerateDigestListFromTx(nil)
	if err != nil {
		apiResp.ApiRespErr(api_code.ApiCodeError500, "GenerateDigestListFromTx err")
		return fmt.Errorf("GenerateDigestListFromTx err: %s", err.Error())
	}
	signInfo.SignList = signList

	log.Info("doTransferCkb: ", txBuilder.TxString())

	apiResp.ApiRespOK(signInfo)
	return nil
}
