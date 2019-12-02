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
	Contract    string `db:"contract" json:"contract"`         //合约地址（如果有）
}

type Block struct {
	Height    uint64   `db:"height" json:"height"`         //区块高度
	Hash      string   `db:"hash" json:"hash"`             //区块hash
	TxCount   int32    `db:"tx_count" json:"tx_count"`     //交易数量
	BlockTime int64    `db:"block_time" json:"block_time"` //区块时间
	Signers   []string `db:"signers" json:"signers"`       //签名者
}

type Address struct {
	Addr       string `db:"addr",json:"addr"`               //地址
	AddTime    int64  `db:"add_time",json:"add_time"`       //添加时间
	UpdateTime int64  `db:"update_time",json:"update_time"` //更新时间
}

type TxReceipt struct {
	logs []ReceiptLog `json:"logs"`
}

type ReceiptLog struct {
	Topics []string `json:"topics"` //为3个的才关心 第2个和第3个分别是from和to
	Data   string   `json:"data"`
}
