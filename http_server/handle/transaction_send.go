package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/http_server/api_code"
	"mybtckb-svr/lib/contract"
	"mybtckb-svr/lib/txbuilder"
	"mybtckb-svr/tables"
	"net/http"
)

type ReqTransactionSend struct {
	txbuilder.SignInfoList
}

type RespTransactionSend struct {
	TxHash string `json:"tx_hash"`
}

func (h *HttpHandle) TransactionSend(ctx *gin.Context) {
	var (
		funcName               = "TransactionSendNew"
		clientIp, remoteAddrIP = GetClientIp(ctx)
		req                    ReqTransactionSend
		apiResp                api_code.ApiResp
		err                    error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName, clientIp, ctx)
		apiResp.ApiRespErr(api_code.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}

	log.Info("ApiReq:", funcName, clientIp, remoteAddrIP, gconv.String(req), ctx)

	if err = h.doTransactionSend(&req, &apiResp); err != nil {
		log.Error("doTransactionSend err:", err.Error(), funcName, clientIp, ctx)
		apiResp.ApiRespErr(api_code.ApiCodeError500, "doTransactionSend err: "+err.Error())
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doTransactionSend(req *ReqTransactionSend, apiResp *api_code.ApiResp) error {
	var resp RespTransactionSend

	txBuilder := txbuilder.NewTxBuilder(
		txbuilder.WithCache(h.TxTool.Cache()),
		txbuilder.WithTxBuilderBase(h.TxTool.TxBuilderBase()),
	)
	sic, err := txBuilder.GetTxBuilderFromCache(req.SignKey)
	if err != nil {
		return err
	}
	if err := txBuilder.AddSignatureForTx(req.SignList); err != nil {
		return err
	}
	hash, err := txBuilder.SendTransaction()
	if err != nil {
		return doSendTransactionError(err, apiResp)
	}
	txHash := hash.String()
	if err := h.doRecordTx(sic, txHash); err != nil {
		return err
	}
	resp.TxHash = txHash
	apiResp.ApiRespOK(resp)
	return nil
}
func (h *HttpHandle) doRecordTx(sic *txbuilder.SignInfoCache, txHash string) error {
	switch sic.Action {
	case contract.ActionTransferXudt:
		xudtTransferRecord := &tables.XudtTransferRecord{
			TxHash:         txHash,
			Address:        sic.TransferXudt.Address,
			ReceiptAddress: sic.TransferXudt.ReceiptAddress,
			Amount:         sic.TransferXudt.Amount,
			TokenId:        sic.TransferXudt.TokenId,
		}
		return h.DbDao.CreateXudtTransferRecord(xudtTransferRecord)
	case contract.ActionTransferSpore:
		sporeTransferRecord := &tables.SporeTransferRecord{
			TxHash:         txHash,
			Address:        sic.TransferSpore.Address,
			ReceiptAddress: sic.TransferSpore.ReceiptAddress,
			SporeId:        sic.TransferSpore.SporeId,
		}
		return h.DbDao.CreateSporeTransferRecord(sporeTransferRecord)
	case contract.ActionTransferCkb:
		ckbTransferRecord := &tables.CkbTransferRecord{
			TxHash:         txHash,
			Address:        sic.TransferCkb.Address,
			ReceiptAddress: sic.TransferCkb.ReceiptAddress,
			Amount:         sic.TransferCkb.Amount,
		}
		return h.DbDao.CreateCkbTransferRecord(ckbTransferRecord)
	}

	return nil
}
