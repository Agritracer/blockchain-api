package handler

import (
	"encoding/json"
	"net/http"

	"agritrace-api/internal/config"
	"agritrace-api/internal/eth"
	"agritrace-api/internal/ethscan"
	"agritrace-api/internal/model"
	"agritrace-api/internal/service"
)

func HandleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ hỗ trợ POST", http.StatusMethodNotAllowed)
		return
	}

	var input model.InputData
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Lỗi JSON", http.StatusBadRequest)
		return
	}

	resp, err := service.ProcessTrace(input)
	if err != nil {
		http.Error(w, "Lỗi xử lý: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// storage.AddTxHash(input.ID, resp.TxHash)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func HandleTrace(w http.ResponseWriter, r *http.Request) {
	txHash := r.URL.Query().Get("tx")
	if txHash == "" {
		http.Error(w, "Thiếu ?tx=tx_hash", http.StatusBadRequest)
		return
	}

	data, err := eth.GetDataFromTransaction(txHash)
	if err != nil {
		http.Error(w, "Không thể truy xuất giao dịch: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"tx_hash": txHash,
		"data":    data,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func HandleTraceByID(w http.ResponseWriter, r *http.Request) {
	config.LoadConfig()

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Thiếu tham số ?id=", http.StatusBadRequest)
		return
	}

	address := config.Cfg.Wallet
	if address == "" {
		http.Error(w, "TRACE_WALLET_ADDRESS chưa được cấu hình", http.StatusInternalServerError)
		return
	}

	txs, err := ethscan.GetTransactionsByAddress(address, id)
	if err != nil {
		http.Error(w, "Lỗi truy vấn: "+err.Error(), http.StatusInternalServerError)
		return
	}

	hashes := make([]string, 0)
	for _, tx := range txs {
		hashes = append(hashes, tx.Hash)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":        id,
		"tx_count":  len(hashes),
		"tx_hashes": hashes,
	})
}
