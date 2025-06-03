package storage

import (
	"sync"
)

var (
	mu      sync.RWMutex
	idToTxs = make(map[string][]string)
)

// Thêm một tx_hash cho ID
func AddTxHash(id string, txHash string) {
	mu.Lock()
	defer mu.Unlock()
	idToTxs[id] = append(idToTxs[id], txHash)
}

// Lấy toàn bộ tx_hash liên quan ID
func GetTxHashes(id string) []string {
	mu.RLock()
	defer mu.RUnlock()
	return idToTxs[id]
}
