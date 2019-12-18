package model

type ChainInfo struct {
	LastBlockIndex string `json:"last_block_index"`
}

type Account struct {
	Nonce   uint64 `json:"nonce"`
	Balance uint64 `json:"balance"`
}
