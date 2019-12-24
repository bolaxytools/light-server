package resp

import "wallet-svc/model"

type TxHistory struct {
	Txs []*model.Tx `json:"txs"`
	Total uint64 `json:"total"`
}

func NewTxHistory(txs []*model.Tx,total uint64) *TxHistory {
	return &TxHistory{Txs:txs,
		Total:total}
}

type SendTxResp struct {
	TxHash string `json:"tx_hash"`
}

type BlockHistory struct {
	Blocks []*model.Block `json:"blocks"`
	Total uint64 `json:"total"`
}

func NewBlockHistory(txs []*model.Block,total uint64) *BlockHistory {
	return &BlockHistory{Blocks:txs,Total:total}
}