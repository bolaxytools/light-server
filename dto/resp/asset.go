package resp

type Asset struct {
	Symbol string `json:"symbol"`
	Balance string `json:"balance"`
}

// swagger:response UpdateUserResponseWrapper
type AssetBox struct {
	MainCoin *Asset `json:"main_coin"`
	ExtCoinList []*Asset `json:"ext_coin_list"`
}