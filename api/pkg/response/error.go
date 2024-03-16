package response

const (
	ResponseErrorCode     = 1000000
	ParamParseErrorCode   = ResponseErrorCode + 1
	LogicErroeCode        = ResponseErrorCode + 2
	NotAllowedError       = ResponseErrorCode + 3
	TransferWowTokenError = ResponseErrorCode + 4
	TransferWowNftError   = ResponseErrorCode + 5
)
