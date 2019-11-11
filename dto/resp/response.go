package resp

import "wallet-service/werror"

/*
	每个响应的结构体前辍为Resp，方便写代码时的识别
	如登录--RespLogin
*/
type RespData struct {
	ErrCode werror.ErrorG `json:"err_code"`
	ErrMsg  string        `json:"err_msg"`
	Data    interface{}   `json:"data"`
}

func BindJsonErrorResp(errMsg string) RespData {
	return RespData{
		ErrCode: werror.BindJsonError,
		ErrMsg:  errMsg,
	}
}

func NewErrorResp(errCode werror.ErrorG, errMsg string) RespData {
	return RespData{
		ErrCode: errCode,
		ErrMsg:  errMsg,
	}
}

func NewSuccessResp(data interface{}) RespData {
	return RespData{
		ErrCode: werror.Success,
		ErrMsg:"成功",
		Data:    data,
	}
}

func NewSuccess() RespData {
	return RespData{
		ErrCode: werror.Success,
		ErrMsg:"成功",
	}
}