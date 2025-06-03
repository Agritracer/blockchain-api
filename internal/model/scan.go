package model

type Tx struct {
	Hash  string `json:"hash"`
	Input string `json:"input"`
	From  string `json:"from"`
	To    string `json:"to"`
	Time  string `json:"timeStamp"`
}

type EtherscanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []Tx   `json:"result"`
}
