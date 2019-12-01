package model

import "math/big"

type ChainInfo struct {
	LastBlockIndex string `json:"last_block_index"`
}

type Account struct {
	Nonce   uint64  `json:"nonce"`
	Balance big.Int `json:"balance"`
}

type Token struct {
	Contract string `db:"contract" json:"contract"`
	Symbol   string `db:"symbol" json:"symbol"`
	Logo     string `db:"logo" json:"logo"`
	Desc     string `db:"desc" json:"desc"`
	Quantity uint64 `db:"quantity" json:"quantity"`
	Followed bool   `db:"followed" json:"followed"`
}

type Follow struct {
	Contract string `db:"contract" json:"contract"`
	Wallet   string `db:"wallet" json:"wallet"`
	Balance  string `db:"balance" json:"balance"`
}

type Asset struct {
	Symbol   string `db:"symbol" json:"symbol"`     //币简称
	Balance  string `db:"balance" json:"balance"`   //余额
	Contract string `db:"contract" json:"contract"` //合约地址
	Logo     string `db:"logo" json:"logo"`         //图标地址
	Desc     string `db:"desc" json:"desc"`         //币名称
	Decimals uint32 `db:"decimals" json:"decimals"` //精度
}
