package model

type Tx struct {
	TxType      int16  `db:"tx_type" json:"tx_type"`           //交易类型 0-转账
	AddrFrom    string `db:"addr_from" json:"addr_from"`       //来自于谁
	AddrTo      string `db:"addr_to" json:"addr_to"`           //谁是接收者
	Amount      string `db:"amount" json:"amount"`             //金额
	MinerFee    string `db:"miner_fee" json:"miner_fee"`       //手续费
	TxHash      string `db:"tx_hash" json:"tx_hash"`           //交易hash
	BlockHeight uint64 `db:"block_height" json:"block_height"` //交易所在的区块块高
	TxTime      int64  `db:"tx_time" json:"tx_time"`           //交易时间(因交易上没有时间，取块生成的时间，即block.time)
	Memo        string `db:"memo" json:"memo"`                 //备注
}
