package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
)

const (
	saltSize  = 16
	nonceSize = 12 // GCM nonce size
	keySize   = 32 // AES-256 key size
	scryptN   = 32768
	scryptR   = 8
	scryptP   = 1
)

// GenerateSalt tạo ra một giá trị salt ngẫu nhiên
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

// Encrypt mã hóa dữ liệu bằng AES-256-GCM, sử dụng khóa được tạo từ secret và salt bởi Scrypt
func Encrypt(data, secret, salt []byte) ([]byte, []byte, error) {
	key, err := scrypt.Key(secret, salt, scryptN, scryptR, scryptP, keySize)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive key from scrypt: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)

	return ciphertext, nonce, nil
}

// Decrypt giải mã dữ liệu bằng AES-256-GCM, sử dụng khóa được tạo từ secret và salt bởi Scrypt
func Decrypt(ciphertext, secret, salt, nonce []byte) ([]byte, error) {
	key, err := scrypt.Key(secret, salt, scryptN, scryptR, scryptP, keySize)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key from scrypt: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}