package models

type Block struct {
	BlockNumber      int64         `json:"blockNumber"`
	Timestamp        uint64        `json:"timestamp"`
	Difficulty       uint64        `json:"difficulty"`
	Hash             string        `json:"hash"`
	TransactionCount int           `json:"transactionCount"`
	Transactions     []Transaction `json:"transactions"`
}

type Transaction struct {
	Hash     string `json:"hash"`
	Value    string `json:"value"`
	Gas      uint64 `json:"gas"`
	GasPrice uint64 `json:"gasPrice"`
	Nonce    uint64 `json:"nonce"`
	To       string `json:"to"`
	Pending  bool   `json:"pending"`
}

type TopAddress struct {
	Address string `json:"address"`
	Count   int    `json:"count"`
}

type Error struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
}
