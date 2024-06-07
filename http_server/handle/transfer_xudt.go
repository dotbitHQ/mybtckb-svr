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

type ReqTransfer struct {
	Address        string `json:"address" binding:"required"`
	ReceiptAddress string `json:"receipt_address" binding:"required"`
	TokenId        string `json:"token_id" binding:"required"`
	Amount         uint64 `json:"amount" binding:"gt=0"`
}

func (h *HttpHandle) TransferXudt(ctx *gin.Context) {
	var (
		funcName               = "TransferXudt"
		clientIp, remoteAddrIP = GetClientIp(ctx)
		req                    ReqTransfer
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

	if err = h.doTransferXudt(&req, &apiResp); err != nil {
		log.Error("doTransferXudt err:", err.Error(), funcName, clientIp, ctx)
		apiResp.ApiRespErr(api_code.ApiCodeError500, "doTransferXudt err: "+err.Error())
	}
	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doTransferXudt(req *ReqTransfer, apiResp *api_code.ApiResp) error {
	txBuilder, err := h.TxTool.Transfer(&txtool.TransferParams{
		TokenId:        req.TokenId,
		Amount:         req.Amount,
		Address:        req.Address,
		ReceiptAddress: req.ReceiptAddress,
	})
	if err != nil {
		apiResp.ApiRespErr(api_code.ApiCodeError500, "Transfer err")
		return err
	}

	genSignInfoParam := &txbuilder.SignInfoCache{
		Action: contract.ActionTransferXudt,
		TransferXudt: &txbuilder.TransferXudtCache{
			Address:        req.Address,
			ReceiptAddress: req.ReceiptAddress,
			TokenId:        req.TokenId,
			Amount:         req.Amount,
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

	txHash, _ := txBuilder.Transaction.ComputeHash()
	log.Info("doTransfer: ", txHash)

	apiResp.ApiRespOK(signInfo)
	return nil
}
