package resp

import "wallet-svc/model"

type TxHistory struct {
	Txs []*model.Tx `json:"txs"`
}

func NewTxHistory(txs []*model.Tx) *TxHistory {
	return &TxHistory{Txs:txs}
}