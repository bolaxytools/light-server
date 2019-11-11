package resp

import "wallet-service/model"

type TxHistory struct {
	Txs []*model.Tx `json:"txs"`
}

func NewTxHistory(txs []*model.Tx) *TxHistory {
	return &TxHistory{Txs:txs}
}