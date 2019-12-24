package resp

import "wallet-svc/model"

// swagger:response UpdateUserResponseWrapper
type AssetBox struct {
	MainCoin    *model.Asset   `json:"main_coin"`
	ExtCoinList []*model.Asset `json:"ext_coin_list"`
}

type NonceObj struct {
	Nonce uint64 `json:"nonce"`
}

type SearchTokenRet struct {
	Total uint64 `json:"total"`
	TokenList []*model.Token `json:"token_list"` //token列表
}

func NewSearchTokenRet(tkns []*model.Token, total uint64) *SearchTokenRet {
	return &SearchTokenRet{
		Total:total,
		TokenList:tkns,
	}
}