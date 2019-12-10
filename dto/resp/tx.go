package resp

import "wallet-svc/model"

type TxHistory struct {
	Txs []*model.Tx `json:"txs"`
}

func NewTxHistory(txs []*model.Tx) *TxHistory {
	return &TxHistory{Txs:txs}
}

type SendTxResp struct {
	TxHash string `json:"tx_hash"`
}

type BlockHistory struct {
	Blocks []*model.Block `json:"blocks"`
}

func NewBlockHistory(txs []*model.Block) *BlockHistory {
	return &BlockHistory{Blocks:txs}
}