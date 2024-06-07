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

type ReqTransferSpore struct {
	SporeId        string `json:"spore_id" binding:"required"`
	ReceiptAddress string `json:"receipt_address" binding:"required"`
}

func (h *HttpHandle) TransferSpore(ctx *gin.Context) {
	var (
		funcName               = "TransferXudt"
		clientIp, remoteAddrIP = GetClientIp(ctx)
		req                    ReqTransferSpore
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

	if err = h.doTransferSpore(&req, &apiResp); err != nil {
		log.Error("doTransferSpore err:", err.Error(), funcName, clientIp, ctx)
		apiResp.ApiRespErr(api_code.ApiCodeError500, "doTransferSpore err: "+err.Error())
	}
	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doTransferSpore(req *ReqTransferSpore, apiResp *api_code.ApiResp) error {
	txBuilder, err := h.TxTool.TransferSpore(&txtool.TransferSporeParams{
		SporeId: req.SporeId,

		ReceiptAddress: req.ReceiptAddress,
	})
	if err != nil {
		apiResp.ApiRespErr(api_code.ApiCodeError500, "Transfer err")
		return err
	}

	genSignInfoParam := &txbuilder.SignInfoCache{
		Action: contract.ActionTransferSpore,
		TransferSpore: &txbuilder.TransferSporeCache{
			//Address:        req.Address,
			ReceiptAddress: req.ReceiptAddress,
			SporeId:        req.SporeId,
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
