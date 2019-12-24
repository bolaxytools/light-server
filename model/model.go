package model

type ChainInfo struct {
	LastBlockIndex string `json:"last_block_index"`
}

type Account struct {
	Nonce   uint64 `json:"nonce"`
	Balance uint64 `json:"balance"`
}

type Token struct {
	Contract string `db:"contract" json:"contract"`
	Symbol   string `db:"symbol" json:"symbol"`
	Logo     string `db:"logo" json:"logo"`
	Desc     string `db:"desc" json:"desc"`
	Followed bool   `db:"followed" json:"followed"`
}

type Follow struct {
	Contract string `db:"contract" json:"contract"`
	Wallet   string `db:"wallet" json:"wallet"`
	Balance  string `db:"balance" json:"balance"`
}

type Asset struct {
	Symbol   string `db:"symbol" json:"symbol"`
	Balance  string `db:"balance" json:"balance"`
	Contract string `db:"contract" json:"contract"`
	Logo     string `db:"logo" json:"logo"`
}
