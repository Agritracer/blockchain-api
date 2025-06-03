package handler

import (
	"net/http"

	"agritrace-api/internal/config"
	"agritrace-api/internal/eth"
	"agritrace-api/internal/ethscan"
	"agritrace-api/internal/model"
	"agritrace-api/internal/service"

	"github.com/gin-gonic/gin"
)

func HandleSubmit(c *gin.Context) {
	var input model.InputData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lỗi JSON"})
		return
	}

	resp, err := service.ProcessTrace(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi xử lý: " + err.Error()})
		return
	}

	// storage.AddTxHash(input.ID, resp.TxHash)

	c.JSON(http.StatusOK, resp)
}

func HandleTrace(c *gin.Context) {
	txHash := c.Query("tx")
	if txHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu tham số ?tx=tx_hash"})
		return
	}

	data, err := eth.GetDataFromTransaction(txHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể truy xuất giao dịch: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tx_hash": txHash,
		"data":    data,
	})
}

func HandleTraceByID(c *gin.Context) {
	config.LoadConfig()

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu tham số ?id="})
		return
	}

	address := config.Cfg.Wallet
	if address == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "TRACE_WALLET_ADDRESS chưa được cấu hình"})
		return
	}

	txs, err := ethscan.GetTransactionsByAddress(address, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi truy vấn: " + err.Error()})
		return
	}

	hashes := make([]string, 0, len(txs))
	for _, tx := range txs {
		hashes = append(hashes, tx.Hash)
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        id,
		"tx_count":  len(hashes),
		"tx_hashes": hashes,
	})
}
