package resp

import "wallet-svc/model"

const (
	Ret_Addr   = "dz"
	Ret_Hash   = "hx"
	Ret_Height = "kg"
)

type IndexRet struct {
	ChainId       string  `json:"chain_id"`        //圈子id
	BlockCount    uint64  `json:"block_count"`     //总区块数
	AddressCount  uint64  `json:"address_count"`   //总地址数
	MainCoinCount float64 `json:"main_coin_count"` //积分币总数
	TxCount       uint64  `json:"tx_count"`        //总交易数
	CrossMax      float64 `json:"cross_max"`       //总跨链币数
	GasCostCount  float64 `json:"gas_cost_count"`  //全网消耗gas总量

	Txs    []*model.Tx    `json:"txs"`    //最新交易
	Blocks []*model.Block `json:"blocks"` //最新区块
}

type SearchRet struct {
	RetType string      `json:"ret_type"` //dz-地址 hx-哈希 kg-块高
	Data    interface{} `json:"data"`     //具体数据
}

func NewSearchRet(tp string, data interface{}) *SearchRet {
	return &SearchRet{
		RetType: tp,
		Data:    data,
	}
}

/*
	地址资产信息
*/
type AssetsInfo struct {
	Total     uint64       `json:"total"`
	AssetList []*AssetInfo `json:"asset_list"` //资产列表
}

func NewAssetList(assets []*AssetInfo, total uint64) *AssetsInfo {
	return &AssetsInfo{
		AssetList: assets,
		Total:     total,
	}
}

type AssetInfo struct {
	Name     string  `json:"name"`     //名称
	Contract string  `json:"contract"` //合约地址
	Type     string  `json:"type"`     //类型
	Symbol   string  `json:"symbol"`   //币种简称
	Quantity float64 `json:"quantity"` //量
	Logo     string  `json:"logo"`     //图标
}
