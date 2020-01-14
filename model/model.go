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
	Decimals int32  `db:"decimals" json:"decimals"`
	Bap		 uint32 `db:"bap" json:"bap"`			//前端用于手续费的计算值
}

type Follow struct {
	Contract string `db:"contract" json:"contract"`
	Wallet   string `db:"wallet" json:"wallet"`
	Balance  string `db:"balance" json:"balance"`
	Followed bool   `db:"followed" json:"followed"`
}

type Asset struct {
	Symbol   string `db:"symbol" json:"symbol"`     //token简称
	Balance  string `db:"balance" json:"balance"`   //余额
	Contract string `db:"contract" json:"contract"` //合约地址
	Logo     string `db:"logo" json:"logo"`         //图标地址
	Desc     string `db:"desc" json:"desc"`         //token名称
	Decimals uint32 `db:"decimals" json:"decimals"` //精度
	Bap		 uint32 `db:"bap" json:"bap"`			//前端用于手续费的计算值
}

type FollowAsset struct {
	Address  string
	Contract string
}

type AlphaPing struct {
	Data interface{}
	Err string
}