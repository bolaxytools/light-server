package req

type ReqSendTx struct {
	Addr string `json:"addr"`
	SignedTx string `json:"signed_tx"`
}

type ReqHistory struct {
	ReqListBase
}