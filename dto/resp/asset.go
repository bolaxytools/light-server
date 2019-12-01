package resp

type Asset struct {
	Symbol   string `json:"symbol"`
	Balance  string `json:"balance"`
	Contract string `json:"contract"`
}

// swagger:response UpdateUserResponseWrapper
type AssetBox struct {
	MainCoin    *Asset   `json:"main_coin"`
	ExtCoinList []*Asset `json:"ext_coin_list"`
}

type NonceObj struct {
	Nonce uint64 `json:"nonce"`
}
