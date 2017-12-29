package common

type RequestHeader struct {
	Version int32  `json:"ver,omitempty"`
	Addr    string `json:"addr,omitempty"`
	Nonce   int64  `json:"nonce,omitempty"`
	Pwd     string `json:"pwd,omitempty"`
	Time    int64  `json:"ts,omitempty"`
	Note    string `json:"note,omitempty"`
}

type ReturnValue struct {
	TxID    string `json:"txid,omitempty"`
	Message string `json:"msg,omitempty"`
}

type RequestCreateAccount struct {
	Pubkey string `json:"pubkey,omitempty"`
}
