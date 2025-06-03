package handler

import (
	"encoding/json"
	"net/http"

	"agritrace-api/internal/eth"
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
