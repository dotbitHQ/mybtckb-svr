package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/http_server/api_code"
	"mybtckb-svr/tables"
	"net/http"
)

type ReqTransferCkbRecord struct {
	api_code.PageData
	Address string `json:"address" binding:"required"`
}

type RespTransferCkbRecord struct {
	api_code.PageData
	List []TransferCkbInfo `json:"list"`
}

type TransferCkbInfo struct {
	Address        string          `json:"address"`
	ReceiptAddress string          `json:"receipt_address"`
	Amount         uint64          `json:"amount"`
	TxHash         string          `json:"tx_hash"`
	Status         tables.TxStatus `json:"status"`
	BlockNum       uint64          `json:"block_num"`
}

func (h *HttpHandle) TransferCkbRecord(ctx *gin.Context) {
	var (
		funcName               = "TransferCkbRecord"
		clientIp, remoteAddrIP = GetClientIp(ctx)
		req                    ReqTransferCkbRecord
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

	if err = h.doTransferCkbRecord(&req, &apiResp); err != nil {
		log.Error("doTransferCkbRecord err:", err.Error(), funcName, clientIp, ctx)
		apiResp.ApiRespErr(api_code.ApiCodeError500, "doTransferCkbRecord err: "+err.Error())
	}
	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doTransferCkbRecord(req *ReqTransferCkbRecord, apiResp *api_code.ApiResp) error {
	var resp RespTransferCkbRecord
	resp.List = make([]TransferCkbInfo, 0)
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	res, total, err := h.DbDao.GetAllCkbTransferRecordByAddress(req.Address, req.Page, req.PageSize)
	if err != nil {
		apiResp.ApiRespErr(api_code.ApiCodeDbError, "db error")
		return nil
	}
	for _, v := range res {
		resp.List = append(resp.List, TransferCkbInfo{
			Address:        v.Address,
			ReceiptAddress: v.ReceiptAddress,
			Amount:         v.Amount,
			TxHash:         v.TxHash,
			Status:         v.Status,
			BlockNum:       v.BlockNum,
		})
	}
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	resp.Total = total
	resp.GetTotalPadge()
	apiResp.ApiRespOK(resp)
	return nil
}
