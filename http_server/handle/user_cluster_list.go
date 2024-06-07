package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/http_server/api_code"
	"net/http"
)

type ReqUserClusterList struct {
	api_code.PageData
	Address string `json:"address" binding:"required"`
}

type RespUserClusterList struct {
	api_code.PageData
	List []ClusterInfo `json:"list"`
}

type ClusterInfo struct {
	ClusterId   string `json:"cluster_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SporesNum   uint64 `json:"spores_num"`
	Outpoint    string `json:"outpoint"`
	BlockNum    uint64 `json:"block_num"`
}

func (h *HttpHandle) UserClusterList(ctx *gin.Context) {
	var (
		funcName               = "UserClusterList"
		clientIp, remoteAddrIP = GetClientIp(ctx)
		req                    ReqUserClusterList
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

	if err = h.doUserClusterList(&req, &apiResp); err != nil {
		log.Error("doUserClusterList err:", err.Error(), funcName, clientIp, ctx)
		apiResp.ApiRespErr(api_code.ApiCodeError500, "doUserClusterList err: "+err.Error())
	}
	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doUserClusterList(req *ReqUserClusterList, apiResp *api_code.ApiResp) error {
	var resp RespUserClusterList
	resp.List = make([]ClusterInfo, 0)
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	res, total, err := h.DbDao.GetClusterByAddress(req.Address, req.Page, req.PageSize)
	if err != nil {
		apiResp.ApiRespErr(api_code.ApiCodeDbError, "db error")
		return nil
	}
	for _, v := range res {
		resp.List = append(resp.List, ClusterInfo{
			ClusterId:   v.ClusterId,
			Name:        v.Name,
			Description: v.Description,
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
