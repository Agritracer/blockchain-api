package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"agritrace-api/internal/eth"
	"agritrace-api/internal/model"
)

func ProcessTrace(input model.InputData) (*model.Response, error) {
	jsonBytes, _ := json.Marshal(input.Data)
	hash := sha256.Sum256(jsonBytes)
	hashHex := hex.EncodeToString(hash[:])
	inputID := input.ID
	statusIF := input.Status
	inputIDEditor := input.IDEditor

	date := time.Now().Format("02/01/2006")
	content := "ID: " + inputID + "\nIDEditor:" + inputIDEditor + "\nSTATUS: " + statusIF + "\nDATE: " + date + "\nDATA: " + hashHex

	txHash, err := eth.SendToEthereum(content)
	if err != nil {
		return nil, err
	}

	return &model.Response{
		ID:       inputID,
		IDEditor: inputIDEditor,
		Status:   statusIF,
		Time:     date,
		SHA256:   hashHex,
		TxHash:   txHash,
	}, nil
}
