package handler

import (
	"net/http"

	"agritrace/internal/config"
	"agritrace/internal/ethscan"

	"github.com/gin-gonic/gin"
)

func HandleList(c *gin.Context) {
	config.LoadConfig()
	address := config.Cfg.Wallet
	apiKey := config.Cfg.EthscanAPI

	if address == "" || apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Thiếu cấu hình WALLET_ADDRESS hoặc ETHERSCAN_API_KEY"})
		return
	}

	txs, err := ethscan.GetAllTxsByAddress(address, apiKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var hashes []string
	for _, tx := range txs {
		if tx.Input != "0x" {
			hashes = append(hashes, tx.Hash)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"wallet":    address,
		"tx_count":  len(hashes),
		"tx_hashes": hashes,
	})
}
