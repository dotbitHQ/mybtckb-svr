package handle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/http_server/api_code"
	"mybtckb-svr/lib/common"
	"mybtckb-svr/lib/math/uint128"
	"mybtckb-svr/tables"
	"net/http"
)

type ReqUserXudtList struct {
	api_code.PageData
	Address string `json:"address" binding:"required"`
}

type RespUserXudtList struct {
	api_code.PageData
	List []TokenInfo `json:"list"`
}

type TokenInfo struct {
	TokenId string `json:"token_id"`
	Decimal uint8  `json:"decimal"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Amount  string `json:"amount"`
}

// @Summary 获取用户的xudt 资产列表

// @Description 获取资产列表1
// @Accept json
// @Produce json
// @Param address query string true "用户地址"
// @Success 200 {object} RespUserXudtList "成功"
// @Router /user/xudt/list [get]
func (h *HttpHandle) UserXudtList(ctx *gin.Context) {
	var (
		funcName               = "UserXudtList"
		clientIp, remoteAddrIP = GetClientIp(ctx)
		req                    ReqUserXudtList
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

	if err = h.doUserXudtList(&req, &apiResp); err != nil {
		log.Error("doUserXudtList err:", err.Error(), funcName, clientIp, ctx)
		apiResp.ApiRespErr(api_code.ApiCodeError500, "doUserXudtList err: "+err.Error())
	}
	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doUserXudtList(req *ReqUserXudtList, apiResp *api_code.ApiResp) error {
	var resp RespUserXudtList
	resp.List = make([]TokenInfo, 0)
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	parseAddr, err := address.Parse(req.Address)
	if err != nil {
		apiResp.ApiRespErr(api_code.ApiCodeParamsInvalid, "ckb_adress error")
		return nil
	}
	searchKey := &indexer.SearchKey{
		Script:     parseAddr.Script,
		ScriptType: indexer.ScriptTypeLock,
		Filter: &indexer.CellsFilter{
			Script: &types.Script{
				CodeHash: h.Contracts.XudtType.CodeHash,
				HashType: types.HashTypeType,
			},
		},
	}

	res, err := h.Contracts.Client().GetCells(h.Ctx, searchKey, indexer.SearchOrderDesc, 1000, "")
	if err != nil {
		apiResp.ApiRespErr(api_code.ApiCodeNetError, "GetCells err")
		return fmt.Errorf("GetCells err: %s", err.Error())
	}
	//发行xudt FAIR 0x79fa45b10a34cf767e9be38394ba62eb551988b3c8363b7d5ac5422af1071f13
	//transfer xudt TEST 0x89230a87dc2b66b44a5dd8d3ef9e2bdbddfb776f5ec754298e76c11971ad7cad
	if len(res.Objects) == 0 {
		apiResp.ApiRespOK(resp)
		return nil
	}
	fmt.Println(len(res.Objects))
	//return
	temp := make(map[string]uint128.Uint128)
	for _, v := range res.Objects {
		amount := uint128.Zero
		if len(v.OutputData) != 16 {
			continue
		}
		tokenId := common.Bytes2Hex(v.Output.Type.Args)
		temp[tokenId] = temp[tokenId].Add(uint128.FromBytes(v.OutputData))
		fmt.Println("0------amount: ", amount)
		fmt.Println("==============")
	}

	xudtInfoList, err := h.DbDao.GetXudtInfo()
	if err != nil {
		apiResp.ApiRespErr(api_code.ApiCodeDbError, "GetXudtInfo err")
		return err
	}
	//todo : cache xudt info to redis
	tempXudtInfo := make(map[string]tables.TableXudtInfo)
	for _, v := range xudtInfoList {
		tempXudtInfo[v.TokenId] = v
	}

	respTokenList := make([]TokenInfo, 0)
	for k, v := range temp {
		tokenInfo := TokenInfo{
			TokenId: k,
			Amount:  v.String(),
		}
		if _, ok := tempXudtInfo[k]; ok {
			tokenInfo.Name = tempXudtInfo[k].Name
			tokenInfo.Decimal = tempXudtInfo[k].Decimal
			tokenInfo.Symbol = tempXudtInfo[k].Symbol
		}
		respTokenList = append(respTokenList, tokenInfo)
	}

	resp.List = respTokenList
	resp.PageData = api_code.PageData{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    int64(len(respTokenList)),
	}
	apiResp.ApiRespOK(resp)
	return nil
}
