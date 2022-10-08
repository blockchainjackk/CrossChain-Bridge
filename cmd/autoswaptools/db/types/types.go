package types

// AddressInfo
type AddressInfo struct {
	Key     string `json:"key"`
	Address string `json:"address"`
	Balance int64  `json:"balance"`
}

type SwapIn struct {
	TxId        string `json:"txid"`
	FromAddress string `json:"address"`
	ToAddress   int64  `json:"balance"`
	SignInfo    string `json:"sign_info"`
	Status      int64  `json:"status"`
}
type SwapOut struct {
	TxId   string `json:"txid"`
	From   string `json:"from"`
	Bind   string `json:"bind"`
	Status int64  `json:"status"`
}
