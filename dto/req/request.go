package req

import (
	"github.com/alecthomas/log4go"
	"reflect"
)

/*
	每个请求的结构体前辍为Req,，方便写代码时的识别
	如：获取余额--ReqGetBalance
*/
type ReqData struct {
	Data interface{} `form:"data" json:"data"` //请求的数据包
	Sign string `form:"sign" json:"sign,omitempty"` //签名信息
}

type ReqBase struct {
	Addr string `json:"addr"`
}

type ReqTokenInfo struct {
	Addr     string `json:"addr"`     //地址
	Contract string `json:"contract"` //合约地址
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

type ReqFollow struct {
	Addr     string `json:"addr"`
	Contract string `json:"contract"`
}

type ReqTxHash struct {
	Txnash string `json:"txnash"`
}

type ReqBlockHeight struct {
	Height uint64 `json:"height"` //区块高度
}

type ReqSearch struct {
	Addr    string `json:"addr"`    //搜索人的地址
	Content string `json:"content"` //搜索内容
}

type ReqListBase struct {
	Addr     string `json:"addr"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

func (reqdata *ReqData) GetData() map[string]interface{} {
	//mp := reqdata.Data
	mp := reqdata.Data.(map[string]interface{})
	return mp
}

func (reqdata *ReqData) Reverse(ptr interface{}) error {
	return reqdata.ReflectReverse(ptr)
}

func (reqdata *ReqData) ReflectReverse(ptr interface{}) error {

	elm := reflect.TypeOf(ptr).Elem()

	vlm := reflect.ValueOf(ptr).Elem()

	mpdata := reqdata.GetData()
	for i := 0; i < elm.NumField(); i++ {
		f := elm.Field(i)
		kd := f.Type.Kind()
		tag := f.Tag.Get("json")
		fld := vlm.Field(i)

		mv := mpdata[tag]
		if mv == nil {
			continue
		}
		log4go.Info("mv.type=%T\n", mv)

		switch kd {
		case reflect.String:
			fld.SetString(mpdata[tag].(string))
		case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			fld.SetUint(uint64(mpdata[tag].(float64)))

		case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			fld.SetInt(int64(mpdata[tag].(float64)))

		case reflect.Float64, reflect.Float32:
			fld.SetFloat(mpdata[tag].(float64))
		case reflect.Bool:
			fld.SetBool(mpdata[tag].(bool))
		default:
			log4go.Info("reflect to struct error,kd=%s\n", kd.String())
		}

	}

	return nil
}
