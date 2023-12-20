package service

const (
	OkCode  = 200
	ErrCode = -1
)

type ApiRsp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data",o`
}

func RspOk(data interface{}) *ApiRsp {
	return &ApiRsp{
		Code: OkCode,
		Msg:  "",
		Data: data,
	}
}
func RspErr(msg string) *ApiRsp {
	return &ApiRsp{
		Code: ErrCode,
		Msg:  msg,
	}
}
