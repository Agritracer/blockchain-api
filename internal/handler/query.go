package handler

import (
	"net/http"

	"agritrace/internal/storage"

	"github.com/gin-gonic/gin"
)

type QueryResponse struct {
	ID       string   `json:"id"`
	TxHashes []string `json:"tx_hashes"`
}

func HandleQuery(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu tham số id"})
		return
	}

	txs := storage.GetTxHashes(id)
	if len(txs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy giao dịch nào cho ID này"})
		return
	}

	resp := QueryResponse{
		ID:       id,
		TxHashes: txs,
	}

	c.JSON(http.StatusOK, resp)
}
