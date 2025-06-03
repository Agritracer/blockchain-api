package model

type InputData struct {
	ID   string                 `json:"id"`
	Data map[string]interface{} `json:"data"`
}

type Response struct {
	SHA256  string `json:"sha256"`
	TxHash  string `json:"tx_hash"`
	Message string `json:"message,omitempty"`
}
