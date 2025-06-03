package ethscan

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"agritrace-api/internal/config"
	"agritrace-api/internal/model"
)

func GetTransactionsByAddress(address string, targetID string) ([]model.Tx, error) {
	config.LoadConfig()
	apiKey := config.Cfg.EthscanAPI
	if apiKey == "" {
		return nil, fmt.Errorf("ETHERSCAN_API_KEY chưa được thiết lập")
	}

	url := fmt.Sprintf(
		"https://api-sepolia.etherscan.io/api?module=account&action=txlist&address=%s&sort=desc&apikey=%s",
		address,
		apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("lỗi truy vấn Etherscan: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var parsed model.EtherscanResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("lỗi phân tích JSON: %v", err)
	}

	if parsed.Status != "1" {
		return nil, fmt.Errorf("etherscan trả về lỗi: %s", parsed.Message)
	}

	var matched []model.Tx
	for _, tx := range parsed.Result {
		if tx.Input == "0x" {
			continue
		}

		dataBytes, err := hex.DecodeString(strings.TrimPrefix(tx.Input, "0x"))
		if err != nil {
			continue
		}

		dataStr := string(dataBytes)
		if strings.Contains(dataStr, targetID) {
			matched = append(matched, tx)
		}
	}

	return matched, nil
}

func GetAllTxsByAddress(address string, apiKey string) ([]model.Tx, error) {
	url := fmt.Sprintf(
		"https://api.etherscan.io/api?module=account&action=txlist&address=%s&sort=desc&apikey=%s",
		address, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var parsed model.EtherscanResponse
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return nil, err
	}

	if parsed.Status != "1" {
		return nil, fmt.Errorf("etherscan trả về lỗi: %s", parsed.Message)
	}

	return parsed.Result, nil
}
