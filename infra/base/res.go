package base

type ResCode int

const (
	ResCodeOk                ResCode = 1000
	ResCodeValidationErr     ResCode = 2000
	ResCodeRequestParamsErr  ResCode = 2100
	ResCodeInternalServerErr ResCode = 5000
	// 业务异常
	ResCodeBizErr                ResCode = 6000
	ResCodeBizTransferredFailure ResCode = 6010
)

type Res struct {
	Code    ResCode     `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
