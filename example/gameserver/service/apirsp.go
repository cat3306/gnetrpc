package service

const (
	OkCode  = 200
	ErrCode = -1
)

type ApiRsp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func (a *ApiRsp) Err(msg string) {
	a.Msg = msg
	a.Code = ErrCode
}

func (a *ApiRsp) Ok(data interface{}) *ApiRsp {
	a.Data = data
	a.Code = OkCode
	return a
}
