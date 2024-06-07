package api_code

type ApiCode = int

const (
	ApiCodeSuccess        ApiCode = 0
	ApiCodeError500       ApiCode = 500
	ApiCodeParamsInvalid  ApiCode = 10000
	ApiCodeMethodNotExist ApiCode = 10001
	ApiCodeDbError        ApiCode = 10002
	ApiCodeCacheError     ApiCode = 10003
	ApiCodeNetError       ApiCode = 10004

	ApiCodeTransactionNotExist ApiCode = 11001
	ApiCodeNotSupportAddress   ApiCode = 11005
	ApiCodeInsufficientBalance ApiCode = 11007
	ApiCodeTxExpired           ApiCode = 11008
	ApiCodeAmountInvalid       ApiCode = 11010
	ApiCodeRejectedOutPoint    ApiCode = 11011
	ApiCodeSyncBlockNumber     ApiCode = 11012
	ApiCodeOperationFrequent   ApiCode = 11013
	ApiCodeNotEnoughChange     ApiCode = 11014
	ApiCodeTransactionSendFail ApiCode = 11015

	ApiCodeTokenIdNotFound     ApiCode = 20001
	ApiCodeTransactionNotFound ApiCode = 20002
	ApiCodeOwnerAddressInvalid ApiCode = 20003

	ApiCodeSystemUpgrade ApiCode = 30019
)

const (
	TextSystemUpgrade = "The service is under maintenance, please try again later."
)

type ApiResp struct {
	ErrNo  ApiCode     `json:"err_no"`
	ErrMsg string      `json:"err_msg"`
	Data   interface{} `json:"data"`
}

func (a *ApiResp) ApiRespErr(errNo ApiCode, errMsg string) {
	a.ErrNo = errNo
	a.ErrMsg = errMsg
}

func (a *ApiResp) ApiRespOK(data interface{}) {
	a.ErrNo = ApiCodeSuccess
	a.Data = data
}

func ApiRespErr(errNo ApiCode, errMsg string) ApiResp {
	return ApiResp{
		ErrNo:  errNo,
		ErrMsg: errMsg,
		Data:   nil,
	}
}
