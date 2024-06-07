package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/http_server/api_code"
	"mybtckb-svr/lib/common"
	"net/http"
)

type ReqUserSporeList struct {
	api_code.PageData
	Address string `json:"address" binding:"required"`
}

type RespUserSporeList struct {
	api_code.PageData
	List []SporeInfo `json:"list"`
}

func (h *HttpHandle) UserSporeList(ctx *gin.Context) {
	var (
		funcName               = "ReqUserSporeList"
		clientIp, remoteAddrIP = GetClientIp(ctx)
		req                    ReqUserSporeList
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

	if err = h.doUserSporeList(&req, &apiResp); err != nil {
		log.Error("doClusterSporeList err:", err.Error(), funcName, clientIp, ctx)
		apiResp.ApiRespErr(api_code.ApiCodeError500, "doClusterSporeList err: "+err.Error())
	}
	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doUserSporeList(req *ReqUserSporeList, apiResp *api_code.ApiResp) error {
	var resp RespUserSporeList
	resp.List = make([]SporeInfo, 0)
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	res, total, err := h.DbDao.GetSporeByAddress(req.Address, req.Page, req.PageSize)
	if err != nil {
		apiResp.ApiRespErr(api_code.ApiCodeDbError, "db error")
		return nil
	}
	for _, v := range res {
		resp.List = append(resp.List, SporeInfo{
			Address:     v.Address,
			SporeId:     v.SporeId,
			ClusterId:   v.ClusterId,
			ContentType: v.ContentType,
			Content:     common.Bytes2Hex(v.Content),
			Outpoint:    v.Outpoint,
			BlockNum:    v.BlockNum,
		})
	}
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	resp.Total = total
	resp.GetTotalPadge()
	apiResp.ApiRespOK(resp)
	return nil
}
