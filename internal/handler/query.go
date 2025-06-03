package handler

import (
	"encoding/json"
	"net/http"

	"agritrace-api/internal/storage"
)

type QueryResponse struct {
	ID       string   `json:"id"`
	TxHashes []string `json:"tx_hashes"`
}

func HandleQuery(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Thiếu tham số id", http.StatusBadRequest)
		return
	}

	txs := storage.GetTxHashes(id)
	if len(txs) == 0 {
		http.Error(w, "Không tìm thấy giao dịch nào cho ID này", http.StatusNotFound)
		return
	}

	resp := QueryResponse{
		ID:       id,
		TxHashes: txs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
