package req

import (
	"github.com/mitchellh/mapstructure"
)

/*
	每个请求的结构体前辍为Req,，方便写代码时的识别
	如：获取余额--ReqGetBalance
*/
type ReqData struct {
	Data interface{} `form:"data" json:"data"` //请求的数据包
	// 此处不能用 interface{}，要换为map[string]interface{} 
	//Data interface{} `form:"data" json:"data,omitempty"` //请求的数据包
	Sign string `form:"sign" json:"sign,omitempty"` //签名信息
}

type ReqBase struct {
	Addr string `json:"addr"`
}

type ReqTxHash struct {
	Txnash string `json:"txnash"`
}

type ReqBlockHeight struct {
	Height string `json:"height"` //区块高度
}

type ReqSearch struct {
	Content string `json:"content"`	//搜索内容
}

type ReqListBase struct {
	Addr string `json:"addr"`
	Page int32 `json:"page"`
	PageSize int32 `json:"page_size"`
}

func (reqdata *ReqData) GetData() map[string]interface{} {
	//mp := reqdata.Data
	mp := reqdata.Data.(map[string]interface{})
	return mp
}

func (reqdata *ReqData) Reverse(ptr interface{}) error {
	return mapstructure.Decode(reqdata.GetData(), ptr)
}
