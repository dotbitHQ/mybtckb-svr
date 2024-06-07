package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/http_server/api_code"
	"mybtckb-svr/lib/common"
	"net/http"
)

type ReqClusterSporeList struct {
	api_code.PageData
	ClusterId string `json:"cluster_id" binding:"required"`
}

type RespClusterSporeList struct {
	api_code.PageData
	List []SporeInfo `json:"list"`
}
type SporeInfo struct {
	Address     string `json:"address"`
	SporeId     string `json:"spore_id"`
	ClusterId   string `json:"cluster_id"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
	Outpoint    string `json:"outpoint"`
	BlockNum    uint64 `json:"block_num"`
}

func (h *HttpHandle) ClusterSporeList(ctx *gin.Context) {
	var (
		funcName               = "ClusterSporeList"
		clientIp, remoteAddrIP = GetClientIp(ctx)
		req                    ReqClusterSporeList
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

	if err = h.doClusterSporeList(&req, &apiResp); err != nil {
		log.Error("doClusterSporeList err:", err.Error(), funcName, clientIp, ctx)
		apiResp.ApiRespErr(api_code.ApiCodeError500, "doClusterSporeList err: "+err.Error())
	}
	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doClusterSporeList(req *ReqClusterSporeList, apiResp *api_code.ApiResp) error {
	var resp RespClusterSporeList
	resp.List = make([]SporeInfo, 0)
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	res, total, err := h.DbDao.GetSporeByClusterId(req.ClusterId, req.Page, req.PageSize)
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
