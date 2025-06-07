package model

type InputData struct {
	ID       string                 `json:"id"`
	IDEditor string                 `json:"idEditor"`
	Status   string                 `json:"status"`
	Data     map[string]interface{} `json:"data"`
}

type Response struct {
	ID       string `json:"id"`
	IDEditor string `json:"idEditor"`
	Status   string `json:"status"`
	Time     string `json:"time"`
	SHA256   string `json:"sha256"`
	TxHash   string `json:"tx_hash"`
	Message  string `json:"message,omitempty"`
}
